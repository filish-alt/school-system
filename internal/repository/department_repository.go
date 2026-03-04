package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type DepartmentRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewDepartmentRepository(db *sql.DB) *DepartmentRepository {
	return &DepartmentRepository{DB: db, Queries: q.New(db)}
}

func (r *DepartmentRepository) Create(ctx context.Context, id, tenantID, name string) error {
	return r.Queries.CreateDepartment(ctx, q.CreateDepartmentParams{
		ID:       id,
		TenantID: sql.NullString{String: tenantID, Valid: true},
		Name:     name,
	})
}

func (r *DepartmentRepository) GetByID(ctx context.Context, id string) (q.Department, error) {
	return r.Queries.GetDepartmentByID(ctx, id)
}

func (r *DepartmentRepository) ListByTenant(ctx context.Context, tenantID string, limit, offset int64) ([]q.Department, error) {
	return r.Queries.ListDepartmentsByTenant(ctx, q.ListDepartmentsByTenantParams{
		TenantID: sql.NullString{String: tenantID, Valid: true},
		Limit:    limit,
		Offset:   offset,
	})
}

func (r *DepartmentRepository) Update(ctx context.Context, id, name string) error {
	return r.Queries.UpdateDepartment(ctx, q.UpdateDepartmentParams{
		ID:   id,
		Name: name,
	})
}

func (r *DepartmentRepository) Delete(ctx context.Context, id string) error {
	return r.Queries.DeleteDepartment(ctx, id)
}

