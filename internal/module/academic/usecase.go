package academic

import (
	"context"
	"school-exam/internal/domain"
	"school-exam/internal/repository"

	"github.com/google/uuid"
)

type Usecase struct {
	ayRepo  *repository.AcademicYearRepository
	setRepo *repository.SystemSettingRepository
}

func NewUsecase(ayRepo *repository.AcademicYearRepository, setRepo *repository.SystemSettingRepository) *Usecase {
	return &Usecase{ayRepo: ayRepo, setRepo: setRepo}
}

func (u *Usecase) GetFirstTenantID(ctx context.Context, out *string) error {
	id, err := u.ayRepo.GetFirstTenantID(ctx)
	if err == nil {
		*out = id
	}
	return err
}

func (u *Usecase) CreateAcademicYear(ctx context.Context, tenantID, yearName string, startDate, endDate *string) (string, error) {
	id := uuid.NewString()
	err := u.ayRepo.Create(ctx, domain.AcademicYear{
		ID:        id,
		TenantID:  tenantID,
		YearName:  yearName,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  false,
	})
	return id, err
}

func (u *Usecase) ListAcademicYears(ctx context.Context, tenantID string) ([]domain.AcademicYear, error) {
	return u.ayRepo.ListByTenant(ctx, tenantID)
}

func (u *Usecase) SetActiveAcademicYear(ctx context.Context, tenantID, id string) error {
	return u.ayRepo.SetActive(ctx, tenantID, id)
}

func (u *Usecase) DeleteAcademicYear(ctx context.Context, tenantID, id string) error {
	return u.ayRepo.Delete(ctx, tenantID, id)
}

func (u *Usecase) GetSystemSetting(ctx context.Context, tenantID, key string) (*domain.SystemSetting, error) {
	return u.setRepo.Get(ctx, tenantID, key)
}

func (u *Usecase) ListSystemSettings(ctx context.Context, tenantID string) ([]domain.SystemSetting, error) {
	return u.setRepo.List(ctx, tenantID)
}

func (u *Usecase) SaveSystemSetting(ctx context.Context, tenantID, key, value string) error {
	id := uuid.NewString()
	return u.setRepo.Upsert(ctx, id, tenantID, key, value)
}
