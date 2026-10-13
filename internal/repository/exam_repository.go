package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type ExamRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewExamRepository(db *sql.DB) *ExamRepository {
	return &ExamRepository{DB: db, Queries: q.New(db)}
}

func (r *ExamRepository) Create(ctx context.Context, p q.CreateExamParams) error {
	return r.Queries.CreateExam(ctx, p)
}

func (r *ExamRepository) Get(ctx context.Context, id string) (q.Exam, error) {
	return r.Queries.GetExam(ctx, id)
}

func (r *ExamRepository) ListByTeacher(ctx context.Context, teacherID string, limit, offset int64) ([]q.Exam, error) {
	return r.Queries.ListExamsByTeacher(ctx, q.ListExamsByTeacherParams{
		CreatedByTeacherID: sql.NullString{String: teacherID, Valid: true},
		Limit:              limit,
		Offset:             offset,
	})
}

func (r *ExamRepository) ListByTenant(ctx context.Context, tenantID string, limit, offset int64) ([]q.Exam, error) {
	query := `
		SELECT id, tenant_id, title, subject_id, section_id, created_by_teacher_id, duration_minutes, start_time, end_time, status, total_marks, shuffle_options
		FROM exams
		WHERE tenant_id = ? OR ? = ''
		ORDER BY rowid DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.DB.QueryContext(ctx, query, tenantID, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []q.Exam
	for rows.Next() {
		var i q.Exam
		if err := rows.Scan(
			&i.ID,
			&i.TenantID,
			&i.Title,
			&i.SubjectID,
			&i.SectionID,
			&i.CreatedByTeacherID,
			&i.DurationMinutes,
			&i.StartTime,
			&i.EndTime,
			&i.Status,
			&i.TotalMarks,
			&i.ShuffleOptions,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ExamRepository) ListBySection(ctx context.Context, sectionID string, limit, offset int64) ([]q.Exam, error) {
	return r.Queries.ListExamsBySection(ctx, q.ListExamsBySectionParams{
		SectionID: sectionID,
		Limit:     limit,
		Offset:    offset,
	})
}

func (r *ExamRepository) ListPublishedBySection(ctx context.Context, sectionID string) ([]q.Exam, error) {
	return r.Queries.ListPublishedExamsBySection(ctx, sectionID)
}

func (r *ExamRepository) Update(ctx context.Context, p q.UpdateExamParams) error {
	return r.Queries.UpdateExam(ctx, p)
}

func (r *ExamRepository) UpdateStatus(ctx context.Context, id, status string) error {
	return r.Queries.UpdateExamStatus(ctx, q.UpdateExamStatusParams{
		ID:     id,
		Status: sql.NullString{String: status, Valid: true},
	})
}

func (r *ExamRepository) UpdateTotalMarks(ctx context.Context, id string) error {
	return r.Queries.UpdateExamTotalMarks(ctx, q.UpdateExamTotalMarksParams{
		ExamID:   sql.NullString{String: id, Valid: true},
		ID: id,
	})
}

func (r *ExamRepository) Delete(ctx context.Context, id string) error {
	return r.Queries.DeleteExam(ctx, id)
}

func (r *ExamRepository) GetExamMarks(ctx context.Context, examID string) ([]q.GetExamMarksRow, error) {
	return r.Queries.GetExamMarks(ctx, sql.NullString{String: examID, Valid: true})
}
