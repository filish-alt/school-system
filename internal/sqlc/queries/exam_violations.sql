-- name: CreateExamViolation :exec
INSERT INTO exam_violations (id, session_id, violation_type, created_at)
VALUES (?, ?, ?, ?);

-- name: ListExamViolationsBySession :many
SELECT id, session_id, violation_type, created_at
FROM exam_violations
WHERE session_id = ?;

-- name: ListAllExamViolations :many
SELECT v.id, v.session_id, v.violation_type, v.created_at, 
       s.student_id, st.first_name, st.last_name, 
       s.exam_id, ex.title as exam_title
FROM exam_violations v
JOIN exam_sessions s ON v.session_id = s.id
JOIN students st ON s.student_id = st.id
JOIN exams ex ON s.exam_id = ex.id
WHERE st.tenant_id = ?;
