package exam_session

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	studentdto "school-exam/internal/dto/student"
	"school-exam/internal/repository"
	q "school-exam/internal/sqlc/gen"
	"school-exam/internal/security"
)

type Usecase struct {
	SessionRepo *repository.ExamSessionRepository
	AnswerRepo  *repository.StudentAnswerRepository
	StudentRepo *repository.StudentRepository
	ExamRepo    *repository.ExamRepository
	QuestionRepo *repository.QuestionRepository
	OptionRepo  *repository.OptionRepository
	Queries     *q.Queries
}

func NewUsecase(db *sql.DB, sRepo *repository.ExamSessionRepository, aRepo *repository.StudentAnswerRepository, stuRepo *repository.StudentRepository, eRepo *repository.ExamRepository, qRepo *repository.QuestionRepository, oRepo *repository.OptionRepository) *Usecase {
	return &Usecase{
		SessionRepo: sRepo,
		AnswerRepo:  aRepo,
		StudentRepo: stuRepo,
		ExamRepo:    eRepo,
		QuestionRepo: qRepo,
		OptionRepo:  oRepo,
		Queries:     q.New(db),
	}
}

func (u *Usecase) studentInfo(ctx context.Context) (string, error) {
	c, ok := ctx.(*gin.Context)
	if !ok {
		return "", fmt.Errorf("invalid context")
	}
	v, ok := c.Get("claims")
	if !ok {
		return "", fmt.Errorf("claims missing")
	}
	claims := v.(*security.Claims)
	s, err := u.StudentRepo.GetByUserID(ctx, claims.UserID)
	if err != nil {
		return "", err
	}
	return s.ID, nil
}

func (u *Usecase) StartSession(ctx context.Context, req studentdto.StartSessionRequest) (string, error) {
	studentID, err := u.studentInfo(ctx)
	if err != nil {
		return "", err
	}

	exam, err := u.ExamRepo.Get(ctx, req.ExamID)
	if err != nil {
		return "", err
	}

	if exam.Status.String != "published" {
		return "", fmt.Errorf("exam is not published")
	}

	// Check for existing active session
	active, err := u.SessionRepo.GetActive(ctx, studentID, req.ExamID)
	if err == nil {
		return active.ID, nil // Resume existing session
	}

	id := uuid.New().String()
	now := time.Now().UTC()
	duration := time.Duration(exam.DurationMinutes) * time.Minute
	endTime := now.Add(duration)

	err = u.SessionRepo.Create(ctx, q.CreateExamSessionParams{
		ID:         id,
		ExamID:     sql.NullString{String: req.ExamID, Valid: true},
		StudentID:  sql.NullString{String: studentID, Valid: true},
		StartTime:  sql.NullTime{Time: now, Valid: true},
		EndTime:    sql.NullTime{Time: endTime, Valid: true},
		Status:     sql.NullString{String: "in_progress", Valid: true},
		TotalScore: sql.NullInt64{Int64: 0, Valid: true},
	})
	return id, err
}

func (u *Usecase) SaveAnswer(ctx context.Context, req studentdto.SaveAnswerRequest) error {
	session, err := u.SessionRepo.Get(ctx, req.SessionID)
	if err != nil {
		return err
	}

	if session.Status.String != "in_progress" {
		return fmt.Errorf("session is not in progress")
	}

	if time.Now().UTC().After(session.EndTime.Time) {
		_ = u.SessionRepo.UpdateStatus(ctx, req.SessionID, "timed_out")
		return fmt.Errorf("session has timed out")
	}

	// For upsert, we need the existing ID if it exists, or create a new one.
	// But our Upsert query uses ON CONFLICT(id). 
	// We should probably rely on (session_id, question_id) being unique in the logic.
	
	existing, err := u.AnswerRepo.Get(ctx, req.SessionID, req.QuestionID)
	id := uuid.New().String()
	if err == nil {
		id = existing.ID
	}

	return u.AnswerRepo.Upsert(ctx, q.UpsertStudentAnswerParams{
		ID:               id,
		SessionID:        sql.NullString{String: req.SessionID, Valid: true},
		QuestionID:       sql.NullString{String: req.QuestionID, Valid: true},
		AnswerText:       toNullString(req.AnswerText),
		SelectedOptionID: toNullString(req.SelectedOptionID),
		Score:            sql.NullInt64{Int64: 0, Valid: true}, // Reset score, will be calculated on submit
	})
}

