package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"
	"go.uber.org/zap"
)

var (
	authManager *AuthManager
)

// generateJTI 生成JWT ID
func generateJTI() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为备选
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func InitAuthManager(logger *zap.Logger) {
	if configManager == nil {
		logger.Fatal("配置管理器未初始化")
		return
	}
	authManager = NewAuthManager(logger, configManager)
}

// LoginHandler 登录接口
func LoginHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req model.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeBadRequest,
				Message: model.GetErrorMessage(model.CodeBadRequest),
				Details: "请求参数格式错误",
			}, http.StatusBadRequest)
			return
		}

		// 输入验证和清理
		req.Username = strings.TrimSpace(req.Username)
		req.Password = strings.TrimSpace(req.Password)

		if req.Username == "" || req.Password == "" {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeMissingParameter,
				Message: model.GetErrorMessage(model.CodeMissingParameter),
				Details: "用户名和密码不能为空",
			}, http.StatusBadRequest)
			return
		}

		// 检查用户名长度
		if len(req.Username) > 50 {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeValidationFailed,
				Message: "用户名长度不能超过50个字符",
				Details: "请使用较短的用户名",
			}, http.StatusBadRequest)
			return
		}

		// 检查密码长度
		if len(req.Password) > 128 {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeValidationFailed,
				Message: "密码长度不能超过128个字符",
				Details: "请使用较短的密码",
			}, http.StatusBadRequest)
			return
		}

		if configManager == nil {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeInternalServerError,
				Message: "系统配置未初始化",
			}, http.StatusInternalServerError)
			return
		}

		authConfig := configManager.GetAuthConfig()
		configUsername := authConfig.Username
		password := authConfig.Password
		clientIP := c.ClientIP()
		username := req.Username

		if authManager != nil && authManager.IsLocked(username, clientIP) {
			remainingAttempts := authManager.GetRemainingAttempts(username, clientIP)
			lockTime := authManager.GetLockTime(username, clientIP)
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeRequestTimeout,
				Message: "登录失败次数过多，请稍后再试",
				Details: map[string]interface{}{
					"remainingAttempts": remainingAttempts,
					"maxFailCount":      authConfig.MaxLoginFail,
					"lockDuration":      authConfig.LockDuration.String(),
					"lockTime":          lockTime.String(),
				},
			}, http.StatusTooManyRequests)
			return
		}

		usernameMatch := req.Username == configUsername
		passwordMatch := false

		// 密码验证（移除调试日志与敏感输出）

		if isHashedPassword(password) {
			pm := NewPasswordManager()
			passwordMatch = pm.VerifyPassword(req.Password, password)
		} else {
			passwordMatch = req.Password == password
		}

		if usernameMatch && passwordMatch {
			logger.Info("用户登录成功",
				zap.String("username", req.Username),
				zap.String("clientIP", c.ClientIP()),
				zap.String("userAgent", c.GetHeader("User-Agent")),
				zap.String("event", "login_success"),
			)

			secret := configManager.GetJWTSecret()
			authConfig := configManager.GetAuthConfig()

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": req.Username,
				"iat":      time.Now().Unix(),                                // 签发时间
				"exp":      time.Now().Add(authConfig.SessionTimeout).Unix(), // 从配置读取过期时间
				"iss":      "k8svision",                                      // 签发者
				"aud":      "k8svision-client",                               // 受众
				"jti":      generateJTI(),                                    // JWT ID，用于撤销
			})
			tokenString, err := token.SignedString(secret)
			if err != nil {
				logger.Error("Token生成失败",
					zap.String("username", req.Username),
					zap.Error(err),
				)
				middleware.ResponseError(c, logger, &model.APIError{
					Code:    model.CodeAuthError,
					Message: model.GetErrorMessage(model.CodeAuthError),
					Details: "Token生成失败",
				}, http.StatusInternalServerError)
				return
			}

			logger.Info("JWT token生成成功",
				zap.String("username", req.Username),
			)

			if authManager != nil {
				authManager.RecordSuccess(username, clientIP)
			}

			middleware.ResponseSuccess(c, map[string]string{
				"token": tokenString,
			}, "登录成功", nil)
			return
		}

		if authManager != nil {
			authManager.RecordFailure(username, clientIP)
		}

		remainingAttempts := authConfig.MaxLoginFail
		if authManager != nil {
			remainingAttempts = authManager.GetRemainingAttempts(username, clientIP)
		}

		// 记录登录失败审计日志
		logger.Warn("用户登录失败",
			zap.String("username", req.Username),
			zap.String("clientIP", c.ClientIP()),
			zap.String("userAgent", c.GetHeader("User-Agent")),
			zap.String("event", "login_failed"),
			zap.Int("remainingAttempts", remainingAttempts),
			zap.Bool("usernameMatch", usernameMatch),
		)

		middleware.ResponseError(c, logger, &model.APIError{
			Code:    model.CodeAuthError,
			Message: "用户名或密码错误",
			Details: map[string]interface{}{
				"remainingAttempts": remainingAttempts,
				"maxFailCount":      authConfig.MaxLoginFail,
			},
		}, http.StatusUnauthorized)
	}
}
