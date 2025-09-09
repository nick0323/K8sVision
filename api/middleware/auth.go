package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nick0323/K8sVision/config"
	"github.com/nick0323/K8sVision/model"
	"go.uber.org/zap"
)

var configManager *config.Manager

// SetConfigManager 设置配置管理器
func SetConfigManager(cm *config.Manager) {
	configManager = cm
}

// InitJWTSecret 初始化 JWT 密钥（保持向后兼容）
func InitJWTSecret(secret string) {
	// 这个函数现在只是占位符，实际的密钥管理通过配置管理器
	_ = secret
}

// getJWTSecret 获取JWT密钥
func getJWTSecret() []byte {
	// 优先从配置管理器获取密钥
	if configManager != nil {
		return configManager.GetJWTSecret()
	}

	// 兼容旧版本，从环境变量获取密钥
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "k8svision-secret-key"
	}
	return []byte(secret)
}

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := getTraceID(c)

		logger.Debug("JWT认证开始",
			zap.String("traceId", traceId),
			zap.String("clientIP", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		// 获取Authorization头
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			logger.Warn("missing authorization header",
				zap.String("traceId", traceId),
				zap.String("clientIP", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
			)
			ResponseError(c, logger, &model.APIError{
				Code:    model.CodeUnauthorized,
				Message: model.GetErrorMessage(model.CodeUnauthorized),
				Details: "缺少Authorization头",
			}, 401)
			c.Abort()
			return
		}

		logger.Debug("Authorization header found",
			zap.String("traceId", traceId),
			zap.String("header", tokenStr[:min(len(tokenStr), 20)]+"..."),
		)

		// 检查Bearer前缀
		if !strings.HasPrefix(tokenStr, "Bearer ") {
			logger.Warn("invalid authorization format",
				zap.String("traceId", traceId),
				zap.String("clientIP", c.ClientIP()),
				zap.String("header", tokenStr),
			)
			ResponseError(c, logger, &model.APIError{
				Code:    model.CodeUnauthorized,
				Message: model.GetErrorMessage(model.CodeUnauthorized),
				Details: "Authorization格式错误，应为Bearer token",
			}, 401)
			c.Abort()
			return
		}

		// 提取token
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		logger.Debug("Parsing JWT token",
			zap.String("traceId", traceId),
			zap.String("token", tokenStr[:min(len(tokenStr), 20)]+"..."),
			zap.Int("tokenLength", len(tokenStr)),
		)

		// 检查token格式
		segments := strings.Split(tokenStr, ".")
		if len(segments) != 3 {
			logger.Warn("invalid token format - wrong number of segments",
				zap.String("traceId", traceId),
				zap.String("clientIP", c.ClientIP()),
				zap.Int("segments", len(segments)),
				zap.String("token", tokenStr),
			)
			ResponseError(c, logger, &model.APIError{
				Code:    model.CodeAuthError,
				Message: model.GetErrorMessage(model.CodeAuthError),
				Details: fmt.Sprintf("Token格式错误：期望3个段，实际%d个段", len(segments)),
			}, 401)
			c.Abort()
			return
		}

		// 解析JWT token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return getJWTSecret(), nil
		})

		if err != nil {
			logger.Warn("jwt parse error",
				zap.String("traceId", traceId),
				zap.String("clientIP", c.ClientIP()),
				zap.Error(err),
			)
			ResponseError(c, logger, &model.APIError{
				Code:    model.CodeAuthError,
				Message: model.GetErrorMessage(model.CodeAuthError),
				Details: "Token解析失败: " + err.Error(),
			}, 401)
			c.Abort()
			return
		}

		// 验证token有效性
		if !token.Valid {
			logger.Warn("invalid jwt token",
				zap.String("traceId", traceId),
				zap.String("clientIP", c.ClientIP()),
			)
			ResponseError(c, logger, &model.APIError{
				Code:    model.CodeAuthError,
				Message: model.GetErrorMessage(model.CodeAuthError),
				Details: "Token无效或已过期",
			}, 401)
			c.Abort()
			return
		}

		// 提取用户信息
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if username, exists := claims["username"].(string); exists {
				c.Set("username", username)
				logger.Debug("User authenticated",
					zap.String("traceId", traceId),
					zap.String("username", username),
				)
			}
		}

		logger.Info("authentication successful",
			zap.String("traceId", traceId),
			zap.String("clientIP", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}

// RequirePermission 权限检查中间件
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以实现具体的权限检查逻辑
		// 目前只是占位符，可以根据实际需求扩展
		username, exists := c.Get("username")
		if !exists {
			c.AbortWithStatus(401)
			return
		}

		// 简单的权限检查示例
		// 在实际应用中，这里应该查询数据库或缓存来验证用户权限
		if permission != "" {
			// 这里可以添加具体的权限验证逻辑
			_ = username // 避免未使用变量警告
		}

		c.Next()
	}
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
