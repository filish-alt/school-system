package exam

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	examdto "school-exam/internal/dto/exam"
	"school-exam/internal/repository"
	q "school-exam/internal/sqlc/gen"
	"school-exam/internal/security"
)

type Usecase struct {
	ExamRepo    *repository.ExamRepository
	EQRepo      *repository.ExamQuestionRepository
	TeacherRepo *repository.TeacherRepository
	StudentRepo *repository.StudentRepository
	OptionRepo  *repository.OptionRepository
	Queries     *q.Queries
}

func NewUsecase(db *sql.DB, examRepo *repository.ExamRepository, eqRepo *repository.ExamQuestionRepository, teacherRepo *repository.TeacherRepository, studentRepo *repository.StudentRepository, optionRepo *repository.OptionRepository) *Usecase {
	return &Usecase{
		ExamRepo:    examRepo,
		EQRepo:      eqRepo,
		TeacherRepo: teacherRepo,
		StudentRepo: studentRepo,
		OptionRepo:  optionRepo,
		Queries:     q.New(db),
	}
}

func (u *Usecase) teacherInfo(ctx context.Context) (string, string, error) {
	c, ok := ctx.(*gin.Context)
	if !ok {
		return "", "", fmt.Errorf("invalid context")
	}
	v, ok := c.Get("claims")
	if !ok {
		return "", "", fmt.Errorf("claims missing")
	}
	claims := v.(*security.Claims)
	t, err := u.TeacherRepo.GetByUserID(ctx, claims.UserID)
	if err != nil {
		return "", "", err
	}
	return t.ID, valueOrEmpty(t.TenantID), nil
}

