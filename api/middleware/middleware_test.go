package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nick0323/K8sVision/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTraceMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// 创建测试路由
	r := gin.New()
	r.Use(TraceMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// 测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Trace-ID"))
	
	// 验证traceId在context中
	assert.Contains(t, w.Header().Get("X-Trace-ID"), "-")
}

func TestResponseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// 创建测试路由
	r := gin.New()
	r.Use(TraceMiddleware())
	r.GET("/error", func(c *gin.Context) {
		apiError := &model.APIError{
			Code:    model.CodeBadRequest,
			Message: "测试错误",
			Details: "错误详情",
		}
		ResponseError(c, nil, apiError, http.StatusBadRequest)
	})

	// 测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response model.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, model.CodeBadRequest, response.Code)
	assert.Equal(t, "测试错误", response.Message)
	assert.NotEmpty(t, response.TraceID)
}

func TestResponseSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// 创建测试路由
	r := gin.New()
	r.Use(TraceMiddleware())
	r.GET("/success", func(c *gin.Context) {
		data := map[string]string{"key": "value"}
		ResponseSuccess(c, data, "操作成功", nil)
	})

	// 测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/success", nil)
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response model.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, model.CodeSuccess, response.Code)
	assert.Equal(t, "操作成功", response.Message)
	assert.NotEmpty(t, response.TraceID)
	assert.NotNil(t, response.Data)
}

func TestValidationError(t *testing.T) {
	error := ValidationError("username", "用户名不能为空")
	assert.Equal(t, model.CodeValidationFailed, error.Code)
	assert.Equal(t, model.GetErrorMessage(model.CodeValidationFailed), error.Message)
	
	details, ok := error.Details.(map[string]string)
	assert.True(t, ok)
	assert.Equal(t, "username", details["field"])
	assert.Equal(t, "用户名不能为空", details["message"])
}

func TestMissingParameterError(t *testing.T) {
	error := MissingParameterError("password")
	assert.Equal(t, model.CodeMissingParameter, error.Code)
	assert.Equal(t, model.GetErrorMessage(model.CodeMissingParameter), error.Message)
	
	details, ok := error.Details.(map[string]string)
	assert.True(t, ok)
	assert.Equal(t, "password", details["parameter"])
}

func TestInvalidParameterError(t *testing.T) {
	error := InvalidParameterError("age", "abc")
	assert.Equal(t, model.CodeInvalidParameter, error.Code)
	assert.Equal(t, model.GetErrorMessage(model.CodeInvalidParameter), error.Message)
	
	details, ok := error.Details.(map[string]string)
	assert.True(t, ok)
	assert.Equal(t, "age", details["parameter"])
	assert.Equal(t, "abc", details["value"])
}

func TestLoggingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// 创建测试logger
	logger, _ := zap.NewDevelopment()
	
	// 创建测试路由
	r := gin.New()
	r.Use(LoggingMiddleware(logger))
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// 测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, 200, w.Code)
}

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// 创建测试logger
	logger, _ := zap.NewDevelopment()
	
	// 创建测试路由
	r := gin.New()
	r.Use(ErrorHandler(logger))
	r.Use(TraceMiddleware())
	r.GET("/panic", func(c *gin.Context) {
		panic("测试panic")
	})

	// 测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	var response model.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, model.CodeInternalServerError, response.Code)
	assert.NotEmpty(t, response.TraceID)
} 