-- name: ListTeacherAssignmentsByTenant :many
SELECT ts.id,
       t.id   AS teacher_id,
       t.first_name,
       t.last_name,
       sub.id AS subject_id,
       sub.name AS subject_name,
       sec.id AS section_id,
       sec.name AS section_name
FROM teacher_subjects ts
JOIN teachers t ON t.id = ts.teacher_id
JOIN subjects sub ON sub.id = ts.subject_id
JOIN sections sec ON sec.id = ts.section_id
WHERE t.tenant_id = ?
ORDER BY t.last_name, t.first_name, subject_name, section_name
LIMIT ? OFFSET ?;

-- name: ListTeacherAssignmentsByTeacher :many
SELECT ts.id,
       t.id   AS teacher_id,
       t.first_name,
       t.last_name,
       sub.id AS subject_id,
       sub.name AS subject_name,
       sec.id AS section_id,
       sec.name AS section_name
FROM teacher_subjects ts
JOIN teachers t ON t.id = ts.teacher_id
JOIN subjects sub ON sub.id = ts.subject_id
JOIN sections sec ON sec.id = ts.section_id
WHERE ts.teacher_id = ?
ORDER BY subject_name, section_name
LIMIT ? OFFSET ?;

