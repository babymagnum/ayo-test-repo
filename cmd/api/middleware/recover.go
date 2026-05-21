package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recoverer(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

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
