package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nick0323/K8sVision/config"
	"github.com/nick0323/K8sVision/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建测试logger
	logger, _ := zap.NewDevelopment()

	// 创建测试配置管理器
	configMgr := config.NewManager(logger)
	configMgr.Load("") // 加载默认配置
	SetConfigManager(configMgr)
	InitAuthManager(logger)

	// 创建测试路由
	r := gin.New()
	r.POST("/login", LoginHandler)

	t.Run("成功登录", func(t *testing.T) {
		reqBody := model.LoginRequest{
			Username: "admin",
			Password: "admin",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, model.CodeSuccess, response.Code)
		assert.Contains(t, response.Data, "token")
	})

	t.Run("用户名错误", func(t *testing.T) {
		reqBody := model.LoginRequest{
			Username: "wronguser",
			Password: "admin",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response model.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, model.CodeAuthError, response.Code)
	})

	t.Run("密码错误", func(t *testing.T) {
		reqBody := model.LoginRequest{
			Username: "admin",
			Password: "wrongpass",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response model.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, model.CodeAuthError, response.Code)
	})

	t.Run("缺少用户名", func(t *testing.T) {
		reqBody := model.LoginRequest{
			Password: "admin",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response model.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, model.CodeMissingParameter, response.Code)
	})

	t.Run("缺少密码", func(t *testing.T) {
		reqBody := model.LoginRequest{
			Username: "admin",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response model.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, model.CodeMissingParameter, response.Code)
	})

	t.Run("无效JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response model.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, model.CodeBadRequest, response.Code)
	})
}

func TestPasswordUtils(t *testing.T) {
	t.Run("密码哈希和验证", func(t *testing.T) {
		password := "testpassword123"

		// 哈希密码
		hashed, err := PasswordUtil.HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, password, hashed)

		// 验证密码
		valid := PasswordUtil.VerifyPassword(password, hashed)
		assert.True(t, valid)

		// 验证错误密码
		invalid := PasswordUtil.VerifyPassword("wrongpassword", hashed)
		assert.False(t, invalid)
	})

	t.Run("密码强度验证", func(t *testing.T) {
		// 有效密码
		err := PasswordUtil.ValidatePasswordStrength("password123")
		assert.NoError(t, err)

		// 密码太短
		err = PasswordUtil.ValidatePasswordStrength("pass1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "至少8位")

		// 密码太长
		longPassword := make([]byte, 200)
		for i := range longPassword {
			longPassword[i] = 'a'
		}
		err = PasswordUtil.ValidatePasswordStrength(string(longPassword))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不能超过128位")

		// 缺少字母
		err = PasswordUtil.ValidatePasswordStrength("12345678")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "至少一个字母")

		// 缺少数字
		err = PasswordUtil.ValidatePasswordStrength("password")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "至少一个数字")
	})

	t.Run("生成安全密码", func(t *testing.T) {
		password, err := PasswordUtil.GenerateSecurePassword(12)
		assert.NoError(t, err)
		assert.Len(t, password, 12)

		// 验证生成的密码强度
		err = PasswordUtil.ValidatePasswordStrength(password)
		assert.NoError(t, err)
	})
}
