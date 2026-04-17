package examdto

import (
	"time"
	q "school-exam/internal/sqlc/gen"
)

type CreateExamRequest struct {
	Title           string     `json:"title" binding:"required"`
	SubjectID       string     `json:"subject_id" binding:"required"`
	SectionID       string     `json:"section_id" binding:"required"`
	DurationMinutes int64      `json:"duration_minutes" binding:"required"`
	StartTime       time.Time  `json:"start_time" binding:"required"`
	EndTime         *time.Time `json:"end_time"`
}

type UpdateExamRequest struct {
	ID              string     `json:"id" binding:"required"`
	Title           *string    `json:"title"`
	SubjectID       *string    `json:"subject_id"`
	SectionID       *string    `json:"section_id"`
	DurationMinutes *int64     `json:"duration_minutes"`
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	ShuffleOptions  *bool      `json:"shuffle_options"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AddQuestionsRequest struct {
	ExamID         string              `json:"exam_id" binding:"required"`
	Questions      []ExamQuestionInput `json:"questions" binding:"required"`
	ShuffleOptions bool                `json:"shuffle_options"`
}

type ExamQuestionInput struct {
	QuestionID string `json:"question_id" binding:"required"`
	Marks      *int64 `json:"marks"`
	OrderIndex int64  `json:"order_index"`
}

type AddRandomQuestionsRequest struct {
	ExamID         string `json:"exam_id" binding:"required"`
	QuestionBankID string `json:"question_bank_id" binding:"required"`
	Count          int64  `json:"count" binding:"required"`
	Marks          *int64 `json:"marks"`
	ShuffleOptions bool   `json:"shuffle_options"`
}

type ListQuery struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}

type ExamQuestionDetail struct {
	ID              string             `json:"id"`
	QuestionID      string             `json:"question_id"`
	QuestionText    string             `json:"question_text"`
	Type            string             `json:"type"`
	Marks           int64              `json:"marks"`
	DifficultyLevel string             `json:"difficulty_level"`
	OrderIndex      int64              `json:"order_index"`
	Options         []q.QuestionOption `json:"options,omitempty"`
}

type GetExamResponse struct {
	Exam      q.Exam               `json:"exam"`
	Questions []ExamQuestionDetail `json:"questions"`
}
