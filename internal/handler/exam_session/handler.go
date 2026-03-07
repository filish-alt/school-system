package exam_session

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sessiondto "school-exam/internal/dto/student"
	m "school-exam/internal/module/exam_session"
)

type Handler struct {
	UC *m.Usecase
}

func NewHandler(uc *m.Usecase) *Handler {
	return &Handler{UC: uc}
}

func (h *Handler) StartSession(c *gin.Context) {
	var req sessiondto.StartSessionRequest
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
	var req sessiondto.SaveAnswerRequest
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
	var req sessiondto.SubmitSessionRequest
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
