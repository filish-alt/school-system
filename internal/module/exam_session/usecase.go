package exam_session

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sessiondto "school-exam/internal/dto/student"
	"school-exam/internal/repository"
	q "school-exam/internal/sqlc/gen"
	"school-exam/internal/security"
)

type Usecase struct {
	SessionRepo *repository.ExamSessionRepository
	AnswerRepo  *repository.StudentAnswerRepository
	StudentRepo *repository.StudentRepository
	ExamRepo    *repository.ExamRepository
	Queries     *q.Queries
}

func NewUsecase(db *sql.DB, sRepo *repository.ExamSessionRepository, aRepo *repository.StudentAnswerRepository, stuRepo *repository.StudentRepository, eRepo *repository.ExamRepository) *Usecase {
	return &Usecase{
		SessionRepo: sRepo,
		AnswerRepo:  aRepo,
		StudentRepo: stuRepo,
		ExamRepo:    eRepo,
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

func (u *Usecase) StartSession(ctx context.Context, req sessiondto.StartSessionRequest) (string, error) {
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

func (u *Usecase) SaveAnswer(ctx context.Context, req sessiondto.SaveAnswerRequest) error {
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
		}
	}

	return u.SessionRepo.UpdateScore(ctx, q.UpdateExamSessionScoreParams{
		ID:         sessionID,
		TotalScore: sql.NullInt64{Int64: totalScore, Valid: true},
		EndTime:    sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
}

func (u *Usecase) GetSession(ctx context.Context, id string) (q.ExamSession, error) {
	session, err := u.SessionRepo.Get(ctx, id)
	if err == nil && session.Status.String == "in_progress" && time.Now().UTC().After(session.EndTime.Time) {
		_ = u.SessionRepo.UpdateStatus(ctx, id, "timed_out")
		session.Status.String = "timed_out"
	}
	return session, err
}

func toNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}
