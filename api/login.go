package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/config"
	"github.com/nick0323/K8sVision/model"
	"go.uber.org/zap"
)

var (
	configManager *config.Manager
	authManager   *AuthManager
)

func SetConfigManager(cm *config.Manager) {
	configManager = cm
	fmt.Printf("DEBUG: SetConfigManager called, configManager is nil: %v\n", configManager == nil)
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

		if req.Username == "" || req.Password == "" {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeMissingParameter,
				Message: model.GetErrorMessage(model.CodeMissingParameter),
				Details: "用户名和密码不能为空",
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

		if strings.Contains(password, ":") {
			pm := NewPasswordManager()
			passwordMatch = pm.VerifyPassword(req.Password, password)
		} else {
			passwordMatch = req.Password == password
		}

		if usernameMatch && passwordMatch {
			logger.Info("用户登录成功，生成JWT token",
				zap.String("username", req.Username),
				zap.String("clientIP", c.ClientIP()),
			)

			secret := configManager.GetJWTSecret()

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": req.Username,
				"exp":      time.Now().Add(24 * time.Hour).Unix(),
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
				zap.String("token", func() string {
					if len(tokenString) > 20 {
						return tokenString[:20] + "..."
					}
					return tokenString
				}()),
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
