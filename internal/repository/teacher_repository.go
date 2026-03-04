package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	q "school-exam/internal/sqlc/gen"
)

type TeacherRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewTeacherRepository(db *sql.DB) *TeacherRepository {
	return &TeacherRepository{DB: db, Queries: q.New(db)}
}

func (r *TeacherRepository) Create(ctx context.Context, p q.CreateTeacherParams) error {
	return r.Queries.CreateTeacher(ctx, p)
}

func (r *TeacherRepository) GetByID(ctx context.Context, id string) (q.Teacher, error) {
	return r.Queries.GetTeacherByID(ctx, id)
}

func (r *TeacherRepository) ListByTenant(ctx context.Context, tenantID string, limit, offset int64) ([]q.Teacher, error) {
	return r.Queries.ListTeachersByTenant(ctx, q.ListTeachersByTenantParams{
		TenantID: sql.NullString{String: tenantID, Valid: true},
		Limit:    limit,
		Offset:   offset,
	})
}

func (r *TeacherRepository) Update(ctx context.Context, p q.UpdateTeacherParams) error {
	return r.Queries.UpdateTeacher(ctx, p)
}

func (r *TeacherRepository) Delete(ctx context.Context, id string) error {
	return r.Queries.DeleteTeacher(ctx, id)
}

func (r *TeacherRepository) Assign(ctx context.Context, teacherID, subjectID, sectionID string) error {
	return r.Queries.AssignTeacherSubjectSection(ctx, q.AssignTeacherSubjectSectionParams{
		ID:        uuid.New().String(),
		TeacherID: sql.NullString{String: teacherID, Valid: true},
		SubjectID: sql.NullString{String: subjectID, Valid: true},
		SectionID: sql.NullString{String: sectionID, Valid: true},
	})
}

func (r *TeacherRepository) Unassign(ctx context.Context, teacherID, subjectID, sectionID string) error {
	return r.Queries.UnassignTeacherSubjectSection(ctx, q.UnassignTeacherSubjectSectionParams{
		TeacherID: sql.NullString{String: teacherID, Valid: true},
		SubjectID: sql.NullString{String: subjectID, Valid: true},
		SectionID: sql.NullString{String: sectionID, Valid: true},
	})
}

func (r *TeacherRepository) ListAssignments(ctx context.Context, teacherID string, limit, offset int64) ([]q.TeacherSubject, error) {
	return r.Queries.ListAssignmentsByTeacher(ctx, q.ListAssignmentsByTeacherParams{
		TeacherID: sql.NullString{String: teacherID, Valid: true},
		Limit:     limit,
		Offset:    offset,
	})
}

