-- name: CreateSection :exec
INSERT INTO sections (id, tenant_id, name, department_id, grade_level, academic_year) VALUES (?, ?, ?, ?, ?, ?);

-- name: GetSectionByID :one
SELECT id, tenant_id, name, department_id, grade_level, academic_year FROM sections WHERE id = ? LIMIT 1;

-- name: ListSectionsByTenant :many
SELECT id, tenant_id, name, department_id, grade_level, academic_year FROM sections WHERE tenant_id = ? ORDER BY name LIMIT ? OFFSET ?;

-- name: UpdateSection :exec
UPDATE sections SET name = ?, department_id = ?, grade_level = ?, academic_year = ? WHERE id = ?;

-- name: DeleteSection :exec
DELETE FROM sections WHERE id = ?;

