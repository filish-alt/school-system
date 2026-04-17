-- name: CreateExam :exec
INSERT INTO exams (id, tenant_id, title, subject_id, section_id, created_by_teacher_id, duration_minutes, start_time, end_time, status, total_marks, shuffle_options)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetExam :one
SELECT id, tenant_id, title, subject_id, section_id, created_by_teacher_id, duration_minutes, start_time, end_time, status, total_marks, shuffle_options
FROM exams WHERE id = ? LIMIT 1;

-- name: ListExamsByTeacher :many
SELECT id, tenant_id, title, subject_id, section_id, created_by_teacher_id, duration_minutes, start_time, end_time, status, total_marks, shuffle_options
FROM exams WHERE created_by_teacher_id = ? ORDER BY rowid DESC LIMIT ? OFFSET ?;

-- name: ListExamsBySection :many
SELECT id, tenant_id, title, subject_id, section_id, created_by_teacher_id, duration_minutes, start_time, end_time, status, total_marks, shuffle_options
FROM exams WHERE section_id = ? ORDER BY rowid DESC LIMIT ? OFFSET ?;

-- name: UpdateExam :exec
UPDATE exams
SET title = ?, subject_id = ?, section_id = ?, duration_minutes = ?, start_time = ?, end_time = ?, shuffle_options = ?
WHERE id = ?;

-- name: UpdateExamStatus :exec
UPDATE exams SET status = ? WHERE id = ?;

-- name: UpdateExamTotalMarks :exec
UPDATE exams SET total_marks = (
    SELECT COALESCE(SUM(COALESCE(NULLIF(eq.marks, 0), q.marks)), 0) 
    FROM exam_questions eq
    JOIN questions q ON q.id = eq.question_id
    WHERE eq.exam_id = ?
) WHERE exams.id = ?;

-- name: DeleteExam :exec
DELETE FROM exams WHERE id = ?;

-- name: ListPublishedExamsBySection :many
SELECT id, tenant_id, title, subject_id, section_id, created_by_teacher_id, duration_minutes, start_time, end_time, status, total_marks, shuffle_options
FROM exams WHERE section_id = ? AND status = 'published' ORDER BY start_time ASC;
