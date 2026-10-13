package exam_control

import (
	"net/http"
	"school-exam/internal/domain"
	"school-exam/internal/middleware"
	"school-exam/internal/module/exam_control"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *exam_control.Usecase
}

func New(uc *exam_control.Usecase) *Handler {
	return &Handler{uc: uc}
}

type ExtensionReq struct {
	ExamID          string  `json:"exam_id" binding:"required"`
	StudentID       *string `json:"student_id"`
	ExtendedMinutes int     `json:"extended_minutes" binding:"required"`
	Reason          string  `json:"reason" binding:"required"`
}

func (h *Handler) AddExtension(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	var req ExtensionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ext := domain.ExamExtension{
		ExamID:          req.ExamID,
		StudentID:       req.StudentID,
		ExtendedMinutes: req.ExtendedMinutes,
		Reason:          req.Reason,
		AuthorizedBy:    teacherID,
	}
	if err := h.uc.AddExtension(c.Request.Context(), ext); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "extension granted"})
}

func getExamID(c *gin.Context) string {
	if id := c.Param("id"); id != "" {
		return id
	}
	return c.Param("examId")
}

func (h *Handler) ListExtensions(c *gin.Context) {
	examID := getExamID(c)
	list, err := h.uc.ListExtensions(c.Request.Context(), examID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

type RemoteAuthReq struct {
	ExamID     string `json:"exam_id" binding:"required"`
	StudentID  string `json:"student_id" binding:"required"`
	ValidUntil string `json:"valid_until" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
}

func (h *Handler) AddRemoteAuth(c *gin.Context) {
	authorizerID := middleware.GetUserID(c)
	var req RemoteAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ra := domain.ExamRemoteAuthorization{
		ExamID:     req.ExamID,
		StudentID:  req.StudentID,
		ValidUntil: req.ValidUntil,
		Reason:     req.Reason,
		ApprovedBy: authorizerID,
	}
	if err := h.uc.AddRemoteAuthorization(c.Request.Context(), ra); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "remote authorization granted"})
}

func (h *Handler) ListRemoteAuths(c *gin.Context) {
	examID := getExamID(c)
	list, err := h.uc.ListRemoteAuthorizations(c.Request.Context(), examID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

type DropQuestionReq struct {
	ExamID     string `json:"exam_id" binding:"required"`
	QuestionID string `json:"question_id" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
	CreditType string `json:"credit_type"` // full_credit or reduce_total
}

func (h *Handler) DropQuestion(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req DropQuestionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creditType := "full_credit"
	if req.CreditType != "" {
		creditType = req.CreditType
	}
	dq := domain.DroppedQuestion{
		ExamID:     req.ExamID,
		QuestionID: req.QuestionID,
		Reason:     req.Reason,
		CreditType: creditType,
		DroppedBy:  userID,
	}
	if err := h.uc.DropQuestion(c.Request.Context(), dq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "question dropped and scores recalculated"})
}

type AnswerKeyCorrectionReq struct {
	ExamID     string  `json:"exam_id" binding:"required"`
	QuestionID string  `json:"question_id" binding:"required"`
	OldAnswer  *string `json:"old_answer"`
	NewAnswer  *string `json:"new_answer" binding:"required"`
	Reason     string  `json:"reason" binding:"required"`
}

func (h *Handler) CorrectAnswerKey(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req AnswerKeyCorrectionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	corr := domain.AnswerKeyCorrection{
		ExamID:      req.ExamID,
		QuestionID:  req.QuestionID,
		OldAnswer:   req.OldAnswer,
		NewAnswer:   req.NewAnswer,
		Reason:      req.Reason,
		CorrectedBy: userID,
	}
	if err := h.uc.CorrectAnswerKey(c.Request.Context(), corr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "answer key corrected and marks recalculated"})
}

type UpdateExamControlsReq struct {
	IsStarted      *bool   `json:"is_started"`
	IsVisible      *bool   `json:"is_visible"`
	IsLocked       *bool   `json:"is_locked"`
	Status         *string `json:"status"`
	ApprovalStatus *string `json:"approval_status"`
	AccessCode     *string `json:"access_code"`
}

func (h *Handler) UpdateControls(c *gin.Context) {
	examID := getExamID(c)
	userID := middleware.GetUserID(c)

	var req UpdateExamControlsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var approvedBy *string
	if req.ApprovalStatus != nil && *req.ApprovalStatus == "approved" {
		approvedBy = &userID
	}

	if err := h.uc.UpdateExamControls(c.Request.Context(), examID, req.IsStarted, req.IsVisible, req.IsLocked, req.Status, req.ApprovalStatus, req.AccessCode, approvedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "exam controls updated"})
}

type MatchPairItem struct {
	LeftPrompt string `json:"left_prompt" binding:"required"`
	RightMatch string `json:"right_match" binding:"required"`
	OrderIndex int    `json:"order_index"`
}

func (h *Handler) SetMatchPairs(c *gin.Context) {
	questionID := c.Param("questionId")
	var req []MatchPairItem
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pairs []domain.QuestionMatchPair
	for _, p := range req {
		pairs = append(pairs, domain.QuestionMatchPair{
			QuestionID: questionID,
			LeftPrompt: p.LeftPrompt,
			RightMatch: p.RightMatch,
			OrderIndex: p.OrderIndex,
		})
	}

	if err := h.uc.AddMatchPairs(c.Request.Context(), pairs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "matching pairs saved"})
}

func (h *Handler) GetMatchPairs(c *gin.Context) {
	questionID := c.Param("questionId")
	list, err := h.uc.ListMatchPairs(c.Request.Context(), questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
