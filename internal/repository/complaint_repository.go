package repository

import (
	"context"
	"database/sql"
	"school-exam/internal/domain"
	"time"
)

type ComplaintRepository struct {
	DB *sql.DB
}

func NewComplaintRepository(db *sql.DB) *ComplaintRepository {
	return &ComplaintRepository{DB: db}
}

func (r *ComplaintRepository) Create(ctx context.Context, c domain.StudentComplaint) error {
	query := `INSERT INTO student_complaints (id, tenant_id, student_id, subject_id, exam_id, coursework_id, question_id, category, description, status, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, c.ID, c.TenantID, c.StudentID, c.SubjectID, c.ExamID, c.CourseworkID, c.QuestionID, c.Category, c.Description, "pending", time.Now())
	return err
}

func (r *ComplaintRepository) ListByStudent(ctx context.Context, studentID string) ([]domain.StudentComplaint, error) {
	query := `SELECT c.id, c.tenant_id, c.student_id, s.student_code, (s.first_name || ' ' || s.last_name) as student_name,
	                 c.subject_id, sub.name as subject_name, c.exam_id, ex.title as exam_title,
	                 c.coursework_id, c.question_id, c.category, c.description, c.status, c.response,
	                 c.resolved_by, u.username as resolver_name, c.resolved_at, c.created_at
	          FROM student_complaints c
	          JOIN students s ON c.student_id = s.id
	          LEFT JOIN subjects sub ON c.subject_id = sub.id
	          LEFT JOIN exams ex ON c.exam_id = ex.id
	          LEFT JOIN users u ON c.resolved_by = u.id
	          WHERE c.student_id = ?
	          ORDER BY c.created_at DESC`
	rows, err := r.DB.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.StudentComplaint
	for rows.Next() {
		var c domain.StudentComplaint
		var stCode, stName, subID, subName, exID, exTitle, cwID, qID, resp, resBy, resName sql.NullString
		var resAt sql.NullTime
		if err := rows.Scan(&c.ID, &c.TenantID, &c.StudentID, &stCode, &stName,
			&subID, &subName, &exID, &exTitle,
			&cwID, &qID, &c.Category, &c.Description, &c.Status, &resp,
			&resBy, &resName, &resAt, &c.CreatedAt); err != nil {
			return nil, err
		}
		if stCode.Valid { c.StudentCode = &stCode.String }
		if stName.Valid { c.StudentName = &stName.String }
		if subID.Valid { c.SubjectID = &subID.String }
		if subName.Valid { c.SubjectName = &subName.String }
		if exID.Valid { c.ExamID = &exID.String }
		if exTitle.Valid { c.ExamTitle = &exTitle.String }
		if cwID.Valid { c.CourseworkID = &cwID.String }
		if qID.Valid { c.QuestionID = &qID.String }
		if resp.Valid { c.Response = &resp.String }
		if resBy.Valid { c.ResolvedBy = &resBy.String }
		if resName.Valid { c.ResolverName = &resName.String }
		if resAt.Valid { c.ResolvedAt = &resAt.Time }
		list = append(list, c)
	}
	return list, nil
}

func (r *ComplaintRepository) ListByTenant(ctx context.Context, tenantID string, statusFilter string) ([]domain.StudentComplaint, error) {
	baseQuery := `SELECT c.id, c.tenant_id, c.student_id, s.student_code, (s.first_name || ' ' || s.last_name) as student_name,
	                     c.subject_id, sub.name as subject_name, c.exam_id, ex.title as exam_title,
	                     c.coursework_id, c.question_id, c.category, c.description, c.status, c.response,
	                     c.resolved_by, u.username as resolver_name, c.resolved_at, c.created_at
	              FROM student_complaints c
	              JOIN students s ON c.student_id = s.id
	              LEFT JOIN subjects sub ON c.subject_id = sub.id
	              LEFT JOIN exams ex ON c.exam_id = ex.id
	              LEFT JOIN users u ON c.resolved_by = u.id
	              WHERE c.tenant_id = ?`

	var rows *sql.Rows
	var err error
	if statusFilter != "" {
		rows, err = r.DB.QueryContext(ctx, baseQuery+" AND c.status = ? ORDER BY c.created_at DESC", tenantID, statusFilter)
	} else {
		rows, err = r.DB.QueryContext(ctx, baseQuery+" ORDER BY c.created_at DESC", tenantID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.StudentComplaint
	for rows.Next() {
		var c domain.StudentComplaint
		var stCode, stName, subID, subName, exID, exTitle, cwID, qID, resp, resBy, resName sql.NullString
		var resAt sql.NullTime
		if err := rows.Scan(&c.ID, &c.TenantID, &c.StudentID, &stCode, &stName,
			&subID, &subName, &exID, &exTitle,
			&cwID, &qID, &c.Category, &c.Description, &c.Status, &resp,
			&resBy, &resName, &resAt, &c.CreatedAt); err != nil {
			return nil, err
		}
		if stCode.Valid { c.StudentCode = &stCode.String }
		if stName.Valid { c.StudentName = &stName.String }
		if subID.Valid { c.SubjectID = &subID.String }
		if subName.Valid { c.SubjectName = &subName.String }
		if exID.Valid { c.ExamID = &exID.String }
		if exTitle.Valid { c.ExamTitle = &exTitle.String }
		if cwID.Valid { c.CourseworkID = &cwID.String }
		if qID.Valid { c.QuestionID = &qID.String }
		if resp.Valid { c.Response = &resp.String }
		if resBy.Valid { c.ResolvedBy = &resBy.String }
		if resName.Valid { c.ResolverName = &resName.String }
		if resAt.Valid { c.ResolvedAt = &resAt.Time }
		list = append(list, c)
	}
	return list, nil
}

func (r *ComplaintRepository) UpdateStatus(ctx context.Context, id, status, response, resolvedBy string) error {
	query := `UPDATE student_complaints SET
	          status = ?,
	          response = ?,
	          resolved_by = ?,
	          resolved_at = CURRENT_TIMESTAMP
	          WHERE id = ?`
	_, err := r.DB.ExecContext(ctx, query, status, response, resolvedBy, id)
	return err
}
