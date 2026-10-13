package result

import (
	"net/http"
	"school-exam/internal/middleware"
	"school-exam/internal/module/result"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *result.Usecase
}

func New(uc *result.Usecase) *Handler {
	return &Handler{uc: uc}
}

type CalculateResultReq struct {
	SectionID      string  `json:"section_id" binding:"required"`
	SubjectID      string  `json:"subject_id" binding:"required"`
	AcademicYearID *string `json:"academic_year_id"`
}

func (h *Handler) CalculateResults(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req CalculateResultReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.CalculateResults(c.Request.Context(), tenantID, req.SectionID, req.SubjectID, req.AcademicYearID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "results calculated and published"})
}

func (h *Handler) GetTransferSheet(c *gin.Context) {
	sectionID := c.Param("sectionId")
	list, err := h.uc.GetTransferSheet(c.Request.Context(), sectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) GetMyResults(c *gin.Context) {
	studentID := middleware.GetUserID(c)
	list, err := h.uc.GetStudentResults(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
