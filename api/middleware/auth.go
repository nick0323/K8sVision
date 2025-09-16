package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nick0323/K8sVision/model"
	"go.uber.org/zap"
)

type ConfigProvider interface {
	GetJWTSecret() []byte
}

func getJWTSecret(provider ConfigProvider) []byte {
	if provider == nil {
		panic("配置提供者未初始化")
	}
	return provider.GetJWTSecret()
}

func JWTAuthMiddleware(logger *zap.Logger, configProvider ConfigProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetString("traceId")

		logger.Debug("JWT认证开始",
			zap.String("traceId", traceId),
			zap.String("clientIP", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

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

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		logger.Debug("Parsing JWT token",
			zap.String("traceId", traceId),
			zap.String("token", tokenStr[:min(len(tokenStr), 20)]+"..."),
			zap.Int("tokenLength", len(tokenStr)),
		)

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

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return getJWTSecret(configProvider), nil
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
