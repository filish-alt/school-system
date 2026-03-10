package exam_session

import (
	"net/http"

	"github.com/gin-gonic/gin"
	studentdto "school-exam/internal/dto/student"
	m "school-exam/internal/module/exam_session"
)

type Handler struct {
	UC *m.Usecase
}

func NewHandler(uc *m.Usecase) *Handler {
	return &Handler{UC: uc}
}

func (h *Handler) StartSession(c *gin.Context) {
	var req studentdto.StartSessionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	id, err := h.UC.StartSession(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) SaveAnswer(c *gin.Context) {
	var req studentdto.SaveAnswerRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.UC.SaveAnswer(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) SubmitSession(c *gin.Context) {
	var req studentdto.SubmitSessionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.UC.SubmitSession(c, req.SessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetSession(c *gin.Context) {
	id := c.Param("id")
	out, err := h.UC.GetSession(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) GetSessionResult(c *gin.Context) {
	id := c.Param("id")
	out, err := h.UC.GetSessionResult(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) ListMySessions(c *gin.Context) {
	out, err := h.UC.ListMySessions(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) ReportViolation(c *gin.Context) {
	var req studentdto.ReportViolationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.UC.ReportViolation(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListAllViolations(c *gin.Context) {
	out, err := h.UC.ListAllViolations(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
