-- name: CreateDepartment :exec
INSERT INTO departments (id, tenant_id, name) VALUES (?, ?, ?);

-- name: GetDepartmentByID :one
SELECT id, tenant_id, name FROM departments WHERE id = ? LIMIT 1;

-- name: ListDepartmentsByTenant :many
SELECT id, tenant_id, name FROM departments WHERE tenant_id = ? ORDER BY name LIMIT ? OFFSET ?;

-- name: UpdateDepartment :exec
UPDATE departments SET name = ? WHERE id = ?;

-- name: DeleteDepartment :exec
DELETE FROM departments WHERE id = ?;

