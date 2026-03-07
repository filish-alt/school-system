-- name: CreateStudent :exec
INSERT INTO students (id, tenant_id, student_code, first_name, last_name, year, section_id, department_id, user_id, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'active');

-- name: GetStudentByID :one
SELECT id, tenant_id, student_code, first_name, last_name, year, section_id, department_id, user_id, status
FROM students WHERE id = ? LIMIT 1;

-- name: ListByTenant :many
SELECT id, tenant_id, student_code, first_name, last_name, year, section_id, department_id, user_id, status
FROM students WHERE tenant_id = ? ORDER BY last_name, first_name LIMIT ? OFFSET ?;

-- name: UpdateStudent :exec
UPDATE students SET first_name = ?, last_name = ?, year = ?, section_id = ?, department_id = ? WHERE id = ?;

-- name: SetStudentStatus :exec
UPDATE students SET status = ? WHERE id = ?;
