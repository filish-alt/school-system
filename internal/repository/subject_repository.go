package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type SubjectRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewSubjectRepository(db *sql.DB) *SubjectRepository {
	return &SubjectRepository{DB: db, Queries: q.New(db)}
}

func (r *SubjectRepository) Create(ctx context.Context, p q.CreateSubjectParams) error {
	return r.Queries.CreateSubject(ctx, p)
}

func (r *SubjectRepository) GetByID(ctx context.Context, id string) (q.Subject, error) {
	return r.Queries.GetSubjectByID(ctx, id)
}

func (r *SubjectRepository) ListByTenant(ctx context.Context, tenantID string, limit, offset int64) ([]q.Subject, error) {
	return r.Queries.ListSubjectsByTenant(ctx, q.ListSubjectsByTenantParams{
		TenantID: sql.NullString{String: tenantID, Valid: true},
		Limit:    limit,
		Offset:   offset,
	})
}

func (r *SubjectRepository) Update(ctx context.Context, p q.UpdateSubjectParams) error {
	return r.Queries.UpdateSubject(ctx, p)
}

func (r *SubjectRepository) Delete(ctx context.Context, id string) error {
	return r.Queries.DeleteSubject(ctx, id)
}

