package schedule

import (
	"net/http"
	"school-exam/internal/domain"
	"school-exam/internal/middleware"
	"school-exam/internal/module/schedule"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *schedule.Usecase
}

func New(uc *schedule.Usecase) *Handler {
	return &Handler{uc: uc}
}

type CreateScheduleReq struct {
	TeacherID    string  `json:"teacher_id" binding:"required"`
	SubjectID    *string `json:"subject_id"`
	SectionID    *string `json:"section_id"`
	DayOfWeek    string  `json:"day_of_week" binding:"required"`
	StartTime    string  `json:"start_time" binding:"required"`
	EndTime      string  `json:"end_time" binding:"required"`
	Room         *string `json:"room"`
	ActivityType string  `json:"activity_type"`
}

func (h *Handler) CreateSchedule(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req CreateScheduleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activityType := "lecture"
	if req.ActivityType != "" {
		activityType = req.ActivityType
	}

	s := domain.TeacherSchedule{
		TenantID:     tenantID,
		TeacherID:    req.TeacherID,
		SubjectID:    req.SubjectID,
		SectionID:    req.SectionID,
		DayOfWeek:    req.DayOfWeek,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Room:         req.Room,
		ActivityType: activityType,
	}

	id, err := h.uc.CreateSchedule(c.Request.Context(), s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) ListMySchedule(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	list, err := h.uc.ListByTeacher(c.Request.Context(), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) ListByTeacher(c *gin.Context) {
	teacherID := c.Param("teacherId")
	list, err := h.uc.ListByTeacher(c.Request.Context(), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) ListBySection(c *gin.Context) {
	sectionID := c.Param("sectionId")
	list, err := h.uc.ListBySection(c.Request.Context(), sectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	if err := h.uc.DeleteSchedule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
