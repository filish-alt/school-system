package complaint

import (
	"net/http"
	"school-exam/internal/domain"
	"school-exam/internal/middleware"
	"school-exam/internal/module/complaint"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *complaint.Usecase
}

func New(uc *complaint.Usecase) *Handler {
	return &Handler{uc: uc}
}

type CreateComplaintReq struct {
	SubjectID    *string `json:"subject_id"`
	ExamID       *string `json:"exam_id"`
	CourseworkID *string `json:"coursework_id"`
	QuestionID   *string `json:"question_id"`
	Category     string  `json:"category" binding:"required"`
	Description  string  `json:"description" binding:"required"`
}

func (h *Handler) CreateComplaint(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	studentID := middleware.GetUserID(c)

	var req CreateComplaintReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comp := domain.StudentComplaint{
		TenantID:     tenantID,
		StudentID:    studentID,
		SubjectID:    req.SubjectID,
		ExamID:       req.ExamID,
		CourseworkID: req.CourseworkID,
		QuestionID:   req.QuestionID,
		Category:     req.Category,
		Description:  req.Description,
	}

	id, err := h.uc.CreateComplaint(c.Request.Context(), comp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) ListMyComplaints(c *gin.Context) {
	studentID := middleware.GetUserID(c)
	list, err := h.uc.ListByStudent(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) ListSchoolComplaints(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	status := c.Query("status")
	list, err := h.uc.ListByTenant(c.Request.Context(), tenantID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

type ResolveComplaintReq struct {
	Status   string `json:"status" binding:"required"` // escalated_to_dept_head, resolved, rejected
	Response string `json:"response" binding:"required"`
}

func (h *Handler) ResolveComplaint(c *gin.Context) {
	id := c.Param("id")
	userID := middleware.GetUserID(c)

	var req ResolveComplaintReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.ResolveComplaint(c.Request.Context(), id, req.Status, req.Response, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "complaint updated"})
}
