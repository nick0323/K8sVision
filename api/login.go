package api

import (
	"net/http"
	"os"
	"strconv"
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
	maxLoginFail  = 5
	lockDuration  = 10 * time.Minute
	configManager *config.Manager
	authManager   *AuthManager
)

// SetConfigManager 设置配置管理器
func SetConfigManager(cm *config.Manager) {
	configManager = cm
	// 从配置管理器获取认证配置
	if cm != nil {
		authConfig := cm.GetAuthConfig()
		maxLoginFail = authConfig.MaxLoginFail
		lockDuration = authConfig.LockDuration
	}
}

// InitAuthManager 初始化认证管理器
func InitAuthManager(logger *zap.Logger) {
	authManager = NewAuthManager(logger)
}

func init() {
	if v := os.Getenv("LOGIN_MAX_FAIL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxLoginFail = n
		}
	}
	if v := os.Getenv("LOGIN_LOCK_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			lockDuration = time.Duration(n) * time.Minute
		}
	}
}

// LoginHandler 登录接口
// @Summary 用户登录
// @Description 登录获取 JWT Token，连续失败5次后10分钟内禁止尝试
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body model.LoginRequest true "登录参数"
// @Success 200 {object} model.APIResponse "登录成功"
// @Failure 400 {object} model.APIResponse "参数错误"
// @Failure 401 {object} model.APIResponse "用户名或密码错误"
// @Failure 429 {object} model.APIResponse "登录失败次数过多"
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	// 创建临时logger用于错误处理
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseError(c, logger, &model.APIError{
			Code:    model.CodeBadRequest,
			Message: model.GetErrorMessage(model.CodeBadRequest),
			Details: "请求参数格式错误",
		}, http.StatusBadRequest)
		return
	}

	// 参数验证
	if req.Username == "" || req.Password == "" {
		middleware.ResponseError(c, logger, &model.APIError{
			Code:    model.CodeMissingParameter,
			Message: model.GetErrorMessage(model.CodeMissingParameter),
			Details: "用户名和密码不能为空",
		}, http.StatusBadRequest)
		return
	}

	// 优先从配置管理器获取认证信息
	var username, password string
	if configManager != nil {
		authConfig := configManager.GetAuthConfig()
		username = authConfig.Username
		password = authConfig.Password
	} else {
		// 兼容旧版本，从环境变量获取
		username = os.Getenv("LOGIN_USERNAME")
		password = os.Getenv("LOGIN_PASSWORD")
	}
	key := req.Username + "|" + c.ClientIP()

	// 检查是否被锁定
	if authManager != nil && authManager.IsLocked(key, maxLoginFail, lockDuration) {
		remainingAttempts := authManager.GetRemainingAttempts(key, maxLoginFail)
		middleware.ResponseError(c, logger, &model.APIError{
			Code:    model.CodeRequestTimeout,
			Message: "登录失败次数过多，请稍后再试",
			Details: map[string]interface{}{
				"remainingAttempts": remainingAttempts,
				"maxFailCount":      maxLoginFail,
				"lockDuration":      lockDuration.String(),
			},
		}, http.StatusTooManyRequests)
		return
	}

	// 验证用户名和密码
	usernameMatch := req.Username == username
	passwordMatch := false

	// 检查密码是否为哈希格式（包含冒号分隔符）
	if strings.Contains(password, ":") {
		// 使用哈希验证
		passwordMatch = PasswordUtil.VerifyPassword(req.Password, password)
	} else {
		// 兼容旧版本明文密码
		passwordMatch = req.Password == password
	}

	if usernameMatch && passwordMatch {
		logger.Info("用户登录成功，生成JWT token",
			zap.String("username", req.Username),
			zap.String("clientIP", c.ClientIP()),
		)

		// 获取JWT密钥
		var secret []byte
		if configManager != nil {
			secret = configManager.GetJWTSecret()
		} else {
			// 兼容旧版本
			secret = jwtSecret
		}

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
			zap.String("token", tokenString[:min(len(tokenString), 20)]+"..."),
		)

		// 登录成功，清除失败记录
		if authManager != nil {
			authManager.ClearLoginFailures(key)
		}

		middleware.ResponseSuccess(c, map[string]string{
			"token": tokenString,
		}, "登录成功", nil)
		return
	}

	// 登录失败，记录失败次数
	if authManager != nil {
		authManager.RecordLoginFailure(key, maxLoginFail, lockDuration)
	}

	remainingAttempts := maxLoginFail
	if authManager != nil {
		remainingAttempts = authManager.GetRemainingAttempts(key, maxLoginFail)
	}

	middleware.ResponseError(c, logger, &model.APIError{
		Code:    model.CodeAuthError,
		Message: "用户名或密码错误",
		Details: map[string]interface{}{
			"remainingAttempts": remainingAttempts,
			"lockTime":          lockDuration,
		},
	}, http.StatusUnauthorized)
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
