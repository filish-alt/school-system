package repository

import (
	"context"
	"database/sql"
	"school-exam/internal/domain"
	"time"
)

type AcademicYearRepository struct {
	DB *sql.DB
}

func NewAcademicYearRepository(db *sql.DB) *AcademicYearRepository {
	return &AcademicYearRepository{DB: db}
}

func (r *AcademicYearRepository) GetFirstTenantID(ctx context.Context) (string, error) {
	var id string
	err := r.DB.QueryRowContext(ctx, "SELECT id FROM tenants WHERE status = 'active' ORDER BY created_at ASC LIMIT 1").Scan(&id)
	if err != nil {
		err = r.DB.QueryRowContext(ctx, "SELECT id FROM tenants ORDER BY created_at ASC LIMIT 1").Scan(&id)
	}
	return id, err
}

func (r *AcademicYearRepository) Create(ctx context.Context, y domain.AcademicYear) error {
	query := `INSERT INTO academic_years (id, tenant_id, year_name, start_date, end_date, is_active, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, y.ID, y.TenantID, y.YearName, y.StartDate, y.EndDate, y.IsActive, time.Now())
	return err
}

func (r *AcademicYearRepository) ListByTenant(ctx context.Context, tenantID string) ([]domain.AcademicYear, error) {
	var query string
	var rows *sql.Rows
	var err error

	if tenantID != "" {
		query = `SELECT id, tenant_id, year_name, start_date, end_date, is_active, created_at
		          FROM academic_years WHERE tenant_id = ? ORDER BY created_at DESC`
		rows, err = r.DB.QueryContext(ctx, query, tenantID)
	} else {
		query = `SELECT id, tenant_id, year_name, start_date, end_date, is_active, created_at
		          FROM academic_years ORDER BY created_at DESC`
		rows, err = r.DB.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.AcademicYear
	for rows.Next() {
		var y domain.AcademicYear
		var isActiveInt int
		var startDate, endDate sql.NullString
		if err := rows.Scan(&y.ID, &y.TenantID, &y.YearName, &startDate, &endDate, &isActiveInt, &y.CreatedAt); err != nil {
			return nil, err
		}
		if startDate.Valid {
			y.StartDate = &startDate.String
		}
		if endDate.Valid {
			y.EndDate = &endDate.String
		}
		y.IsActive = isActiveInt == 1
		list = append(list, y)
	}
	return list, nil
}

func (r *AcademicYearRepository) SetActive(ctx context.Context, tenantID, id string) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	resolvedTenantID := tenantID
	if resolvedTenantID == "" {
		_ = tx.QueryRowContext(ctx, "SELECT tenant_id FROM academic_years WHERE id = ?", id).Scan(&resolvedTenantID)
	}

	if resolvedTenantID != "" {
		if _, err := tx.ExecContext(ctx, "UPDATE academic_years SET is_active = 0 WHERE tenant_id = ?", resolvedTenantID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "UPDATE academic_years SET is_active = 1 WHERE tenant_id = ? AND id = ?", resolvedTenantID, id); err != nil {
			return err
		}
	} else {
		if _, err := tx.ExecContext(ctx, "UPDATE academic_years SET is_active = 1 WHERE id = ?", id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *AcademicYearRepository) Delete(ctx context.Context, tenantID, id string) error {
	if tenantID != "" {
		_, err := r.DB.ExecContext(ctx, "DELETE FROM academic_years WHERE tenant_id = ? AND id = ?", tenantID, id)
		return err
	}
	_, err := r.DB.ExecContext(ctx, "DELETE FROM academic_years WHERE id = ?", id)
	return err
}
