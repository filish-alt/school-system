package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"school-exam/internal/security"
)

func Auth(ts security.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(h, "Bearer ")
		claims, err := ts.Parse(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

