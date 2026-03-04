package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	mod "school-exam/internal/module/auth"
)

type Handler struct {
	Auth *mod.Usecase
}

func New(a *mod.Usecase) *Handler {
	return &Handler{Auth: a}
}

func (h *Handler) Login(c *gin.Context) {
	var req mod.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	res, err := h.Auth.Login(c, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, res)
}

