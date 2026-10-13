package result

import (
	"context"
	"school-exam/internal/domain"
	"school-exam/internal/repository"
)

type Usecase struct {
	resRepo *repository.ResultRepository
}

func NewUsecase(resRepo *repository.ResultRepository) *Usecase {
	return &Usecase{resRepo: resRepo}
}

func (u *Usecase) CalculateResults(ctx context.Context, tenantID, sectionID, subjectID string, academicYearID *string) error {
	return u.resRepo.CalculateAndSaveResults(ctx, tenantID, sectionID, subjectID, academicYearID)
}

func (u *Usecase) GetTransferSheet(ctx context.Context, sectionID string) ([]domain.AssessmentResult, error) {
	return u.resRepo.GetTransferSheet(ctx, sectionID)
}

func (u *Usecase) GetStudentResults(ctx context.Context, studentID string) ([]domain.AssessmentResult, error) {
	return u.resRepo.GetStudentResults(ctx, studentID)
}
