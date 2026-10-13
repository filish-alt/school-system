package exam_control

import (
	"context"
	"school-exam/internal/domain"
	"school-exam/internal/repository"

	"github.com/google/uuid"
)

type Usecase struct {
	ecRepo *repository.ExamControlRepository
}

func NewUsecase(ecRepo *repository.ExamControlRepository) *Usecase {
	return &Usecase{ecRepo: ecRepo}
}

func (u *Usecase) AddExtension(ctx context.Context, ext domain.ExamExtension) error {
	if ext.ID == "" {
		ext.ID = uuid.NewString()
	}
	return u.ecRepo.AddExtension(ctx, ext)
}

func (u *Usecase) ListExtensions(ctx context.Context, examID string) ([]domain.ExamExtension, error) {
	return u.ecRepo.ListExtensions(ctx, examID)
}

func (u *Usecase) AddRemoteAuthorization(ctx context.Context, auth domain.ExamRemoteAuthorization) error {
	if auth.ID == "" {
		auth.ID = uuid.NewString()
	}
	return u.ecRepo.AddRemoteAuthorization(ctx, auth)
}

func (u *Usecase) ListRemoteAuthorizations(ctx context.Context, examID string) ([]domain.ExamRemoteAuthorization, error) {
	return u.ecRepo.ListRemoteAuthorizations(ctx, examID)
}

func (u *Usecase) DropQuestion(ctx context.Context, dq domain.DroppedQuestion) error {
	if dq.ID == "" {
		dq.ID = uuid.NewString()
	}
	return u.ecRepo.DropQuestion(ctx, dq)
}

func (u *Usecase) CorrectAnswerKey(ctx context.Context, corr domain.AnswerKeyCorrection) error {
	if corr.ID == "" {
		corr.ID = uuid.NewString()
	}
	return u.ecRepo.CorrectAnswerKey(ctx, corr)
}

func (u *Usecase) UpdateExamControls(ctx context.Context, examID string, isStarted, isVisible, isLocked *bool, status, approvalStatus, accessCode *string, approvedBy *string) error {
	return u.ecRepo.UpdateExamControls(ctx, examID, isStarted, isVisible, isLocked, status, approvalStatus, accessCode, approvedBy)
}

func (u *Usecase) AddMatchPairs(ctx context.Context, pairs []domain.QuestionMatchPair) error {
	for i := range pairs {
		if pairs[i].ID == "" {
			pairs[i].ID = uuid.NewString()
		}
	}
	return u.ecRepo.AddMatchPairs(ctx, pairs)
}

func (u *Usecase) ListMatchPairs(ctx context.Context, questionID string) ([]domain.QuestionMatchPair, error) {
	return u.ecRepo.ListMatchPairs(ctx, questionID)
}
