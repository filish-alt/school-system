package studentdto

import (
	"time"
)

type ReportViolationRequest struct {
	SessionID     string `json:"session_id" binding:"required"`
	ViolationType string `json:"violation_type" binding:"required"`
}

type ViolationResponse struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id"`
	ViolationType string    `json:"violation_type"`
	CreatedAt     time.Time `json:"created_at"`
	StudentID     string    `json:"student_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	ExamID        string    `json:"exam_id"`
	ExamTitle     string    `json:"exam_title"`
}
