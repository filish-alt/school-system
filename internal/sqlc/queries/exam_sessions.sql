-- name: CreateExamSession :exec
INSERT INTO exam_sessions (id, exam_id, student_id, start_time, end_time, status, total_score)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetExamSession :one
SELECT id, exam_id, student_id, start_time, end_time, status, total_score
FROM exam_sessions WHERE id = ? LIMIT 1;

-- name: GetActiveSessionByStudent :one
SELECT id, exam_id, student_id, start_time, end_time, status, total_score
FROM exam_sessions 
WHERE student_id = ? AND exam_id = ? AND status = 'in_progress' 
LIMIT 1;

-- name: UpdateExamSessionStatus :exec
UPDATE exam_sessions SET status = ? WHERE id = ?;

-- name: UpdateExamSessionScore :exec
UPDATE exam_sessions SET total_score = ?, status = 'submitted', end_time = ? WHERE id = ?;

-- name: ListStudentSessions :many
SELECT id, exam_id, student_id, start_time, end_time, status, total_score
FROM exam_sessions WHERE student_id = ? ORDER BY start_time DESC;
