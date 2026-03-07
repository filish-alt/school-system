package repository

import (
	"context"
	"database/sql"
	q "school-exam/internal/sqlc/gen"
)

type StudentRepository struct {
	DB      *sql.DB
	Queries *q.Queries
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{DB: db, Queries: q.New(db)}
}

func (r *StudentRepository) Create(ctx context.Context, params q.CreateStudentParams) error {
	return r.Queries.CreateStudent(ctx, params)
}

func (r *StudentRepository) GetByID(ctx context.Context, id string) (q.Student, error) {
	return r.Queries.GetStudentByID(ctx, id)
}

func (r *StudentRepository) ListByTenant(ctx context.Context, tenantID string, limit, offset int64) ([]q.Student, error) {
	return r.Queries.ListByTenant(ctx, q.ListByTenantParams{
		TenantID: sql.NullString{String: tenantID, Valid: true},
		Limit:    limit,
		Offset:   offset,
	})
}

func (r *StudentRepository) Update(ctx context.Context, id string, firstName, lastName, year *string, sectionID, departmentID *string) error {
	return r.Queries.UpdateStudent(ctx, q.UpdateStudentParams{
		ID:           id,
		FirstName:    toNullString(firstName),
		LastName:     toNullString(lastName),
		Year:         toNullString(year),
		SectionID:    toNullString(sectionID),
		DepartmentID: toNullString(departmentID),
	})
}

func (r *StudentRepository) SetStatus(ctx context.Context, id string, status string) error {
	return r.Queries.SetStudentStatus(ctx, q.SetStudentStatusParams{
		ID:     id,
		Status: sql.NullString{String: status, Valid: true},
	})
}

func (r *StudentRepository) GetByUserID(ctx context.Context, userID string) (q.Student, error) {
	return r.Queries.GetStudentByUserID(ctx, sql.NullString{String: userID, Valid: true})
}

func toNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}
