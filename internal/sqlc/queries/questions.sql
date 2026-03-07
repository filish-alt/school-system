-- name: CreateQuestion :exec
INSERT INTO questions (id, question_bank_id, type, question_text, marks, difficulty_level) VALUES (?, ?, ?, ?, ?, ?);

-- name: GetQuestion :one
SELECT id, question_bank_id, type, question_text, marks, difficulty_level FROM questions WHERE id = ? LIMIT 1;

-- name: ListQuestionsByBank :many
SELECT id, question_bank_id, type, question_text, marks, difficulty_level FROM questions WHERE question_bank_id = ? LIMIT ? OFFSET ?;

-- name: UpdateQuestion :exec
UPDATE questions SET type = ?, question_text = ?, marks = ?, difficulty_level = ? WHERE id = ?;

-- name: DeleteQuestion :exec
DELETE FROM questions WHERE id = ?;

