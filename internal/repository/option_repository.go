package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type OptionRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewOptionRepository(db *sql.DB) *OptionRepository {
	return &OptionRepository{DB: db, Queries: q.New(db)}
}

func (r *OptionRepository) Create(ctx context.Context, p q.CreateOptionParams) error {
	return r.Queries.CreateOption(ctx, p)
}

func (r *OptionRepository) Get(ctx context.Context, id string) (q.QuestionOption, error) {
	return r.Queries.GetOption(ctx, id)
}

func (r *OptionRepository) Update(ctx context.Context, p q.UpdateOptionParams) error {
	return r.Queries.UpdateOption(ctx, p)
}

func (r *OptionRepository) ResetCorrectOptions(ctx context.Context, qid string, excludeID string) error {
	return r.Queries.ResetCorrectOptions(ctx, q.ResetCorrectOptionsParams{
		QuestionID: sql.NullString{String: qid, Valid: true},
		ID:         excludeID,
	})
}

func (r *OptionRepository) Delete(ctx context.Context, id string) error {
	return r.Queries.DeleteOption(ctx, id)
}

func (r *OptionRepository) ListByQuestion(ctx context.Context, qid string, limit, offset int64) ([]q.QuestionOption, error) {
	return r.Queries.ListOptionsByQuestion(ctx, q.ListOptionsByQuestionParams{
		QuestionID: sql.NullString{String: qid, Valid: true},
		Limit:      limit,
		Offset:     offset,
	})
}

