package academic

import (
	"net/http"
	"school-exam/internal/middleware"
	"school-exam/internal/module/academic"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *academic.Usecase
}

func New(uc *academic.Usecase) *Handler {
	return &Handler{uc: uc}
}

type CreateAcademicYearReq struct {
	TenantID  *string `json:"tenant_id"`
	YearName  string  `json:"year_name" binding:"required"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

func (h *Handler) CreateAcademicYear(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req CreateAcademicYearReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(tenantID) == "" && req.TenantID != nil && strings.TrimSpace(*req.TenantID) != "" {
		tenantID = strings.TrimSpace(*req.TenantID)
	}
	if strings.TrimSpace(tenantID) == "" {
		tenantID = strings.TrimSpace(c.Query("tenant_id"))
	}
	if strings.TrimSpace(tenantID) == "" {
		var firstID string
		_ = h.uc.GetFirstTenantID(c.Request.Context(), &firstID)
		if firstID != "" {
			tenantID = firstID
		}
	}

	if strings.TrimSpace(tenantID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	id, err := h.uc.CreateAcademicYear(c.Request.Context(), tenantID, req.YearName, req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "tenant_id": tenantID})
}

func (h *Handler) ListAcademicYears(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if strings.TrimSpace(tenantID) == "" {
		tenantID = strings.TrimSpace(c.Query("tenant_id"))
	}
	list, err := h.uc.ListAcademicYears(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) SetActiveAcademicYear(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if strings.TrimSpace(tenantID) == "" {
		tenantID = strings.TrimSpace(c.Query("tenant_id"))
	}
	id := c.Param("id")
	if err := h.uc.SetActiveAcademicYear(c.Request.Context(), tenantID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "active year updated"})
}

func (h *Handler) DeleteAcademicYear(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if strings.TrimSpace(tenantID) == "" {
		tenantID = strings.TrimSpace(c.Query("tenant_id"))
	}
	id := c.Param("id")
	if err := h.uc.DeleteAcademicYear(c.Request.Context(), tenantID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

type SettingReq struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func (h *Handler) GetSettings(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if strings.TrimSpace(tenantID) == "" {
		tenantID = strings.TrimSpace(c.Query("tenant_id"))
	}
	if strings.TrimSpace(tenantID) == "" {
		var firstID string
		_ = h.uc.GetFirstTenantID(c.Request.Context(), &firstID)
		if firstID != "" {
			tenantID = firstID
		}
	}
	list, err := h.uc.ListSystemSettings(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(list) == 0 {
		var firstID string
		_ = h.uc.GetFirstTenantID(c.Request.Context(), &firstID)
		if firstID != "" && firstID != tenantID {
			if fallbackList, err := h.uc.ListSystemSettings(c.Request.Context(), firstID); err == nil && len(fallbackList) > 0 {
				c.JSON(http.StatusOK, fallbackList)
				return
			}
		}
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) SaveSetting(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if strings.TrimSpace(tenantID) == "" {
		tenantID = strings.TrimSpace(c.Query("tenant_id"))
	}
	if strings.TrimSpace(tenantID) == "" {
		var firstID string
		_ = h.uc.GetFirstTenantID(c.Request.Context(), &firstID)
		if firstID != "" {
			tenantID = firstID
		}
	}
	var req SettingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.SaveSystemSetting(c.Request.Context(), tenantID, req.Key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "setting saved"})
}
