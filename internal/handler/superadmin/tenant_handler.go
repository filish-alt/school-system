package superadmin

import (
	"net/http"
	"strings"

	"school-exam/internal/module/superadmin"

	sdto "school-exam/internal/dto/superadmin"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UC *superadmin.Usecase
}

func New(uc *superadmin.Usecase) *Handler {
	return &Handler{UC: uc}
}

func (h *Handler) CreateTenant(c *gin.Context) {
	var req sdto.CreateTenantRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	res, err := h.UC.CreateTenant(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateTenant(c *gin.Context) {
	var req sdto.UpdateTenantRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.UC.UpdateTenant(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetTenant(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	out, err := h.UC.GetTenant(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) ListTenants(c *gin.Context) {
	status := c.Query("status")
	var ptr *string
	if strings.TrimSpace(status) != "" {
		s := strings.TrimSpace(status)
		ptr = &s
	}
	var q sdto.TenantListQuery
	_ = c.ShouldBindQuery(&q)
	out, err := h.UC.ListTenants(c, ptr, q.Page, q.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) DeleteTenant(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	if err := h.UC.DeleteTenant(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ActivateTenant(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	if err := h.UC.ActivateTenant(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeactivateTenant(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	if err := h.UC.DeactivateTenant(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateStudent(c *gin.Context) {
	var req sdto.CreateStudentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	res, err := h.UC.CreateStudent(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateStudent(c *gin.Context) {
	var req sdto.UpdateStudentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	id := c.Param("id")
	if id != "" {
		req.ID = id
	}
	if err := h.UC.UpdateStudent(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetStudent(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	out, err := h.UC.Students.GetByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) ListStudents(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if strings.TrimSpace(tenantID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id required"})
		return
	}
	var q sdto.ListQuery
	_ = c.ShouldBindQuery(&q)
	out, err := h.UC.Students.ListByTenant(c, tenantID, q.PageSize, (q.Page-1)*q.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	if err := h.UC.Students.SetStatus(c, id, "inactive"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ImportStudents(c *gin.Context) {
	tenantID := c.PostForm("tenant_id")
	if strings.TrimSpace(tenantID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id required"})
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot open file"})
		return
	}
	defer f.Close()
	results, err := h.UC.ImportStudents(c, tenantID, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}
