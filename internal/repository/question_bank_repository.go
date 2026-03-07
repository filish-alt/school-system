package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type QuestionBankRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewQuestionBankRepository(db *sql.DB) *QuestionBankRepository {
	return &QuestionBankRepository{DB: db, Queries: q.New(db)}
}

func (r *QuestionBankRepository) Create(ctx context.Context, p q.CreateQuestionBankParams) error {
	return r.Queries.CreateQuestionBank(ctx, p)
}

func (r *QuestionBankRepository) ListByTeacher(ctx context.Context, teacherID string, limit, offset int64) ([]q.QuestionBank, error) {
	return r.Queries.ListQuestionBanksByTeacher(ctx, q.ListQuestionBanksByTeacherParams{
		CreatedByTeacherID: sql.NullString{String: teacherID, Valid: true},
		Limit:              limit,
		Offset:             offset,
	})
}

func (r *QuestionBankRepository) Delete(ctx context.Context, id string) error {
	return r.Queries.DeleteQuestionBank(ctx, id)
}

