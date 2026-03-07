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

func (r *ExamRepository) ListBySection(ctx context.Context, sectionID string, limit, offset int64) ([]q.Exam, error) {
	return r.Queries.ListExamsBySection(ctx, q.ListExamsBySectionParams{
		SectionID: sql.NullString{String: sectionID, Valid: true},
		Limit:     limit,
		Offset:    offset,
	})
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
