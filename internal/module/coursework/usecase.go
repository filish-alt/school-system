package coursework

import (
	"context"
	"school-exam/internal/domain"
	"school-exam/internal/repository"

	"github.com/google/uuid"
)

type Usecase struct {
	cwRepo *repository.CourseworkRepository
}

func NewUsecase(cwRepo *repository.CourseworkRepository) *Usecase {
	return &Usecase{cwRepo: cwRepo}
}

func (u *Usecase) CreateCoursework(ctx context.Context, c domain.CourseworkAssessment) (string, error) {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	err := u.cwRepo.Create(ctx, c)
	return c.ID, err
}

func (u *Usecase) ListBySectionSubject(ctx context.Context, sectionID, subjectID string) ([]domain.CourseworkAssessment, error) {
	return u.cwRepo.ListBySectionSubject(ctx, sectionID, subjectID)
}

func (u *Usecase) ListByTeacher(ctx context.Context, teacherID string) ([]domain.CourseworkAssessment, error) {
	return u.cwRepo.ListByTeacher(ctx, teacherID)
}

func (u *Usecase) SaveMark(ctx context.Context, m domain.StudentCourseworkMark) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return u.cwRepo.UpsertMark(ctx, m)
}

func (u *Usecase) ListMarksByCoursework(ctx context.Context, courseworkID string) ([]domain.StudentCourseworkMark, error) {
	return u.cwRepo.ListMarksByCoursework(ctx, courseworkID)
}

func (u *Usecase) ListMarksByStudent(ctx context.Context, studentID string) ([]domain.StudentCourseworkMark, error) {
	return u.cwRepo.ListMarksByStudent(ctx, studentID)
}
