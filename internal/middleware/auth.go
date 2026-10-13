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
func GetClaims(c *gin.Context) *security.Claims {
	v, ok := c.Get("claims")
	if !ok {
		return nil
	}
	claims, ok := v.(*security.Claims)
	if !ok {
		return nil
	}
	return claims
}

func GetUserID(c *gin.Context) string {
	claims := GetClaims(c)
	if claims != nil {
		return claims.UserID
	}
	return ""
}

func GetTenantID(c *gin.Context) string {
	claims := GetClaims(c)
	if claims != nil && claims.TenantID != nil {
		return *claims.TenantID
	}
	return ""
}

func GetRole(c *gin.Context) *string {
	claims := GetClaims(c)
	if claims != nil {
		return claims.Role
	}
	return nil
}
