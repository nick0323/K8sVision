package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware 请求日志记录中间件
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method
		traceId := c.GetString("traceId")

		// 记录请求开始日志
		logger.Info("request started",
			zap.String("traceId", traceId),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.String("clientIP", c.ClientIP()),
			zap.String("userAgent", c.Request.UserAgent()),
		)

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		// 记录请求完成日志
		logger.Info("request completed",
			zap.String("traceId", traceId),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("statusCode", statusCode),
			zap.Duration("latency", latency),
			zap.Int("bodySize", bodySize),
			zap.Strings("errors", c.Errors.Errors()),
		)
	}
}

// TraceMiddleware 请求追踪中间件
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取traceId，如果没有则生成新的
		traceId := c.GetHeader("X-Trace-ID")
		if traceId == "" {
			traceId = generateTraceID()
		}

		// 设置traceId到context中
		c.Set("traceId", traceId)

		// 在响应头中返回traceId
		c.Header("X-Trace-ID", traceId)

		c.Next()
	}
}

// generateTraceID 生成追踪ID
func generateTraceID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
