-- name: CreateQuestionBank :exec
INSERT INTO question_banks (id, tenant_id, subject_id, created_by_teacher_id, title) VALUES (?, ?, ?, ?, ?);

-- name: GetQuestionBank :one
SELECT id, tenant_id, subject_id, created_by_teacher_id, title FROM question_banks WHERE id = ? LIMIT 1;

-- name: ListQuestionBanksByTeacher :many
SELECT id, tenant_id, subject_id, created_by_teacher_id, title
FROM question_banks WHERE created_by_teacher_id = ? ORDER BY title LIMIT ? OFFSET ?;

-- name: DeleteQuestionBank :exec
DELETE FROM question_banks WHERE id = ?;