func valueOrEmpty(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func (u *Usecase) CreateExam(ctx context.Context, req examdto.CreateExamRequest) (string, error) {
	teacherID, tenantID, err := u.teacherInfo(ctx)
	if err != nil {
		return "", err
	}

	id := uuid.New().String()
	err = u.ExamRepo.Create(ctx, q.CreateExamParams{
		ID:                 id,
		TenantID:           sql.NullString{String: tenantID, Valid: tenantID != ""},
		Title:              sql.NullString{String: req.Title, Valid: true},
		SubjectID:          req.SubjectID,
		SectionID:          req.SectionID,
		CreatedByTeacherID: sql.NullString{String: teacherID, Valid: true},
		DurationMinutes:    req.DurationMinutes,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
		Status:             sql.NullString{String: "draft", Valid: true},
		TotalMarks:         sql.NullInt64{Int64: 0, Valid: true},
		ShuffleOptions:    sql.NullInt64{Int64: 0, Valid: true},
	})
	return id, err
}

func (u *Usecase) GetExam(ctx context.Context, id string) (examdto.GetExamResponse, error) {
	exam, err := u.ExamRepo.Get(ctx, id)
	if err != nil {
		return examdto.GetExamResponse{}, err
	}
	return u.getExamDetail(ctx, exam, true)
}

func (u *Usecase) GetTeacherExamQuestions(ctx context.Context, id string) (examdto.GetExamResponse, error) {
	teacherID, _, err := u.teacherInfo(ctx)
	if err != nil {
		return examdto.GetExamResponse{}, err
	}

	exam, err := u.ExamRepo.Get(ctx, id)
	if err != nil {
		return examdto.GetExamResponse{}, err
	}

	if exam.CreatedByTeacherID.String != teacherID {
		return examdto.GetExamResponse{}, fmt.Errorf("you are not authorized to view questions for this exam")
	}

	return u.getExamDetail(ctx, exam, true)
}

func (u *Usecase) GetStudentExam(ctx context.Context, id string) (examdto.GetExamResponse, error) {
	c, ok := ctx.(*gin.Context)
	if !ok {
		return examdto.GetExamResponse{}, fmt.Errorf("invalid context")
	}
	v, ok := c.Get("claims")
	if !ok {
		return examdto.GetExamResponse{}, fmt.Errorf("claims missing")
	}
	claims := v.(*security.Claims)

	student, err := u.StudentRepo.GetByUserID(ctx, claims.UserID)
	if err != nil {
		return examdto.GetExamResponse{}, err
	}

	exam, err := u.ExamRepo.Get(ctx, id)
	if err != nil {
		return examdto.GetExamResponse{}, err
	}

	if exam.SectionID != student.SectionID.String {
		return examdto.GetExamResponse{}, fmt.Errorf("this exam is not for your section")
	}

	if exam.Status.String != "published" {
		return examdto.GetExamResponse{}, fmt.Errorf("exam is not published")
	}

	return u.getExamDetail(ctx, exam, false)
}

func (u *Usecase) getExamDetail(ctx context.Context, exam q.Exam, isTeacher bool) (examdto.GetExamResponse, error) {
	rows, err := u.EQRepo.List(ctx, exam.ID)
	if err != nil {
		return examdto.GetExamResponse{}, err
	}

	var questions []examdto.ExamQuestionDetail
	for _, r := range rows {
		options, _ := u.OptionRepo.ListByQuestion(ctx, valueOrEmpty(r.QuestionID), 100, 0)
		
		if !isTeacher {
			for i := range options {
				options[i].IsCorrect = sql.NullInt64{Int64: 0, Valid: false}
			}
		}

		if exam.ShuffleOptions.Valid && exam.ShuffleOptions.Int64 == 1 {
			shuffleOptions(options)
		}

		questions = append(questions, examdto.ExamQuestionDetail{
			ID:              r.ID,
			QuestionID:      valueOrEmpty(r.QuestionID),
			QuestionText:    valueOrEmpty(r.QuestionText),
			Type:            valueOrEmpty(r.Type),
			Marks:           r.Marks.Int64,
			DifficultyLevel: valueOrEmpty(r.DifficultyLevel),
			OrderIndex:      r.OrderIndex.Int64,
			Options:         options,
		})
	}

	return examdto.GetExamResponse{
		Exam:      exam,
		Questions: questions,
	}, nil
}

func (u *Usecase) ListExams(ctx context.Context, page, pageSize int64) ([]q.Exam, error) {
	teacherID, _, err := u.teacherInfo(ctx)
	if err != nil {
		return nil, err
	}
	limit, offset := normalize(page, pageSize)
	return u.ExamRepo.ListByTeacher(ctx, teacherID, limit, offset)
}

func (u *Usecase) ListStudentExams(ctx context.Context) ([]q.Exam, error) {
	c, ok := ctx.(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("invalid context")
	}
	v, ok := c.Get("claims")
	if !ok {
		return nil, fmt.Errorf("claims missing")
	}
	claims := v.(*security.Claims)
	
	student, err := u.StudentRepo.GetByUserID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if !student.SectionID.Valid {
		return nil, fmt.Errorf("student section not found")
	}

	return u.ExamRepo.ListPublishedBySection(ctx, student.SectionID.String)
}

func (u *Usecase) UpdateExam(ctx context.Context, req examdto.UpdateExamRequest) error {
	exam, err := u.ExamRepo.Get(ctx, req.ID)
	if err != nil {
		return err
	}

	params := q.UpdateExamParams{
		ID:              req.ID,
		Title:           toNullString(req.Title, exam.Title),
		SubjectID:       valueOrOldString(req.SubjectID, exam.SubjectID),
		SectionID:       valueOrOldString(req.SectionID, exam.SectionID),
		DurationMinutes: valueOrOldInt(req.DurationMinutes, exam.DurationMinutes),
		StartTime:       valueOrOldTime(req.StartTime, exam.StartTime),
		EndTime:         valueOrOldTime(req.EndTime, exam.EndTime),
		ShuffleOptions:  toNullBoolInt(req.ShuffleOptions, exam.ShuffleOptions),
	}

	return u.ExamRepo.Update(ctx, params)
}

func (u *Usecase) UpdateStatus(ctx context.Context, id, status string) error {
	return u.ExamRepo.UpdateStatus(ctx, id, status)
}

func (u *Usecase) DeleteExam(ctx context.Context, id string) error {
	exam, err := u.ExamRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if exam.Status.Valid && exam.Status.String != "draft" {
		return fmt.Errorf("cannot delete non-draft exam")
	}
	return u.ExamRepo.Delete(ctx, id)
}

func (u *Usecase) AddQuestions(ctx context.Context, req examdto.AddQuestionsRequest) error {
	for _, qi := range req.Questions {
		err := u.EQRepo.Add(ctx, q.AddExamQuestionParams{
			ID:         uuid.New().String(),
			ExamID:     sql.NullString{String: req.ExamID, Valid: true},
			QuestionID: sql.NullString{String: qi.QuestionID, Valid: true},
			Marks:      sql.NullInt64{Int64: qi.Marks, Valid: true},
			OrderIndex: sql.NullInt64{Int64: qi.OrderIndex, Valid: true},
		})
		if err != nil {
			return err
		}
	}
	
	if req.ShuffleOptions {
		_ = u.Queries.UpdateExam(ctx, q.UpdateExamParams{
			ID: req.ExamID,
			ShuffleOptions: sql.NullInt64{Int64: 1, Valid: true},
		})
	}

	return u.ExamRepo.UpdateTotalMarks(ctx, req.ExamID)
}

func (u *Usecase) AddRandomQuestions(ctx context.Context, req examdto.AddRandomQuestionsRequest) error {
	questions, err := u.EQRepo.GetRandom(ctx, req.QuestionBankID, req.Count)
	if err != nil {
		return err
	}

	for i, qs := range questions {
		err := u.EQRepo.Add(ctx, q.AddExamQuestionParams{
			ID:         uuid.New().String(),
			ExamID:     sql.NullString{String: req.ExamID, Valid: true},
			QuestionID: sql.NullString{String: qs.ID, Valid: true},
			Marks:      sql.NullInt64{Int64: req.Marks, Valid: true},
			OrderIndex: sql.NullInt64{Int64: int64(i + 1), Valid: true},
		})
		if err != nil {
			return err
		}
	}

	if req.ShuffleOptions {
		_ = u.Queries.UpdateExam(ctx, q.UpdateExamParams{
			ID: req.ExamID,
			ShuffleOptions: sql.NullInt64{Int64: 1, Valid: true},
		})
	}
	
	return u.ExamRepo.UpdateTotalMarks(ctx, req.ExamID)
}

func (u *Usecase) RemoveQuestion(ctx context.Context, examID, eqID string) error {
	err := u.EQRepo.Remove(ctx, eqID)
	if err != nil {
		return err
	}
	return u.ExamRepo.UpdateTotalMarks(ctx, examID)
}

func shuffleOptions(options []q.QuestionOption) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})
}

func normalize(page, pageSize int64) (int64, int64) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}
	return pageSize, (page - 1) * pageSize
}

func toNullString(s *string, old sql.NullString) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return old
}

func toNullInt(i *int64, old sql.NullInt64) sql.NullInt64 {
	if i != nil {
		return sql.NullInt64{Int64: *i, Valid: true}
	}
	return old
}

func toNullTime(t *time.Time, old sql.NullTime) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return old
}

func toNullBoolInt(b *bool, old sql.NullInt64) sql.NullInt64 {
	if b != nil {
		val := int64(0)
		if *b {
			val = 1
		}
		return sql.NullInt64{Int64: val, Valid: true}
	}
	return old
}

func valueOrOldString(s *string, old string) string {
	if s != nil {
		return *s
	}
	return old
}

func valueOrOldInt(i *int64, old int64) int64 {
	if i != nil {
		return *i
	}
	return old
}

func valueOrOldTime(t *time.Time, old time.Time) time.Time {
	if t != nil {
		return *t
	}
	return old
}
