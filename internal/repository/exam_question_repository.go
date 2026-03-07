package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type ExamQuestionRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewExamQuestionRepository(db *sql.DB) *ExamQuestionRepository {
	return &ExamQuestionRepository{DB: db, Queries: q.New(db)}
}

func (r *ExamQuestionRepository) Add(ctx context.Context, p q.AddExamQuestionParams) error {
	return r.Queries.AddExamQuestion(ctx, p)
}

func (r *ExamQuestionRepository) Remove(ctx context.Context, id string) error {
	return r.Queries.RemoveExamQuestion(ctx, id)
}

func (r *ExamQuestionRepository) List(ctx context.Context, examID string) ([]q.ListExamQuestionsRow, error) {
	return r.Queries.ListExamQuestions(ctx, sql.NullString{String: examID, Valid: true})
}

func (r *ExamQuestionRepository) GetRandom(ctx context.Context, bankID string, count int64) ([]q.Question, error) {
	return r.Queries.GetRandomQuestionsFromBank(ctx, q.GetRandomQuestionsFromBankParams{
		QuestionBankID: sql.NullString{String: bankID, Valid: true},
		Limit:          count,
	})
}
