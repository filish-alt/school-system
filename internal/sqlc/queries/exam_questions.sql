-- name: AddExamQuestion :exec
INSERT INTO exam_questions (id, exam_id, question_id, marks, order_index) VALUES (?, ?, ?, ?, ?);

-- name: RemoveExamQuestion :exec
DELETE FROM exam_questions WHERE id = ?;

-- name: ListExamQuestions :many
SELECT eq.id, eq.exam_id, eq.question_id, eq.marks, eq.order_index,
       q.type, q.question_text, q.difficulty_level
FROM exam_questions eq
JOIN questions q ON q.id = eq.question_id
WHERE eq.exam_id = ?
ORDER BY eq.order_index ASC;

-- name: GetRandomQuestionsFromBank :many
SELECT id, question_bank_id, type, question_text, marks, difficulty_level
FROM questions
WHERE question_bank_id = ?
ORDER BY RANDOM()
LIMIT ?;
