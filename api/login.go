package api

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nick0323/K8sVision/model"
	"github.com/nick0323/K8sVision/api/middleware"
	"go.uber.org/zap"
)

var loginFailMap = make(map[string]struct {
	Count    int
	LastFail time.Time
})

var (
	maxLoginFail = 5
	lockDuration = 10 * time.Minute
)

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

	username := os.Getenv("LOGIN_USERNAME")
	password := os.Getenv("LOGIN_PASSWORD")
	key := req.Username + "|" + c.ClientIP()
	failInfo := loginFailMap[key]
	
	// 超过锁定时间自动重置失败次数
	if time.Since(failInfo.LastFail) >= lockDuration {
		failInfo.Count = 0
	}
	
	if failInfo.Count >= maxLoginFail {
		middleware.ResponseError(c, logger, &model.APIError{
			Code:    model.CodeRequestTimeout,
			Message: "登录失败次数过多，请10分钟后再试",
			Details: map[string]interface{}{
				"remainingTime": lockDuration - time.Since(failInfo.LastFail),
				"maxFailCount":  maxLoginFail,
			},
		}, http.StatusTooManyRequests)
		return
	}
	
	if req.Username == username && req.Password == password {
		logger.Info("用户登录成功，生成JWT token",
			zap.String("username", req.Username),
			zap.String("clientIP", c.ClientIP()),
		)
		
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": req.Username,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		})
		tokenString, err := token.SignedString(jwtSecret)
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
		delete(loginFailMap, key)
		
		middleware.ResponseSuccess(c, map[string]string{
			"token": tokenString,
		}, "登录成功", nil)
		return
	}
	
	// 登录失败
	loginFailMap[key] = struct {
		Count    int
		LastFail time.Time
	}{failInfo.Count + 1, time.Now()}
	
	middleware.ResponseError(c, logger, &model.APIError{
		Code:    model.CodeAuthError,
		Message: "用户名或密码错误",
		Details: map[string]interface{}{
			"remainingAttempts": maxLoginFail - (failInfo.Count + 1),
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
