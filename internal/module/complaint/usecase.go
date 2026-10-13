package complaint

import (
	"context"
	"school-exam/internal/domain"
	"school-exam/internal/repository"

	"github.com/google/uuid"
)

type Usecase struct {
	compRepo *repository.ComplaintRepository
}

func NewUsecase(compRepo *repository.ComplaintRepository) *Usecase {
	return &Usecase{compRepo: compRepo}
}

func (u *Usecase) CreateComplaint(ctx context.Context, c domain.StudentComplaint) (string, error) {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	err := u.compRepo.Create(ctx, c)
	return c.ID, err
}

func (u *Usecase) ListByStudent(ctx context.Context, studentID string) ([]domain.StudentComplaint, error) {
	return u.compRepo.ListByStudent(ctx, studentID)
}

func (u *Usecase) ListByTenant(ctx context.Context, tenantID, statusFilter string) ([]domain.StudentComplaint, error) {
	return u.compRepo.ListByTenant(ctx, tenantID, statusFilter)
}

func (u *Usecase) ResolveComplaint(ctx context.Context, id, status, response, resolvedBy string) error {
	return u.compRepo.UpdateStatus(ctx, id, status, response, resolvedBy)
}