func (u *Usecase) SubmitSession(ctx context.Context, sessionID string) error {
	session, err := u.SessionRepo.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	if session.Status.String != "in_progress" {
		return fmt.Errorf("session is not in progress")
	}

	answers, err := u.AnswerRepo.ListBySession(ctx, sessionID)
	if err != nil {
		return err
	}

	var totalScore int64
	for _, ans := range answers {
		correctOptionID, err := u.AnswerRepo.GetCorrectOption(ctx, ans.QuestionID.String)
		if err == nil && ans.SelectedOptionID.Valid && ans.SelectedOptionID.String == correctOptionID {
			marks, _ := u.AnswerRepo.GetQuestionMarks(ctx, session.ExamID.String, ans.QuestionID.String)
			totalScore += marks
			
			// Update individual answer score (optional but good for records)
			_ = u.AnswerRepo.Upsert(ctx, q.UpsertStudentAnswerParams{
				ID:               ans.ID,
				SessionID:        ans.SessionID,
				QuestionID:       ans.QuestionID,
				AnswerText:       ans.AnswerText,
				SelectedOptionID: ans.SelectedOptionID,
				Score:            sql.NullInt64{Int64: marks, Valid: true},
			})
		} else {
			// Ensure score is 0 if incorrect
			_ = u.AnswerRepo.Upsert(ctx, q.UpsertStudentAnswerParams{
				ID:               ans.ID,
				SessionID:        ans.SessionID,
				QuestionID:       ans.QuestionID,
				AnswerText:       ans.AnswerText,
				SelectedOptionID: ans.SelectedOptionID,
				Score:            sql.NullInt64{Int64: 0, Valid: true},
			})
		}
	}

	return u.SessionRepo.UpdateScore(ctx, q.UpdateExamSessionScoreParams{
		ID:         sessionID,
		TotalScore: sql.NullInt64{Int64: totalScore, Valid: true},
		EndTime:    sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
}

func (u *Usecase) GetSessionResult(ctx context.Context, sessionID string) (studentdto.ExamResultResponse, error) {
	studentID, err := u.studentInfo(ctx)
	if err != nil {
		return studentdto.ExamResultResponse{}, err
	}

	session, err := u.SessionRepo.Get(ctx, sessionID)
	if err != nil {
		return studentdto.ExamResultResponse{}, err
	}

	if session.StudentID.String != studentID {
		return studentdto.ExamResultResponse{}, fmt.Errorf("unauthorized access to session")
	}

	exam, err := u.ExamRepo.Get(ctx, session.ExamID.String)
	if err != nil {
		return studentdto.ExamResultResponse{}, err
	}

	answers, err := u.AnswerRepo.ListBySession(ctx, sessionID)
	if err != nil {
		return studentdto.ExamResultResponse{}, err
	}

	var breakdown []studentdto.AnswerReviewDetail
	for _, ans := range answers {
		question, _ := u.QuestionRepo.Get(ctx, ans.QuestionID.String)
		options, _ := u.OptionRepo.ListByQuestion(ctx, ans.QuestionID.String, 100, 0)
		correctOptionID, _ := u.AnswerRepo.GetCorrectOption(ctx, ans.QuestionID.String)
		maxMarks, _ := u.AnswerRepo.GetQuestionMarks(ctx, exam.ID, ans.QuestionID.String)

		isCorrect := ans.SelectedOptionID.Valid && ans.SelectedOptionID.String == correctOptionID

		breakdown = append(breakdown, studentdto.AnswerReviewDetail{
			QuestionID:       ans.QuestionID.String,
			QuestionText:     question.QuestionText.String,
			QuestionType:     question.Type.String,
			SelectedOptionID: ans.SelectedOptionID.String,
			CorrectOptionID:  correctOptionID,
			IsCorrect:        isCorrect,
			Score:            ans.Score.Int64,
			MaxMarks:         maxMarks,
			Options:          options,
		})
	}

	return studentdto.ExamResultResponse{
		Summary: studentdto.ResultSummaryResponse{
			SessionID:  session.ID,
			ExamID:     exam.ID,
			ExamTitle:  exam.Title.String,
			TotalScore: session.TotalScore.Int64,
			Status:     session.Status.String,
			StartTime:  session.StartTime.Time,
			EndTime:    session.EndTime.Time,
		},
		Breakdown: breakdown,
	}, nil
}

func (u *Usecase) GetSession(ctx context.Context, id string) (q.ExamSession, error) {
	session, err := u.SessionRepo.Get(ctx, id)
	if err == nil && session.Status.String == "in_progress" && time.Now().UTC().After(session.EndTime.Time) {
		_ = u.SessionRepo.UpdateStatus(ctx, id, "timed_out")
		session.Status.String = "timed_out"
	}
	return session, err
}

func (u *Usecase) ListMySessions(ctx context.Context) ([]studentdto.SessionResponse, error) {
	studentID, err := u.studentInfo(ctx)
	if err != nil {
		return nil, err
	}

	sessions, err := u.SessionRepo.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	var res []studentdto.SessionResponse
	for _, s := range sessions {
		res = append(res, studentdto.SessionResponse{
			ID:         s.ID,
			ExamID:     s.ExamID.String,
			StudentID:  s.StudentID.String,
			StartTime:  s.StartTime.Time,
			EndTime:    s.EndTime.Time,
			Status:     s.Status.String,
			TotalScore: s.TotalScore.Int64,
		})
	}
	return res, nil
}

func toNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}
