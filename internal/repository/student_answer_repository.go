package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type StudentAnswerRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewStudentAnswerRepository(db *sql.DB) *StudentAnswerRepository {
	return &StudentAnswerRepository{DB: db, Queries: q.New(db)}
}

func (r *StudentAnswerRepository) Upsert(ctx context.Context, p q.UpsertStudentAnswerParams) error {
	return r.Queries.UpsertStudentAnswer(ctx, p)
}

func (r *StudentAnswerRepository) Get(ctx context.Context, sessionID, questionID string) (q.StudentAnswer, error) {
	return r.Queries.GetStudentAnswer(ctx, q.GetStudentAnswerParams{
		SessionID:  sql.NullString{String: sessionID, Valid: true},
		QuestionID: sql.NullString{String: questionID, Valid: true},
	})
}

func (r *StudentAnswerRepository) ListBySession(ctx context.Context, sessionID string) ([]q.StudentAnswer, error) {
	return r.Queries.GetStudentAnswersBySession(ctx, sql.NullString{String: sessionID, Valid: true})
}

func (r *StudentAnswerRepository) GetCorrectOption(ctx context.Context, questionID string) (string, error) {
	return r.Queries.GetCorrectOptionForQuestion(ctx, sql.NullString{String: questionID, Valid: true})
}

func (r *StudentAnswerRepository) GetQuestionMarks(ctx context.Context, examID, questionID string) (int64, error) {
	val, err := r.Queries.GetExamQuestionMarks(ctx, q.GetExamQuestionMarksParams{
		ExamID:     sql.NullString{String: examID, Valid: true},
		QuestionID: sql.NullString{String: questionID, Valid: true},
	})
	return val.Int64, err
}
