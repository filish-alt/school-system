package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"school-exam/internal/usecase"
)

type AuthHandler struct {
	Auth *usecase.AuthUsecase
}

func NewAuthHandler(a *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{Auth: a}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req usecase.LoginRequest
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

