package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	Secret string
	TTL    time.Duration
}

type Claims struct {
	UserID   string `json:"uid"`
	TenantID *string `json:"tid,omitempty"`
	Role     *string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

func (t TokenService) Sign(userID string, tenantID *string, role *string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(t.TTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.Secret))
}

func (t TokenService) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.Secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}
	if c, ok := token.Claims.(*Claims); ok && token.Valid {
		return c, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

