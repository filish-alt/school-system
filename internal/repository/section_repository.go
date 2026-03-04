package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type SectionRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewSectionRepository(db *sql.DB) *SectionRepository {
	return &SectionRepository{DB: db, Queries: q.New(db)}
}

func (r *SectionRepository) Create(ctx context.Context, p q.CreateSectionParams) error {
	return r.Queries.CreateSection(ctx, p)
}

func (r *SectionRepository) GetByID(ctx context.Context, id string) (q.Section, error) {
	return r.Queries.GetSectionByID(ctx, id)
}

func (r *SectionRepository) ListByTenant(ctx context.Context, tenantID string, limit, offset int64) ([]q.Section, error) {
	return r.Queries.ListSectionsByTenant(ctx, q.ListSectionsByTenantParams{
		TenantID: sql.NullString{String: tenantID, Valid: true},
		Limit:    limit,
		Offset:   offset,
	})
}

func (r *SectionRepository) Update(ctx context.Context, p q.UpdateSectionParams) error {
	return r.Queries.UpdateSection(ctx, p)
}

func (r *SectionRepository) Delete(ctx context.Context, id string) error {
	return r.Queries.DeleteSection(ctx, id)
}

