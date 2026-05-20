package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recoverer captures panics, logs the stack trace, and returns a 500 error.
func Recoverer(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Set up the defer function to run after the handler chain finishes (or panics)
		defer func() {
			if rvr := recover(); rvr != nil {
				logger.Error(
					"PANIC RECOVERED",
					zap.Any("panic", rvr),
					zap.Any("stack", string(debug.Stack())),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Internal server error",
				})
			}
		}()

		c.Next()
	}
}
