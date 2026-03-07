package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type ExamSessionRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewExamSessionRepository(db *sql.DB) *ExamSessionRepository {
	return &ExamSessionRepository{DB: db, Queries: q.New(db)}
}

func (r *ExamSessionRepository) Create(ctx context.Context, p q.CreateExamSessionParams) error {
	return r.Queries.CreateExamSession(ctx, p)
}

func (r *ExamSessionRepository) Get(ctx context.Context, id string) (q.ExamSession, error) {
	return r.Queries.GetExamSession(ctx, id)
}

func (r *ExamSessionRepository) GetActive(ctx context.Context, studentID, examID string) (q.ExamSession, error) {
	return r.Queries.GetActiveSessionByStudent(ctx, q.GetActiveSessionByStudentParams{
		StudentID: sql.NullString{String: studentID, Valid: true},
		ExamID:    sql.NullString{String: examID, Valid: true},
	})
}

func (r *ExamSessionRepository) UpdateStatus(ctx context.Context, id, status string) error {
	return r.Queries.UpdateExamSessionStatus(ctx, q.UpdateExamSessionStatusParams{
		ID:     id,
		Status: sql.NullString{String: status, Valid: true},
	})
}

func (r *ExamSessionRepository) UpdateScore(ctx context.Context, p q.UpdateExamSessionScoreParams) error {
	return r.Queries.UpdateExamSessionScore(ctx, p)
}

func (r *ExamSessionRepository) ListByStudent(ctx context.Context, studentID string) ([]q.ExamSession, error) {
	return r.Queries.ListStudentSessions(ctx, sql.NullString{String: studentID, Valid: true})
}
