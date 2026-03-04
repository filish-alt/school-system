-- name: CreateSubject :exec
INSERT INTO subjects (id, tenant_id, name, department_id) VALUES (?, ?, ?, ?);

-- name: GetSubjectByID :one
SELECT id, tenant_id, name, department_id FROM subjects WHERE id = ? LIMIT 1;

-- name: ListSubjectsByTenant :many
SELECT id, tenant_id, name, department_id FROM subjects WHERE tenant_id = ? ORDER BY name LIMIT ? OFFSET ?;

-- name: UpdateSubject :exec
UPDATE subjects SET name = ?, department_id = ? WHERE id = ?;

-- name: DeleteSubject :exec
DELETE FROM subjects WHERE id = ?;

