package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"school-exam/internal/security"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	roleSet := map[string]struct{}{}
	for _, r := range roles {
		roleSet[r] = struct{}{}
	}
	return func(c *gin.Context) {
		v, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		claims := v.(*security.Claims)
		if claims.Role == nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if _, ok := roleSet[*claims.Role]; !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

