package studentdto

import (
    "time"
)

type StartSessionRequest struct {
	ExamID string `json:"exam_id" binding:"required"`
}

type SaveAnswerRequest struct {
	SessionID        string  `json:"session_id" binding:"required"`
	QuestionID       string  `json:"question_id" binding:"required"`
	AnswerText       *string `json:"answer_text"`
	SelectedOptionID *string `json:"selected_option_id"`
}

type SubmitSessionRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

type SessionResponse struct {
	ID         string    `json:"id"`
	ExamID     string    `json:"exam_id"`
	StudentID  string    `json:"student_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status"`
	TotalScore int64     `json:"total_score"`
}
