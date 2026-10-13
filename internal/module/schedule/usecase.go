package schedule

import (
	"context"
	"school-exam/internal/domain"
	"school-exam/internal/repository"

	"github.com/google/uuid"
)

type Usecase struct {
	schRepo *repository.ScheduleRepository
}

func NewUsecase(schRepo *repository.ScheduleRepository) *Usecase {
	return &Usecase{schRepo: schRepo}
}

func (u *Usecase) CreateSchedule(ctx context.Context, s domain.TeacherSchedule) (string, error) {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	err := u.schRepo.Create(ctx, s)
	return s.ID, err
}

func (u *Usecase) ListByTeacher(ctx context.Context, teacherID string) ([]domain.TeacherSchedule, error) {
	return u.schRepo.ListByTeacher(ctx, teacherID)
}

func (u *Usecase) ListBySection(ctx context.Context, sectionID string) ([]domain.TeacherSchedule, error) {
	return u.schRepo.ListBySection(ctx, sectionID)
}

func (u *Usecase) DeleteSchedule(ctx context.Context, id string) error {
	return u.schRepo.Delete(ctx, id)
}
