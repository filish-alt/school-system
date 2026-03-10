package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	mod "school-exam/internal/module/auth"
	"school-exam/internal/security"
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

func (h *Handler) UpdatePassword(c *gin.Context) {
	var req mod.UpdatePasswordRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	v, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims := v.(*security.Claims)

	if err := h.Auth.UpdatePassword(c, claims.UserID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

