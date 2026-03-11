-- name: GetExamMarks :many
SELECT 
    s.student_code, 
    s.first_name, 
    s.last_name, 
    sec.name as section_name, 
    es.total_score
FROM exam_sessions es
JOIN students s ON es.student_id = s.id
JOIN sections sec ON s.section_id = sec.id
WHERE es.exam_id = ? AND es.status IN ('submitted', 'timed_out')
ORDER BY s.last_name, s.first_name;
