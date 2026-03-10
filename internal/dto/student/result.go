package studentdto

import (
	"time"
	q "school-exam/internal/sqlc/gen"
)

type ResultSummaryResponse struct {
	SessionID  string    `json:"session_id"`
	ExamID     string    `json:"exam_id"`
	ExamTitle  string    `json:"exam_title"`
	TotalScore int64     `json:"total_score"`
	Status     string    `json:"status"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}

type AnswerReviewDetail struct {
	QuestionID       string             `json:"question_id"`
	QuestionText     string             `json:"question_text"`
	QuestionType     string             `json:"question_type"`
	SelectedOptionID string             `json:"selected_option_id"`
	CorrectOptionID  string             `json:"correct_option_id"`
	IsCorrect        bool               `json:"is_correct"`
	Score            int64              `json:"score"`
	MaxMarks         int64              `json:"max_marks"`
	Options          []q.QuestionOption `json:"options"`
}

type ExamResultResponse struct {
	Summary   ResultSummaryResponse `json:"summary"`
	Breakdown []AnswerReviewDetail  `json:"breakdown"`
}
