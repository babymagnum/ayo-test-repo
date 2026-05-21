package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ctxKey string

const (
	CtxRequestID ctxKey = "request-id"
)

func Logging(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		if strings.HasPrefix(c.Request.URL.Path, "/v1/swagger") {
			c.Next()
			return
		}

		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set(CtxRequestID, reqID)
		c.Header("Content-Type", "application/json")

		const maxSize = 2048

		if c.Request.Method != http.MethodGet {
			body, _ := io.ReadAll(c.Request.Body)
			if len(body) > 0 {
				truncated := body
				if len(truncated) > maxSize {
					truncated = truncated[:maxSize]
				}
				logger.Info("REQUEST",
					zap.String("RequestId", reqID),
					zap.String("Method", c.Request.Method),
					zap.String("Path", c.Request.URL.Path),
					zap.String("Body", string(truncated)),
				)

				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		} else {
			logger.Info("REQUEST",
				zap.String("RequestId", reqID),
				zap.String("Method", c.Request.Method),
				zap.String("Path", c.Request.URL.String()),
			)
		}

		c.Next()

		logger.Info("RESPONSE",
			zap.String("RequestId", reqID),
			zap.Int("Code", c.Writer.Status()),
			zap.String("Method", c.Request.Method),
			zap.String("Path", c.Request.URL.Path),
			zap.String("Latency", time.Since(start).String()),
		)
	}
}
