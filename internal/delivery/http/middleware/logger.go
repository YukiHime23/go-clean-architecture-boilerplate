package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapLogger returns a Gin middleware that logs each request/response using Uber Zap.
// It captures: method, path, status code, latency, client IP, and any request-specific errors.
func ZapLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process the request.
		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		errorMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
		}

		if query != "" {
			fields = append(fields, zap.String("query", query))
		}
		if errorMsg != "" {
			fields = append(fields, zap.String("error", errorMsg))
		}

		// Use different log levels based on HTTP status code.
		switch {
		case statusCode >= 500:
			log.Error("Server Error", fields...)
		case statusCode >= 400:
			log.Warn("Client Error", fields...)
		default:
			log.Info("Request", fields...)
		}
	}
}
