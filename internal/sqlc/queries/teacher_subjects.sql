-- name: AssignTeacherSubjectSection :exec
INSERT INTO teacher_subjects (id, teacher_id, subject_id, section_id) VALUES (?, ?, ?, ?);

-- name: UnassignTeacherSubjectSection :exec
DELETE FROM teacher_subjects WHERE teacher_id = ? AND subject_id = ? AND section_id = ?;

-- name: ListAssignmentsByTeacher :many
SELECT id, teacher_id, subject_id, section_id FROM teacher_subjects WHERE teacher_id = ? LIMIT ? OFFSET ?;

