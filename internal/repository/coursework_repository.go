package repository

import (
	"context"
	"database/sql"
	"school-exam/internal/domain"
	"time"
)

type CourseworkRepository struct {
	DB *sql.DB
}

func NewCourseworkRepository(db *sql.DB) *CourseworkRepository {
	return &CourseworkRepository{DB: db}
}

func (r *CourseworkRepository) Create(ctx context.Context, c domain.CourseworkAssessment) error {
	// 1. Resolve actual teacher.id from user_id if needed
	var actualTeacherID string
	var teacherTenant sql.NullString
	err := r.DB.QueryRowContext(ctx, "SELECT id, tenant_id FROM teachers WHERE user_id = ? OR id = ? LIMIT 1", c.TeacherID, c.TeacherID).Scan(&actualTeacherID, &teacherTenant)
	if err == nil && actualTeacherID != "" {
		c.TeacherID = actualTeacherID
		if c.TenantID == "" && teacherTenant.Valid {
			c.TenantID = teacherTenant.String
		}
	}

	// 2. Resolve tenant_id from section or subject if still empty
	if c.TenantID == "" {
		var secTenant sql.NullString
		_ = r.DB.QueryRowContext(ctx, "SELECT tenant_id FROM sections WHERE id = ?", c.SectionID).Scan(&secTenant)
		if secTenant.Valid {
			c.TenantID = secTenant.String
		}
	}

	// 3. Resolve academic_year_id (must be valid academic_years.id or NULL)
	if c.AcademicYearID != nil && *c.AcademicYearID == "" {
		c.AcademicYearID = nil
	}
	if c.AcademicYearID == nil {
		var activeAY sql.NullString
		_ = r.DB.QueryRowContext(ctx, "SELECT id FROM academic_years WHERE is_active = 1 LIMIT 1").Scan(&activeAY)
		if activeAY.Valid {
			c.AcademicYearID = &activeAY.String
		}
	}

	query := `INSERT INTO coursework_assessments (id, tenant_id, academic_year_id, subject_id, section_id, teacher_id, title, type, max_marks, weight_percentage, due_date, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = r.DB.ExecContext(ctx, query, c.ID, c.TenantID, c.AcademicYearID, c.SubjectID, c.SectionID, c.TeacherID, c.Title, c.Type, c.MaxMarks, c.WeightPercentage, c.DueDate, time.Now())
	return err
}

func (r *CourseworkRepository) ListBySectionSubject(ctx context.Context, sectionID, subjectID string) ([]domain.CourseworkAssessment, error) {
	query := `SELECT c.id, c.tenant_id, c.academic_year_id, c.subject_id, sub.name as subject_name,
	                 c.section_id, sec.name as section_name, c.teacher_id, (t.first_name || ' ' || t.last_name) as teacher_name,
	                 c.title, c.type, c.max_marks, c.weight_percentage, c.due_date, c.created_at
	          FROM coursework_assessments c
	          LEFT JOIN subjects sub ON c.subject_id = sub.id
	          LEFT JOIN sections sec ON c.section_id = sec.id
	          LEFT JOIN teachers t ON c.teacher_id = t.id
	          WHERE c.section_id = ? AND c.subject_id = ?
	          ORDER BY c.created_at DESC`
	rows, err := r.DB.QueryContext(ctx, query, sectionID, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.CourseworkAssessment
	for rows.Next() {
		var c domain.CourseworkAssessment
		var ayID, subName, secName, teaName, dueDate sql.NullString
		if err := rows.Scan(&c.ID, &c.TenantID, &ayID, &c.SubjectID, &subName,
			&c.SectionID, &secName, &c.TeacherID, &teaName,
			&c.Title, &c.Type, &c.MaxMarks, &c.WeightPercentage, &dueDate, &c.CreatedAt); err != nil {
			return nil, err
		}
		if ayID.Valid { c.AcademicYearID = &ayID.String }
		if subName.Valid { c.SubjectName = &subName.String }
		if secName.Valid { c.SectionName = &secName.String }
		if teaName.Valid { c.TeacherName = &teaName.String }
		if dueDate.Valid { c.DueDate = &dueDate.String }
		list = append(list, c)
	}
	return list, nil
}

func (r *CourseworkRepository) ListByTeacher(ctx context.Context, teacherID string) ([]domain.CourseworkAssessment, error) {
	query := `SELECT c.id, c.tenant_id, c.academic_year_id, c.subject_id, sub.name as subject_name,
	                 c.section_id, sec.name as section_name, c.teacher_id, (t.first_name || ' ' || t.last_name) as teacher_name,
	                 c.title, c.type, c.max_marks, c.weight_percentage, c.due_date, c.created_at
	          FROM coursework_assessments c
	          LEFT JOIN subjects sub ON c.subject_id = sub.id
	          LEFT JOIN sections sec ON c.section_id = sec.id
	          LEFT JOIN teachers t ON c.teacher_id = t.id
	          WHERE c.teacher_id = ? OR t.user_id = ? OR c.teacher_id = (SELECT id FROM teachers WHERE user_id = ?)
	          ORDER BY c.created_at DESC`
	rows, err := r.DB.QueryContext(ctx, query, teacherID, teacherID, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.CourseworkAssessment
	for rows.Next() {
		var c domain.CourseworkAssessment
		var ayID, subName, secName, teaName, dueDate sql.NullString
		if err := rows.Scan(&c.ID, &c.TenantID, &ayID, &c.SubjectID, &subName,
			&c.SectionID, &secName, &c.TeacherID, &teaName,
			&c.Title, &c.Type, &c.MaxMarks, &c.WeightPercentage, &dueDate, &c.CreatedAt); err != nil {
			return nil, err
		}
		if ayID.Valid { c.AcademicYearID = &ayID.String }
		if subName.Valid { c.SubjectName = &subName.String }
		if secName.Valid { c.SectionName = &secName.String }
		if teaName.Valid { c.TeacherName = &teaName.String }
		if dueDate.Valid { c.DueDate = &dueDate.String }
		list = append(list, c)
	}
	return list, nil
}

func (r *CourseworkRepository) UpsertMark(ctx context.Context, m domain.StudentCourseworkMark) error {
	var actualStudentID string
	err := r.DB.QueryRowContext(ctx, "SELECT id FROM students WHERE id = ? OR user_id = ? LIMIT 1", m.StudentID, m.StudentID).Scan(&actualStudentID)
	if err == nil && actualStudentID != "" {
		m.StudentID = actualStudentID
	}

	query := `INSERT INTO student_coursework_marks (id, coursework_id, student_id, marks_obtained, remarks, graded_by, graded_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)
	          ON CONFLICT(coursework_id, student_id) DO UPDATE SET
	          marks_obtained = excluded.marks_obtained,
	          remarks = excluded.remarks,
	          graded_by = excluded.graded_by,
	          graded_at = excluded.graded_at`
	_, err = r.DB.ExecContext(ctx, query, m.ID, m.CourseworkID, m.StudentID, m.MarksObtained, m.Remarks, m.GradedBy, time.Now())
	return err
}

func (r *CourseworkRepository) ListMarksByCoursework(ctx context.Context, courseworkID string) ([]domain.StudentCourseworkMark, error) {
	query := `SELECT m.id, m.coursework_id, m.student_id, s.student_code, (s.first_name || ' ' || s.last_name) as student_name,
	                 m.marks_obtained, m.remarks, m.graded_by, m.graded_at
	          FROM student_coursework_marks m
	          JOIN students s ON m.student_id = s.id
	          WHERE m.coursework_id = ?
	          ORDER BY s.student_code ASC`
	rows, err := r.DB.QueryContext(ctx, query, courseworkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.StudentCourseworkMark
	for rows.Next() {
		var m domain.StudentCourseworkMark
		var code, name, remarks, gradedBy sql.NullString
		if err := rows.Scan(&m.ID, &m.CourseworkID, &m.StudentID, &code, &name, &m.MarksObtained, &remarks, &gradedBy, &m.GradedAt); err != nil {
			return nil, err
		}
		if code.Valid { m.StudentCode = &code.String }
		if name.Valid { m.StudentName = &name.String }
		if remarks.Valid { m.Remarks = &remarks.String }
		if gradedBy.Valid { m.GradedBy = &gradedBy.String }
		list = append(list, m)
	}
	return list, nil
}

func (r *CourseworkRepository) ListMarksByStudent(ctx context.Context, studentID string) ([]domain.StudentCourseworkMark, error) {
	query := `SELECT m.id, m.coursework_id, m.student_id, s.student_code, (s.first_name || ' ' || s.last_name) as student_name,
	                 m.marks_obtained, m.remarks, m.graded_by, m.graded_at
	          FROM student_coursework_marks m
	          JOIN students s ON m.student_id = s.id
	          WHERE m.student_id = ? OR s.user_id = ? OR m.student_id = (SELECT id FROM students WHERE user_id = ?)
	          ORDER BY m.graded_at DESC`
	rows, err := r.DB.QueryContext(ctx, query, studentID, studentID, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.StudentCourseworkMark
	for rows.Next() {
		var m domain.StudentCourseworkMark
		var code, name, remarks, gradedBy sql.NullString
		if err := rows.Scan(&m.ID, &m.CourseworkID, &m.StudentID, &code, &name, &m.MarksObtained, &remarks, &gradedBy, &m.GradedAt); err != nil {
			return nil, err
		}
		if code.Valid { m.StudentCode = &code.String }
		if name.Valid { m.StudentName = &name.String }
		if remarks.Valid { m.Remarks = &remarks.String }
		if gradedBy.Valid { m.GradedBy = &gradedBy.String }
		list = append(list, m)
	}
	return list, nil
}

