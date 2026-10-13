package repository

import (
	"context"
	"database/sql"
	"school-exam/internal/domain"
	"time"
)

type AuditLogRepository struct {
	DB *sql.DB
}

func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
	return &AuditLogRepository{DB: db}
}

func (r *AuditLogRepository) Log(ctx context.Context, log domain.AuditLog) error {
	query := `INSERT INTO audit_logs (id, tenant_id, user_id, action, entity_type, entity_id, details, ip_address, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, log.ID, log.TenantID, log.UserID, log.Action, log.EntityType, log.EntityID, log.Details, log.IPAddress, time.Now())
	return err
}

func (r *AuditLogRepository) List(ctx context.Context, tenantID string, limit, offset int) ([]domain.AuditLog, error) {
	query := `SELECT a.id, a.tenant_id, a.user_id, u.username, a.action, a.entity_type, a.entity_id, a.details, a.ip_address, a.created_at
	          FROM audit_logs a
	          LEFT JOIN users u ON a.user_id = u.id
	          WHERE a.tenant_id = ? OR a.tenant_id IS NULL
	          ORDER BY a.created_at DESC
	          LIMIT ? OFFSET ?`
	rows, err := r.DB.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.AuditLog
	for rows.Next() {
		var l domain.AuditLog
		var tenID, uID, uName, eType, eID, details, ip sql.NullString
		if err := rows.Scan(&l.ID, &tenID, &uID, &uName, &l.Action, &eType, &eID, &details, &ip, &l.CreatedAt); err != nil {
			return nil, err
		}
		if tenID.Valid { l.TenantID = &tenID.String }
		if uID.Valid { l.UserID = &uID.String }
		if uName.Valid { l.Username = &uName.String }
		if eType.Valid { l.EntityType = &eType.String }
		if eID.Valid { l.EntityID = &eID.String }
		if details.Valid { l.Details = &details.String }
		if ip.Valid { l.IPAddress = &ip.String }
		list = append(list, l)
	}
	return list, nil
}
