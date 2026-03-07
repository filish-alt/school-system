package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type QuestionRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewQuestionRepository(db *sql.DB) *QuestionRepository {
	return &QuestionRepository{DB: db, Queries: q.New(db)}
}

func (r *QuestionRepository) Create(ctx context.Context, p q.CreateQuestionParams) error {
	return r.Queries.CreateQuestion(ctx, p)
}

func (r *QuestionRepository) Update(ctx context.Context, p q.UpdateQuestionParams) error {
	return r.Queries.UpdateQuestion(ctx, p)
}

func (r *QuestionRepository) Delete(ctx context.Context, id string) error {
	return r.Queries.DeleteQuestion(ctx, id)
}

func (r *QuestionRepository) ListByBank(ctx context.Context, bankID string, limit, offset int64) ([]q.Question, error) {
	return r.Queries.ListQuestionsByBank(ctx, q.ListQuestionsByBankParams{
		QuestionBankID: sql.NullString{String: bankID, Valid: true},
		Limit:          limit,
		Offset:         offset,
	})
}

