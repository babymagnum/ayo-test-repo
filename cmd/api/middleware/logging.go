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

		// Skip Swagger & static assets
		if strings.HasPrefix(c.Request.URL.Path, "/v1/swagger") {
			c.Next()
			return
		}

		// Get request-id from header or generate new
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set(CtxRequestID, reqID)
		c.Header("Content-Type", "application/json")

		// Limit size
		const maxSize = 2048

		// var query strings.Builder

		// first := true
		// for key, values := range r.URL.Query() {
		// 	for _, v := range values {
		// 		if !first {
		// 			query.WriteString("&")
		// 		} else {
		// 			query.WriteString("?")
		// 		}
		// 		first = false
		// 		query.WriteString(fmt.Sprintf("%s=%s", key, v))
		// 	}
		// }

		// Log request body for non-GET
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
				// Restore body for downstream
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		} else {
			logger.Info("REQUEST",
				zap.String("RequestId", reqID),
				zap.String("Method", c.Request.Method),
				zap.String("Path", c.Request.URL.String()),
			)
		}

		// --- Execute Handler Chain ---
		c.Next()

		// Log the captured status, path, latency, and captured response body
		logger.Info("RESPONSE",
			zap.String("RequestId", reqID),
			zap.Int("Code", c.Writer.Status()),
			zap.String("Method", c.Request.Method),
			// zap.String("Path", fmt.Sprintf("%v%v", r.URL.Path, query.String())),
			zap.String("Path", c.Request.URL.Path),
			zap.String("Latency", time.Since(start).String()),
		)
	}
}
