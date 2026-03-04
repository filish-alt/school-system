package auth

import (
	"context"
	"errors"
	"school-exam/internal/domain"
	"school-exam/internal/repository"
	"school-exam/internal/security"
)

type Usecase struct {
	Users        *repository.UserRepository
	TokenService security.TokenService
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func NewAuthUsecase(repo *repository.UserRepository, ts security.TokenService) *Usecase {
	return &Usecase{Users: repo, TokenService: ts}
}

func (a *Usecase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	u, err := a.Users.GetByUsername(ctx, req.Username)
	if err != nil || u == nil {
		return nil, ErrInvalidCredentials
	}
	if err := security.CheckPassword(u.PasswordHash, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}
	token, err := a.TokenService.Sign(u.ID, u.TenantID, u.RoleName)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{Token: token}, nil
}

func (a *Usecase) SeedSuperAdmin(ctx context.Context, username, password string) error {
	existing, err := a.Users.GetByUsername(ctx, username)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}
	h, err := security.HashPassword(password)
	if err != nil {
		return err
	}
	roleID := int64(1)
	u := domain.User{
		ID:           username,
		Username:     username,
		PasswordHash: h,
		Status:       "active",
		RoleID:       &roleID,
	}
	return a.Users.Create(ctx, u)
}
