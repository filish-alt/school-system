-- name: CreateTeacher :exec
INSERT INTO teachers (id, tenant_id, teacher_code, first_name, last_name, department_id, user_id) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetTeacherByID :one
SELECT id, tenant_id, teacher_code, first_name, last_name, department_id, user_id FROM teachers WHERE id = ? LIMIT 1;

-- name: GetTeacherByUserID :one
SELECT id, tenant_id, teacher_code, first_name, last_name, department_id, user_id FROM teachers WHERE user_id = ? LIMIT 1;

-- name: ListTeachersByTenant :many
SELECT id, tenant_id, teacher_code, first_name, last_name, department_id, user_id FROM teachers WHERE tenant_id = ? ORDER BY last_name, first_name LIMIT ? OFFSET ?;

-- name: UpdateTeacher :exec
UPDATE teachers SET first_name = ?, last_name = ?, department_id = ?, teacher_code = ? WHERE id = ?;

-- name: DeleteTeacher :exec
DELETE FROM teachers WHERE id = ?;
