package repository

import (
	"context"
	"database/sql"
	"errors"
	"school-exam/internal/domain"
	queries "school-exam/internal/sqlc/gen"
)

type UserRepository struct {
	DB      *sql.DB
	Queries *queries.Queries
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db, Queries: queries.New(db)}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	row, err := r.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	var u domain.User
	u.ID = row.ID
	if row.TenantID.Valid {
		u.TenantID = &row.TenantID.String
	}
	u.Username = row.Username
	u.PasswordHash = row.PasswordHash
	if row.Email.Valid {
		u.Email = &row.Email.String
	}
	if row.RoleID.Valid {
		v := int64(row.RoleID.Int64)
		u.RoleID = &v
	}
	if row.Status.Valid {
		u.Status = row.Status.String
	}
	if row.RoleName.Valid {
		n := row.RoleName.String
		u.RoleName = &n
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	var tenant sql.NullString
	if u.TenantID != nil {
		tenant = sql.NullString{String: *u.TenantID, Valid: true}
	}
	var email sql.NullString
	if u.Email != nil {
		email = sql.NullString{String: *u.Email, Valid: true}
	}
	var role sql.NullInt64
	if u.RoleID != nil {
		role = sql.NullInt64{Int64: *u.RoleID, Valid: true}
	}
	return r.Queries.CreateUser(ctx, queries.CreateUserParams{
		ID:           u.ID,
		TenantID:     tenant,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Email:        email,
		RoleID:       role,
		Status:       sql.NullString{String: u.Status, Valid: true},
	})
}
