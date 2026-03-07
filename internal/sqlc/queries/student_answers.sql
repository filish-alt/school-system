-- name: UpsertStudentAnswer :exec
INSERT INTO student_answers (id, session_id, question_id, answer_text, selected_option_id, score)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(session_id, question_id) DO UPDATE SET
    answer_text = excluded.answer_text,
    selected_option_id = excluded.selected_option_id,
    score = excluded.score;

-- name: GetStudentAnswer :one
SELECT id, session_id, question_id, answer_text, selected_option_id, score
FROM student_answers WHERE session_id = ? AND question_id = ? LIMIT 1;

-- name: GetStudentAnswersBySession :many
SELECT id, session_id, question_id, answer_text, selected_option_id, score
FROM student_answers WHERE session_id = ?;

-- name: GetCorrectOptionForQuestion :one
SELECT id FROM question_options WHERE question_id = ? AND is_correct = 1 LIMIT 1;

-- name: GetExamQuestionMarks :one
SELECT marks FROM exam_questions WHERE exam_id = ? AND question_id = ? LIMIT 1;
