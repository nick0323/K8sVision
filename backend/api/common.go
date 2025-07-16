package api

import (
	"time"

	"github.com/nick0323/K8sVision/backend/model"

	"strings"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var jwtSecret []byte

// InitJWTSecret 初始化 JWT 密钥，优先环境变量，其次 viper 配置，最后默认
func InitJWTSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = viper.GetString("jwt.secret")
	}
	if secret == "" {
		secret = "k8svision-secret-key"
	}
	jwtSecret = []byte(secret)
}

// GetTraceID 从 gin.Context 获取 traceId
func GetTraceID(c *gin.Context) string {
	if v, ok := c.Get("traceId"); ok {
		if tid, ok := v.(string); ok {
			return tid
		}
	}
	return ""
}

// ResponseError 统一错误响应，支持日志分级，返回标准 APIResponse 结构体
func ResponseError(c *gin.Context, logger *zap.Logger, err error, code int) {
	traceId := GetTraceID(c)
	if logger != nil {
		if code >= 500 {
			logger.Error("api error", zap.String("traceId", traceId), zap.Error(err))
		} else {
			logger.Warn("api warn", zap.String("traceId", traceId), zap.Error(err))
		}
	}
	c.JSON(code, model.APIResponse{
		Code:      code,
		Message:   err.Error(),
		Data:      nil,
		TraceID:   traceId,
		Timestamp: time.Now().Unix(),
	})
}

// ResponseOK 统一成功响应
func ResponseOK(c *gin.Context, data interface{}, msg string, page *model.PageMeta) {
	traceId := GetTraceID(c)
	c.JSON(200, model.APIResponse{
		Code:      0,
		Message:   msg,
		Data:      data,
		TraceID:   traceId,
		Timestamp: time.Now().Unix(),
		Page:      page,
	})
}

// Paginate 对任意切片类型进行分页，返回 offset-limit 范围内的子切片
func Paginate[T any](list []T, offset, limit int) []T {
	if offset > len(list) {
		return []T{}
	}
	end := offset + limit
	if end > len(list) {
		end = len(list)
	}
	return list[offset:end]
}

// JWTAuthMiddleware Gin中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"msg": "未认证"})
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"msg": "token无效或已过期"})
			return
		}
		c.Next()
	}
}
