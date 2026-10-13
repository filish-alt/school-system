package repository

import (
	"context"
	"database/sql"
	"school-exam/internal/domain"
	"time"
)

type SystemSettingRepository struct {
	DB *sql.DB
}

func NewSystemSettingRepository(db *sql.DB) *SystemSettingRepository {
	return &SystemSettingRepository{DB: db}
}

func (r *SystemSettingRepository) Get(ctx context.Context, tenantID, key string) (*domain.SystemSetting, error) {
	query := `SELECT id, tenant_id, setting_key, setting_value, updated_at
	          FROM system_settings WHERE tenant_id = ? AND setting_key = ?`
	var s domain.SystemSetting
	err := r.DB.QueryRowContext(ctx, query, tenantID, key).Scan(&s.ID, &s.TenantID, &s.SettingKey, &s.SettingValue, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SystemSettingRepository) List(ctx context.Context, tenantID string) ([]domain.SystemSetting, error) {
	query := `SELECT id, tenant_id, setting_key, setting_value, updated_at
	          FROM system_settings WHERE tenant_id = ?`
	rows, err := r.DB.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.SystemSetting
	for rows.Next() {
		var s domain.SystemSetting
		if err := rows.Scan(&s.ID, &s.TenantID, &s.SettingKey, &s.SettingValue, &s.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *SystemSettingRepository) Upsert(ctx context.Context, id, tenantID, key, value string) error {
	query := `INSERT INTO system_settings (id, tenant_id, setting_key, setting_value, updated_at)
	          VALUES (?, ?, ?, ?, ?)
	          ON CONFLICT(tenant_id, setting_key) DO UPDATE SET
	          setting_value = excluded.setting_value,
	          updated_at = excluded.updated_at`
	_, err := r.DB.ExecContext(ctx, query, id, tenantID, key, value, time.Now())
	return err
}
