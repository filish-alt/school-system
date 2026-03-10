package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type ExamViolationRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewExamViolationRepository(db *sql.DB) *ExamViolationRepository {
	return &ExamViolationRepository{DB: db, Queries: q.New(db)}
}

func (r *ExamViolationRepository) Create(ctx context.Context, p q.CreateExamViolationParams) error {
	return r.Queries.CreateExamViolation(ctx, p)
}

func (r *ExamViolationRepository) ListAll(ctx context.Context, tenantID sql.NullString) ([]q.ListAllExamViolationsRow, error) {
	return r.Queries.ListAllExamViolations(ctx, tenantID)
}

func (r *ExamViolationRepository) ListBySession(ctx context.Context, sessionID string) ([]q.ExamViolation, error) {
	return r.Queries.ListExamViolationsBySession(ctx, sql.NullString{String: sessionID, Valid: true})
}
