package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go-clean-api/pkg/apperror"
	pkgjwt "go-clean-api/pkg/jwt"
	"go-clean-api/pkg/response"
)

const UserIDKey = "userID"

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, apperror.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := pkgjwt.Parse(tokenStr, jwtSecret)
		if err != nil {
			response.Error(c, apperror.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}
