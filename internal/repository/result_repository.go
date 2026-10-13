package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"school-exam/internal/domain"
	"time"
)

type ResultRepository struct {
	DB *sql.DB
}

func NewResultRepository(db *sql.DB) *ResultRepository {
	return &ResultRepository{DB: db}
}

func (r *ResultRepository) CalculateAndSaveResults(ctx context.Context, tenantID, sectionID, subjectID string, academicYearID *string) error {
	// 1. Fetch all students in section
	queryStudents := `SELECT id, student_code, (first_name || ' ' || last_name) as name FROM students WHERE section_id = ?`
	rows, err := r.DB.QueryContext(ctx, queryStudents, sectionID)
	if err != nil {
		return err
	}
	defer rows.Close()

	type studentInfo struct {
		id, code, name string
	}
	var students []studentInfo
	for rows.Next() {
		var s studentInfo
		if err := rows.Scan(&s.id, &s.code, &s.name); err == nil {
			students = append(students, s)
		}
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Fetch configured assessment weights and grading thresholds if available
	var finalExamWeight float64 = 60.0
	var passThresh, goodThresh, excThresh float64 = 50.0, 70.0, 85.0

	var settingsJSON string
	if err := tx.QueryRowContext(ctx, "SELECT setting_value FROM system_settings WHERE tenant_id = ? AND setting_key = 'assessment_weights'", tenantID).Scan(&settingsJSON); err == nil {
		// e.g. {"final": 60, "quiz": 10, ...}
		// if final weight is defined in settings, use it
		var parsed map[string]interface{}
		if jsonErr := json.Unmarshal([]byte(settingsJSON), &parsed); jsonErr == nil {
			if f, ok := parsed["final"].(float64); ok && f > 0 {
				finalExamWeight = f
			}
		}
	}

	var threshJSON string
	if err := tx.QueryRowContext(ctx, "SELECT setting_value FROM system_settings WHERE tenant_id = ? AND setting_key = 'performance_thresholds'", tenantID).Scan(&threshJSON); err == nil {
		var parsed map[string]interface{}
		if jsonErr := json.Unmarshal([]byte(threshJSON), &parsed); jsonErr == nil {
			if p, ok := parsed["pass"].(float64); ok && p > 0 { passThresh = p }
			if g, ok := parsed["good"].(float64); ok && g > 0 { goodThresh = g }
			if e, ok := parsed["excellent"].(float64); ok && e > 0 { excThresh = e }
		}
	}

	for _, s := range students {
		// 1. Calculate Quiz score (dynamically scaled by coursework weight_percentage)
		var quizScore float64
		_ = tx.QueryRowContext(ctx, `SELECT COALESCE(SUM(m.marks_obtained * c.weight_percentage / c.max_marks), 0)
		                             FROM student_coursework_marks m
		                             JOIN coursework_assessments c ON m.coursework_id = c.id
		                             WHERE m.student_id = ? AND c.subject_id = ? AND c.section_id = ? AND c.type = 'quiz'`,
			s.id, subjectID, sectionID).Scan(&quizScore)

		// 2. Calculate Mid score (dynamically scaled by coursework weight_percentage)
		var midScore float64
		_ = tx.QueryRowContext(ctx, `SELECT COALESCE(SUM(m.marks_obtained * c.weight_percentage / c.max_marks), 0)
		                             FROM student_coursework_marks m
		                             JOIN coursework_assessments c ON m.coursework_id = c.id
		                             WHERE m.student_id = ? AND c.subject_id = ? AND c.section_id = ? AND c.type = 'mid_exam'`,
			s.id, subjectID, sectionID).Scan(&midScore)

		// 3. Calculate Assignment score (dynamically scaled by coursework weight_percentage)
		var assignScore float64
		_ = tx.QueryRowContext(ctx, `SELECT COALESCE(SUM(m.marks_obtained * c.weight_percentage / c.max_marks), 0)
		                             FROM student_coursework_marks m
		                             JOIN coursework_assessments c ON m.coursework_id = c.id
		                             WHERE m.student_id = ? AND c.subject_id = ? AND c.section_id = ? AND c.type = 'assignment'`,
			s.id, subjectID, sectionID).Scan(&assignScore)

		// 4. Calculate Project score (dynamically scaled by coursework weight_percentage)
		var projScore float64
		_ = tx.QueryRowContext(ctx, `SELECT COALESCE(SUM(m.marks_obtained * c.weight_percentage / c.max_marks), 0)
		                             FROM student_coursework_marks m
		                             JOIN coursework_assessments c ON m.coursework_id = c.id
		                             WHERE m.student_id = ? AND c.subject_id = ? AND c.section_id = ? AND c.type = 'project'`,
			s.id, subjectID, sectionID).Scan(&projScore)

		// 5. Calculate Final Exam score (dynamically scaled to configured finalExamWeight)
		var finalScore float64
		_ = tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(es.total_score * ? / ex.total_marks), 0)
		                             FROM exam_sessions es
		                             JOIN exams ex ON es.exam_id = ex.id
		                             WHERE es.student_id = ? AND ex.subject_id = ? AND ex.section_id = ? AND (ex.exam_type = 'final' OR ex.title LIKE '%Final%')`,
			finalExamWeight, s.id, subjectID, sectionID).Scan(&finalScore)

		caScore := quizScore + midScore + assignScore + projScore
		totalScore := caScore + finalScore

		grade := "F"
		if totalScore >= excThresh {
			grade = "A"
		} else if totalScore >= goodThresh {
			grade = "B"
		} else if totalScore >= passThresh {
			grade = "C"
		} else if totalScore >= (passThresh - 10) {
			grade = "D"
		}

		resID := s.id + "_" + subjectID
		queryUpsert := `INSERT INTO results (id, tenant_id, academic_year_id, student_id, subject_id, section_id, quiz_score, mid_score, assignment_score, project_score, ca_score, final_score, total_score, grade, is_published, published_at)
		                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?)
		                ON CONFLICT(id) DO UPDATE SET
		                tenant_id = excluded.tenant_id,
		                academic_year_id = excluded.academic_year_id,
		                section_id = excluded.section_id,
		                quiz_score = excluded.quiz_score,
		                mid_score = excluded.mid_score,
		                assignment_score = excluded.assignment_score,
		                project_score = excluded.project_score,
		                ca_score = excluded.ca_score,
		                final_score = excluded.final_score,
		                total_score = excluded.total_score,
		                grade = excluded.grade,
		                is_published = 1,
		                published_at = excluded.published_at`
		if _, err := tx.ExecContext(ctx, queryUpsert, resID, tenantID, academicYearID, s.id, subjectID, sectionID, quizScore, midScore, assignScore, projScore, caScore, finalScore, totalScore, grade, time.Now()); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ResultRepository) GetTransferSheet(ctx context.Context, sectionID string) ([]domain.AssessmentResult, error) {
	list := make([]domain.AssessmentResult, 0)
	query := `SELECT r.id, r.tenant_id, r.academic_year_id, r.student_id, 
	                 COALESCE(s.student_code, ''), 
	                 COALESCE((s.first_name || ' ' || s.last_name), ''),
	                 r.subject_id, 
	                 COALESCE(sub.name, ''), 
	                 r.section_id, 
	                 COALESCE(sec.name, ''),
	                 COALESCE(r.quiz_score, 0), 
	                 COALESCE(r.mid_score, 0), 
	                 COALESCE(r.assignment_score, 0), 
	                 COALESCE(r.project_score, 0), 
	                 COALESCE(r.ca_score, 0), 
	                 COALESCE(r.final_score, 0), 
	                 COALESCE(r.total_score, 0), 
	                 COALESCE(r.grade, 'F'), 
	                 COALESCE(r.is_published, 0), 
	                 r.published_at
	          FROM results r
	          LEFT JOIN students s ON r.student_id = s.id
	          LEFT JOIN subjects sub ON r.subject_id = sub.id
	          LEFT JOIN sections sec ON COALESCE(r.section_id, s.section_id) = sec.id
	          WHERE r.section_id = ? OR s.section_id = ?
	          ORDER BY COALESCE(s.student_code, ''), COALESCE(sub.name, '') ASC`
	rows, err := r.DB.QueryContext(ctx, query, sectionID, sectionID)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var res domain.AssessmentResult
		var tenID, ayID, secID, secName sql.NullString
		var stdCode, stdName, subName, grade sql.NullString
		var pubAt sql.NullTime
		var isPubInt int
		if err := rows.Scan(&res.ID, &tenID, &ayID, &res.StudentID, 
			&stdCode, &stdName,
			&res.SubjectID, &subName, 
			&secID, &secName,
			&res.QuizScore, &res.MidScore, &res.AssignmentScore, &res.ProjectScore,
			&res.CAScore, &res.FinalScore, &res.TotalScore, &grade, &isPubInt, &pubAt); err != nil {
			return list, err
		}
		if tenID.Valid { res.TenantID = &tenID.String }
		if ayID.Valid { res.AcademicYearID = &ayID.String }
		if secID.Valid { res.SectionID = &secID.String }
		if secName.Valid { res.SectionName = &secName.String }
		if stdCode.Valid { res.StudentCode = stdCode.String }
		if stdName.Valid { res.StudentName = stdName.String }
		if subName.Valid { res.SubjectName = subName.String }
		if grade.Valid { res.Grade = grade.String }
		if pubAt.Valid { res.PublishedAt = &pubAt.Time }
		res.IsPublished = isPubInt == 1
		list = append(list, res)
	}
	return list, nil
}

func (r *ResultRepository) GetStudentResults(ctx context.Context, studentID string) ([]domain.AssessmentResult, error) {
	list := make([]domain.AssessmentResult, 0)
	query := `SELECT r.id, r.tenant_id, r.academic_year_id, r.student_id, 
	                 COALESCE(s.student_code, ''), 
	                 COALESCE((s.first_name || ' ' || s.last_name), ''),
	                 r.subject_id, 
	                 COALESCE(sub.name, ''), 
	                 r.section_id, 
	                 COALESCE(sec.name, ''),
	                 COALESCE(r.quiz_score, 0), 
	                 COALESCE(r.mid_score, 0), 
	                 COALESCE(r.assignment_score, 0), 
	                 COALESCE(r.project_score, 0), 
	                 COALESCE(r.ca_score, 0), 
	                 COALESCE(r.final_score, 0), 
	                 COALESCE(r.total_score, 0), 
	                 COALESCE(r.grade, 'F'), 
	                 COALESCE(r.is_published, 0), 
	                 r.published_at
	          FROM results r
	          LEFT JOIN students s ON r.student_id = s.id
	          LEFT JOIN subjects sub ON r.subject_id = sub.id
	          LEFT JOIN sections sec ON COALESCE(r.section_id, s.section_id) = sec.id
	          WHERE r.student_id = ? AND r.is_published = 1
	          ORDER BY COALESCE(sub.name, '') ASC`
	rows, err := r.DB.QueryContext(ctx, query, studentID)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var res domain.AssessmentResult
		var tenID, ayID, secID, secName sql.NullString
		var stdCode, stdName, subName, grade sql.NullString
		var pubAt sql.NullTime
		var isPubInt int
		if err := rows.Scan(&res.ID, &tenID, &ayID, &res.StudentID, 
			&stdCode, &stdName,
			&res.SubjectID, &subName, 
			&secID, &secName,
			&res.QuizScore, &res.MidScore, &res.AssignmentScore, &res.ProjectScore,
			&res.CAScore, &res.FinalScore, &res.TotalScore, &grade, &isPubInt, &pubAt); err != nil {
			return list, err
		}
		if tenID.Valid { res.TenantID = &tenID.String }
		if ayID.Valid { res.AcademicYearID = &ayID.String }
		if secID.Valid { res.SectionID = &secID.String }
		if secName.Valid { res.SectionName = &secName.String }
		if stdCode.Valid { res.StudentCode = stdCode.String }
		if stdName.Valid { res.StudentName = stdName.String }
		if subName.Valid { res.SubjectName = subName.String }
		if grade.Valid { res.Grade = grade.String }
		if pubAt.Valid { res.PublishedAt = &pubAt.Time }
		res.IsPublished = isPubInt == 1
		list = append(list, res)
	}
	return list, nil
}
