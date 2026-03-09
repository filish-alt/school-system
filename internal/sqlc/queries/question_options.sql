-- name: CreateOption :exec
INSERT INTO question_options (id, question_id, option_text, is_correct) VALUES (?, ?, ?, ?);

-- name: GetOption :one
SELECT id, question_id, option_text, is_correct FROM question_options WHERE id = ? LIMIT 1;

-- name: ListOptionsByQuestion :many
SELECT id, question_id, option_text, is_correct FROM question_options WHERE question_id = ? LIMIT ? OFFSET ?;

-- name: UpdateOption :exec
UPDATE question_options SET option_text = ?, is_correct = ? WHERE id = ?;

-- name: DeleteOption :exec
DELETE FROM question_options WHERE id = ?;

-- name: ResetCorrectOptions :exec
UPDATE question_options SET is_correct = 0 WHERE question_id = ? AND id != ?;
