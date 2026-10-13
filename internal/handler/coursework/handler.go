package coursework

import (
	"net/http"
	"school-exam/internal/domain"
	"school-exam/internal/middleware"
	"school-exam/internal/module/coursework"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *coursework.Usecase
}

func New(uc *coursework.Usecase) *Handler {
	return &Handler{uc: uc}
}

type CreateCourseworkReq struct {
	AcademicYearID   *string `json:"academic_year_id"`
	SubjectID        string  `json:"subject_id" binding:"required"`
	SectionID        string  `json:"section_id" binding:"required"`
	Title            string  `json:"title" binding:"required"`
	Type             string  `json:"type" binding:"required"` // quiz, mid_exam, assignment, project
	MaxMarks         float64 `json:"max_marks" binding:"required"`
	WeightPercentage float64 `json:"weight_percentage" binding:"required"`
	DueDate          *string `json:"due_date"`
}

func (h *Handler) CreateCoursework(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	teacherID := middleware.GetUserID(c) // or teacher lookup

	var req CreateCourseworkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DueDate != nil && *req.DueDate == "" {
		req.DueDate = nil
	}
	if req.AcademicYearID != nil && *req.AcademicYearID == "" {
		req.AcademicYearID = nil
	}

	cw := domain.CourseworkAssessment{
		TenantID:         tenantID,
		AcademicYearID:   req.AcademicYearID,
		SubjectID:        req.SubjectID,
		SectionID:        req.SectionID,
		TeacherID:        teacherID,
		Title:            req.Title,
		Type:             req.Type,
		MaxMarks:         req.MaxMarks,
		WeightPercentage: req.WeightPercentage,
		DueDate:          req.DueDate,
	}

	id, err := h.uc.CreateCoursework(c.Request.Context(), cw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) ListBySectionSubject(c *gin.Context) {
	sectionID := c.Query("section_id")
	subjectID := c.Query("subject_id")
	if sectionID == "" || subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "section_id and subject_id required"})
		return
	}
	list, err := h.uc.ListBySectionSubject(c.Request.Context(), sectionID, subjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) ListMyCourseworks(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	list, err := h.uc.ListByTeacher(c.Request.Context(), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

type SaveMarkReq struct {
	CourseworkID  string  `json:"coursework_id" binding:"required"`
	StudentID     string  `json:"student_id" binding:"required"`
	MarksObtained float64 `json:"marks_obtained"`
	Remarks       *string `json:"remarks"`
}

func (h *Handler) SaveMark(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	var req SaveMarkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mark := domain.StudentCourseworkMark{
		CourseworkID:  req.CourseworkID,
		StudentID:     req.StudentID,
		MarksObtained: req.MarksObtained,
		Remarks:       req.Remarks,
		GradedBy:      &teacherID,
	}

	if err := h.uc.SaveMark(c.Request.Context(), mark); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "mark saved"})
}

func (h *Handler) ListMarksByCoursework(c *gin.Context) {
	courseworkID := c.Param("courseworkId")
	list, err := h.uc.ListMarksByCoursework(c.Request.Context(), courseworkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) ListMyMarks(c *gin.Context) {
	studentID := middleware.GetUserID(c)
	list, err := h.uc.ListMarksByStudent(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
