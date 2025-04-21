package middleware

import (
	"JavaCode/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

// Logger returns a Gin middleware that logs the beginning and end of each HTTP request.
//
// It logs the request method, path, status code, client IP, and request latency
// using a structured logger (logrus).
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		utils.Logger.WithFields(logrus.Fields{
			"stage":    "begin",
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"clientIP": c.ClientIP(),
		}).Info("request started")

		c.Next()

		latency := time.Since(start)

		utils.Logger.WithFields(logrus.Fields{
			"stage":    "finish",
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   c.Writer.Status(),
			"latency":  latency,
			"clientIP": c.ClientIP(),
		}).Info("request finished")

	}
}
