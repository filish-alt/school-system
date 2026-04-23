package teacher

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	teacherdto "school-exam/internal/dto/teacher"
	"school-exam/internal/repository"
	q "school-exam/internal/sqlc/gen"
	"school-exam/internal/security"
)

type Usecase struct {
	Banks    *repository.QuestionBankRepository
	Questions *repository.QuestionRepository
	Options  *repository.OptionRepository
	Teachers *repository.TeacherRepository
	Queries  *q.Queries
}

func NewUsecase(db *sql.DB, banks *repository.QuestionBankRepository, questions *repository.QuestionRepository, options *repository.OptionRepository, teachers *repository.TeacherRepository) *Usecase {
	return &Usecase{Banks: banks, Questions: questions, Options: options, Teachers: teachers, Queries: q.New(db)}
}

func (u *Usecase) teacherFromContext(ctx context.Context) (string, string, error) {
	c := ctx.(*gin.Context)
	v, ok := c.Get("claims")
	if !ok {
		return "", "", fmt.Errorf("claims missing")
	}
	claims := v.(*security.Claims)
	t, err := u.Teachers.GetByUserID(ctx, claims.UserID)
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

func (u *Usecase) CreateQuestionBank(ctx context.Context, req teacherdto.CreateQuestionBankRequest) (string, error) {
	tid, _, err := u.teacherFromContext(ctx)
	if err != nil {
		return "", err
	}
	id := uuid.New().String()
	return id, u.Banks.Create(ctx, q.CreateQuestionBankParams{
		ID:                   id,
		TenantID:             sql.NullString{}, // optional; schema allows null
		SubjectID:            sql.NullString{String: req.SubjectID, Valid: true},
		CreatedByTeacherID:   sql.NullString{String: tid, Valid: true},
		Title:                sql.NullString{String: req.Title, Valid: true},
	})
}

func (u *Usecase) ListQuestionBanks(ctx context.Context, page, pageSize int64) ([]q.QuestionBank, error) {
	teacherID, _, err := u.teacherFromContext(ctx)
	if err != nil {
		return nil, err
	}
	limit, offset := normalize(page, pageSize)
	return u.Banks.ListByTeacher(ctx, teacherID, limit, offset)
}

func (u *Usecase) CreateQuestion(ctx context.Context, req teacherdto.CreateQuestionRequest) (string, error) {
	if !validQuestionType(req.Type) {
		return "", fmt.Errorf("invalid question type")
	}
	id := uuid.New().String()
	return id, u.Questions.Create(ctx, q.CreateQuestionParams{
		ID:             id,
		QuestionBankID: sql.NullString{String: req.QuestionBankID, Valid: true},
		Type:           sql.NullString{String: req.Type, Valid: true},
		QuestionText:   sql.NullString{String: req.QuestionText, Valid: true},
		ImageURL:       toNullStringSimple(req.ImageURL),
		Marks:          toNullInt64Simple(req.Marks),
		DifficultyLevel: toNullStringSimple(req.Difficulty),
	})
}

func (u *Usecase) UpdateQuestion(ctx context.Context, req teacherdto.UpdateQuestionRequest) error {
	existing, err := u.Questions.Get(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("question not found: %w", err)
	}

	if req.Type != nil && !validQuestionType(*req.Type) {
		return fmt.Errorf("invalid question type")
	}

	return u.Questions.Update(ctx, q.UpdateQuestionParams{
		ID:             req.ID,
		Type:           toNullString(req.Type, existing.Type),
		QuestionText:   toNullString(req.QuestionText, existing.QuestionText),
		ImageURL:       toNullString(req.ImageURL, existing.ImageURL),
		Marks:          toNullInt64(req.Marks, existing.Marks),
		DifficultyLevel: toNullStringWithDefault(req.Difficulty, existing.DifficultyLevel),
	})
}

func (u *Usecase) DeleteQuestion(ctx context.Context, id string) error {
	return u.Questions.Delete(ctx, id)
}

func (u *Usecase) ListQuestions(ctx context.Context, bankID string, page, pageSize int64) ([]q.Question, error) {
	limit, offset := normalize(page, pageSize)
	return u.Questions.ListByBank(ctx, bankID, limit, offset)
}

func (u *Usecase) CreateOption(ctx context.Context, req teacherdto.CreateOptionRequest) (string, error) {
	id := uuid.New().String()
	
	if req.IsCorrect {
		_ = u.Options.ResetCorrectOptions(ctx, req.QuestionID, id)
	}

	return id, u.Options.Create(ctx, q.CreateOptionParams{
		ID:         id,
		QuestionID: sql.NullString{String: req.QuestionID, Valid: true},
		OptionText: sql.NullString{String: req.OptionText, Valid: true},
		IsCorrect:  toNullBoolSimple(req.IsCorrect),
	})
}

func (u *Usecase) UpdateOption(ctx context.Context, req teacherdto.UpdateOptionRequest) error {
	existing, err := u.Options.Get(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("option not found: %w", err)
	}

	if req.IsCorrect != nil && *req.IsCorrect {
		_ = u.Options.ResetCorrectOptions(ctx, existing.QuestionID.String, req.ID)
	}

	return u.Options.Update(ctx, q.UpdateOptionParams{
		ID:         req.ID,
		OptionText: toNullString(req.OptionText, existing.OptionText),
		IsCorrect:  toNullBoolOpt(req.IsCorrect, existing.IsCorrect),
	})
}


func (u *Usecase) DeleteOption(ctx context.Context, id string) error {
	return u.Options.Delete(ctx, id)
}

func (u *Usecase) ListOptions(ctx context.Context, questionID string, page, pageSize int64) ([]q.QuestionOption, error) {
	limit, offset := normalize(page, pageSize)
	return u.Options.ListByQuestion(ctx, questionID, limit, offset)
}

func (u *Usecase) ListMyStudents(ctx context.Context, page, pageSize int64) ([]q.Student, error) {
	teacherID, _, err := u.teacherFromContext(ctx)
	if err != nil {
		return nil, err
	}
	limit, offset := normalize(page, pageSize)
	return u.Teachers.ListStudents(ctx, teacherID, limit, offset)
}

func (u *Usecase) ListMyAssignments(ctx context.Context, page, pageSize int64) ([]q.ListTeacherAssignmentsByTeacherRow, error) {
	teacherID, _, err := u.teacherFromContext(ctx)
	if err != nil {
		return nil, err
	}
	limit, offset := normalize(page, pageSize)
	return u.Teachers.ListAssignmentsDetailedByTeacher(ctx, teacherID, limit, offset)
}

func (u *Usecase) ImportQuestionsFromCSV(ctx context.Context, bankID string, records [][]string) error {
	importedCount := 0
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}
		if len(record) < 7 {
			continue
		}

		questionText := strings.TrimSpace(record[0])
		if questionText == "" {
			continue
		}

		optionA := strings.TrimSpace(record[1])
		optionB := strings.TrimSpace(record[2])
		optionC := strings.TrimSpace(record[3])
		optionD := strings.TrimSpace(record[4])
		correctAnswer := strings.ToUpper(strings.TrimSpace(record[5]))
		marksStr := strings.TrimSpace(record[6])

		marks, _ := strconv.ParseInt(marksStr, 10, 64)

		// 1. Create Question
		qid := uuid.New().String()
		err := u.Questions.Create(ctx, q.CreateQuestionParams{
			ID:             qid,
			QuestionBankID: sql.NullString{String: bankID, Valid: true},
			Type:           sql.NullString{String: "mcq", Valid: true},
			QuestionText:   sql.NullString{String: questionText, Valid: true},
			Marks:          sql.NullInt64{Int64: marks, Valid: true},
			DifficultyLevel: sql.NullString{String: "medium", Valid: true},
		})
		if err != nil {
			return fmt.Errorf("row %d: %w", i+1, err)
		}

		// 2. Create Options
		options := []struct {
			text   string
			letter string
		}{
			{optionA, "A"},
			{optionB, "B"},
			{optionC, "C"},
			{optionD, "D"},
		}

		for _, opt := range options {
			isCorrect := int64(0)
			if opt.letter == correctAnswer {
				isCorrect = 1
			}

			oid := uuid.New().String()
			err = u.Options.Create(ctx, q.CreateOptionParams{
				ID:         oid,
				QuestionID: sql.NullString{String: qid, Valid: true},
				OptionText: sql.NullString{String: opt.text, Valid: true},
				IsCorrect:  sql.NullInt64{Int64: isCorrect, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("row %d option %s: %w", i+1, opt.letter, err)
			}
		}
		importedCount++
	}

	if importedCount == 0 {
		return fmt.Errorf("no valid question rows found in CSV (expected 7 columns: question, opt_a, opt_b, opt_c, opt_d, answer, marks)")
	}
	return nil
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

func toNullStringSimple(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}

func toNullStringWithDefault(s *string, old sql.NullString) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return old
}

func toNullInt64(i *int64, old sql.NullInt64) sql.NullInt64 {
	if i != nil {
		return sql.NullInt64{Int64: *i, Valid: true}
	}
	return old
}

func toNullInt64Simple(i *int64) sql.NullInt64 {
	if i != nil {
		return sql.NullInt64{Int64: *i, Valid: true}
	}
	return sql.NullInt64{}
}

func toNullBoolOpt(b *bool, old sql.NullInt64) sql.NullInt64 {
	if b != nil {
		val := int64(0)
		if *b {
			val = 1
		}
		return sql.NullInt64{Int64: val, Valid: true}
	}
	return old
}

func toNullBoolSimple(b bool) sql.NullInt64 {
	val := int64(0)
	if b {
		val = 1
	}
	return sql.NullInt64{Int64: val, Valid: true}
}

func validQuestionType(t string) bool {
	switch strings.ToLower(t) {
	case "mcq", "true_false", "short", "essay":
		return true
	default:
		return false
	}
}
