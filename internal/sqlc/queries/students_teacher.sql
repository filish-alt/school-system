-- name: ListStudentsByTeacher :many
SELECT s.id, s.tenant_id, s.student_code, s.first_name, s.last_name, s.year, s.section_id, s.department_id, s.user_id, s.status
FROM students s
JOIN teacher_subjects ts ON ts.section_id = s.section_id
WHERE ts.teacher_id = ? ORDER BY s.last_name, s.first_name LIMIT ? OFFSET ?;

