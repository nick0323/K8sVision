package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		traceId := getTraceID(c)

		logger.Error("panic recovered",
			zap.String("traceId", traceId),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Any("error", recovered),
			zap.String("stack", string(debug.Stack())),
		)

		ResponseError(c, logger, &model.APIError{
			Code:    model.CodeInternalServerError,
			Message: "服务器内部错误",
			Details: "系统发生未知错误",
		}, http.StatusInternalServerError)
	})
}

// ErrorHandlerWithMetrics 带指标记录的错误处理中间件
func ErrorHandlerWithMetrics(logger *zap.Logger, metricsRecorder interface{}) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		traceId := getTraceID(c)

		logger.Error("panic recovered",
			zap.String("traceId", traceId),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Any("error", recovered),
			zap.String("stack", string(debug.Stack())),
		)

		// 记录错误指标
		if recorder, ok := metricsRecorder.(interface {
			RecordError(err string)
		}); ok {
			recorder.RecordError(fmt.Sprintf("panic: %v", recovered))
		}

		ResponseError(c, logger, &model.APIError{
			Code:    model.CodeInternalServerError,
			Message: "服务器内部错误",
			Details: "系统发生未知错误",
		}, http.StatusInternalServerError)
	})
}

// APIError 自定义API错误结构
type APIError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ResponseError 统一错误响应处理
func ResponseError(c *gin.Context, logger *zap.Logger, err error, httpCode int) {
	traceId := getTraceID(c)

	var apiError *model.APIError

	switch e := err.(type) {
	case *model.APIError:
		apiError = e
	case *errors.StatusError:
		apiError = convertK8sError(e)
	default:
		apiError = &model.APIError{
			Code:    model.CodeInternalServerError,
			Message: err.Error(),
		}
	}

	if logger != nil {
		if httpCode >= 500 {
			logger.Error("server error",
				zap.String("traceId", traceId),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("httpCode", httpCode),
				zap.Int("errorCode", apiError.Code),
				zap.String("errorMessage", apiError.Message),
				zap.Error(err),
			)
		} else {
			logger.Warn("client error",
				zap.String("traceId", traceId),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("httpCode", httpCode),
				zap.Int("errorCode", apiError.Code),
				zap.String("errorMessage", apiError.Message),
			)
		}
	}

	c.JSON(httpCode, model.APIResponse{
		Code:      apiError.Code,
		Message:   apiError.Message,
		Data:      apiError.Details,
		TraceID:   traceId,
		Timestamp: time.Now().Unix(),
	})
}

// ResponseSuccess 统一成功响应处理
func ResponseSuccess(c *gin.Context, data interface{}, message string, page *model.PageMeta) {
	traceId := getTraceID(c)

	c.JSON(http.StatusOK, model.APIResponse{
		Code:      model.CodeSuccess,
		Message:   message,
		Data:      data,
		TraceID:   traceId,
		Timestamp: time.Now().Unix(),
		Page:      page,
	})
}

// convertK8sError 转换Kubernetes错误为API错误
func convertK8sError(k8sErr *errors.StatusError) *model.APIError {
	switch k8sErr.Status().Code {
	case http.StatusNotFound:
		return &model.APIError{
			Code:    model.CodeResourceNotFound,
			Message: model.GetErrorMessage(model.CodeResourceNotFound),
			Details: k8sErr.Error(),
		}
	case http.StatusConflict:
		return &model.APIError{
			Code:    model.CodeResourceExists,
			Message: model.GetErrorMessage(model.CodeResourceExists),
			Details: k8sErr.Error(),
		}
	case http.StatusForbidden:
		return &model.APIError{
			Code:    model.CodePermissionDenied,
			Message: model.GetErrorMessage(model.CodePermissionDenied),
			Details: k8sErr.Error(),
		}
	case http.StatusUnauthorized:
		return &model.APIError{
			Code:    model.CodeUnauthorized,
			Message: model.GetErrorMessage(model.CodeUnauthorized),
			Details: k8sErr.Error(),
		}
	case http.StatusBadRequest:
		return &model.APIError{
			Code:    model.CodeBadRequest,
			Message: model.GetErrorMessage(model.CodeBadRequest),
			Details: k8sErr.Error(),
		}
	default:
		return &model.APIError{
			Code:    model.CodeK8sAPIError,
			Message: model.GetErrorMessage(model.CodeK8sAPIError),
			Details: k8sErr.Error(),
		}
	}
}

// getTraceID 从gin.Context获取traceId
func getTraceID(c *gin.Context) string {
	if v, ok := c.Get("traceId"); ok {
		if tid, ok := v.(string); ok {
			return tid
		}
	}
	return ""
}

// ValidationError 参数验证错误
func ValidationError(field, message string) *model.APIError {
	return &model.APIError{
		Code:    model.CodeValidationFailed,
		Message: model.GetErrorMessage(model.CodeValidationFailed),
		Details: map[string]string{
			"field":   field,
			"message": message,
		},
	}
}

// MissingParameterError 缺少参数错误
func MissingParameterError(param string) *model.APIError {
	return &model.APIError{
		Code:    model.CodeMissingParameter,
		Message: model.GetErrorMessage(model.CodeMissingParameter),
		Details: map[string]string{
			"parameter": param,
		},
	}
}

// InvalidParameterError 无效参数错误
func InvalidParameterError(param, value string) *model.APIError {
	return &model.APIError{
		Code:    model.CodeInvalidParameter,
		Message: model.GetErrorMessage(model.CodeInvalidParameter),
		Details: map[string]string{
			"parameter": param,
			"value":     value,
		},
	}
}
