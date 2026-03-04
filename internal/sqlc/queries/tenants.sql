-- name: CreateTenant :exec
INSERT INTO tenants (id, name, address, phone, status) VALUES (?, ?, ?, ?, 'active');

-- name: GetTenantByID :one
SELECT id, name, address, phone, status, created_at FROM tenants WHERE id = ? LIMIT 1;

-- name: ListTenants :many
SELECT id, name, address, phone, status, created_at FROM tenants ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: ListTenantsByStatus :many
SELECT id, name, address, phone, status, created_at FROM tenants WHERE status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: UpdateTenant :exec
UPDATE tenants SET name = ?, address = ?, phone = ? WHERE id = ?;

-- name: SetTenantStatus :exec
UPDATE tenants SET status = ? WHERE id = ?;
