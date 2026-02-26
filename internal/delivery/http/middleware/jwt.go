package middleware

import (
	"fmt"
	"strings"

	"go-clean-architecture-boilerplate/internal/delivery/http/response"
	"go-clean-architecture-boilerplate/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// ContextUserID is the key used to store the authenticated user's ID in the Gin context.
	ContextUserID = "userID"
	// ContextUserEmail is the key used to store the authenticated user's email in the Gin context.
	ContextUserEmail = "userEmail"
)

// JWTAuth returns a Gin middleware that validates the Bearer token in the Authorization header.
// The jwtSecret is injected at construction time from config — no os.Getenv inside.
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			response.Unauthorized(c, "authorization header format must be: Bearer {token}")
			c.Abort()
			return
		}

		claims := &usecase.JWTClaims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUserEmail, claims.Email)
		c.Next()
	}
}
