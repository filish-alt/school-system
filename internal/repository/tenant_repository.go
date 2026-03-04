package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type TenantRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{DB: db, Queries: q.New(db)}
}

func (r *TenantRepository) Create(ctx context.Context, id, name string, address, phone *string) error {
	var a, p sql.NullString
	if address != nil {
		a = sql.NullString{String: *address, Valid: true}
	}
	if phone != nil {
		p = sql.NullString{String: *phone, Valid: true}
	}
	return r.Queries.CreateTenant(ctx, q.CreateTenantParams{
		ID:      id,
		Name:    name,
		Address: a,
		Phone:   p,
	})
}

func (r *TenantRepository) Update(ctx context.Context, id, name string, address, phone *string) error {
	var a, p sql.NullString
	if address != nil {
		a = sql.NullString{String: *address, Valid: true}
	}
	if phone != nil {
		p = sql.NullString{String: *phone, Valid: true}
	}
	return r.Queries.UpdateTenant(ctx, q.UpdateTenantParams{
		ID:      id,
		Name:    name,
		Address: a,
		Phone:   p,
	})
}

func (r *TenantRepository) SetStatus(ctx context.Context, id, status string) error {
	return r.Queries.SetTenantStatus(ctx, q.SetTenantStatusParams{
		ID:     id,
		Status: sql.NullString{String: status, Valid: true},
	})
}

func (r *TenantRepository) GetByID(ctx context.Context, id string) (q.Tenant, error) {
	return r.Queries.GetTenantByID(ctx, id)
}

func (r *TenantRepository) List(ctx context.Context, limit, offset int64) ([]q.Tenant, error) {
	return r.Queries.ListTenants(ctx, q.ListTenantsParams{Limit: limit, Offset: offset})
}

func (r *TenantRepository) ListByStatus(ctx context.Context, status string, limit, offset int64) ([]q.Tenant, error) {
	return r.Queries.ListTenantsByStatus(ctx, q.ListTenantsByStatusParams{
		Status: sql.NullString{String: status, Valid: true},
		Limit:  limit,
		Offset: offset,
	})
}
