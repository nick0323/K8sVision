package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nick0323/K8sVision/model"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
)

// ErrorHandler 错误处理器接口
type ErrorHandler interface {
	Handle(ctx context.Context, err error) *model.APIError
	GetHTTPStatus(err error) int
}

// DefaultErrorHandler 默认错误处理器
type DefaultErrorHandler struct {
	logger *zap.Logger
}

// NewDefaultErrorHandler 创建默认错误处理器
func NewDefaultErrorHandler(logger *zap.Logger) *DefaultErrorHandler {
	return &DefaultErrorHandler{
		logger: logger,
	}
}

// Handle 处理错误
func (h *DefaultErrorHandler) Handle(ctx context.Context, err error) *model.APIError {
	if err == nil {
		return nil
	}

	// 记录错误日志
	h.logger.Error("处理错误",
		zap.Error(err),
		zap.String("traceId", getTraceIDFromContext(ctx)),
	)

	// 根据错误类型进行分类处理
	switch e := err.(type) {
	case *model.APIError:
		return e
	case *errors.StatusError:
		return h.handleK8sError(e)
	default:
		return h.handleGenericError(err)
	}
}

// GetHTTPStatus 获取HTTP状态码
func (h *DefaultErrorHandler) GetHTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch e := err.(type) {
	case *model.APIError:
		return h.getHTTPStatusFromCode(e.Code)
	case *errors.StatusError:
		return int(e.Status().Code)
	default:
		return http.StatusInternalServerError
	}
}

// handleK8sError 处理Kubernetes错误
func (h *DefaultErrorHandler) handleK8sError(k8sErr *errors.StatusError) *model.APIError {
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

// handleGenericError 处理通用错误
func (h *DefaultErrorHandler) handleGenericError(err error) *model.APIError {
	return &model.APIError{
		Code:    model.CodeInternalServerError,
		Message: model.GetErrorMessage(model.CodeInternalServerError),
		Details: err.Error(),
	}
}

// getHTTPStatusFromCode 根据错误码获取HTTP状态码
func (h *DefaultErrorHandler) getHTTPStatusFromCode(code int) int {
	switch {
	case code >= 1000 && code < 2000:
		return http.StatusBadRequest
	case code >= 2000 && code < 3000:
		return http.StatusInternalServerError
	case code >= 3000 && code < 4000:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// getTraceIDFromContext 从上下文获取追踪ID
func getTraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if traceID, ok := ctx.Value("traceId").(string); ok {
		return traceID
	}

	return ""
}

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	errorHandler := NewDefaultErrorHandler(logger)

	return func(c *gin.Context) {
		c.Next()

		// 处理gin.Context中的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			apiError := errorHandler.Handle(c.Request.Context(), err.Err)
			httpStatus := errorHandler.GetHTTPStatus(err.Err)

			c.JSON(httpStatus, model.APIResponse{
				Code:      apiError.Code,
				Message:   apiError.Message,
				Data:      apiError.Details,
				TraceID:   getTraceIDFromContext(c.Request.Context()),
				Timestamp: time.Now().Unix(),
			})
		}
	}
}

// ValidationError 创建验证错误
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

// MissingParameterError 创建缺少参数错误
func MissingParameterError(param string) *model.APIError {
	return &model.APIError{
		Code:    model.CodeMissingParameter,
		Message: model.GetErrorMessage(model.CodeMissingParameter),
		Details: map[string]string{
			"parameter": param,
		},
	}
}

// InvalidParameterError 创建无效参数错误
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

// ResourceNotFoundError 创建资源不存在错误
func ResourceNotFoundError(resourceType, name string) *model.APIError {
	return &model.APIError{
		Code:    model.CodeResourceNotFound,
		Message: model.GetErrorMessage(model.CodeResourceNotFound),
		Details: map[string]string{
			"resourceType": resourceType,
			"name":         name,
		},
	}
}

// PermissionDeniedError 创建权限不足错误
func PermissionDeniedError(action, resource string) *model.APIError {
	return &model.APIError{
		Code:    model.CodePermissionDenied,
		Message: model.GetErrorMessage(model.CodePermissionDenied),
		Details: map[string]string{
			"action":   action,
			"resource": resource,
		},
	}
}

// K8sClientError 创建K8s客户端错误
func K8sClientError(operation string, err error) *model.APIError {
	return &model.APIError{
		Code:    model.CodeK8sClientError,
		Message: model.GetErrorMessage(model.CodeK8sClientError),
		Details: map[string]string{
			"operation": operation,
			"error":     err.Error(),
		},
	}
}

// K8sAPIError 创建K8s API错误
func K8sAPIError(operation string, err error) *model.APIError {
	return &model.APIError{
		Code:    model.CodeK8sAPIError,
		Message: model.GetErrorMessage(model.CodeK8sAPIError),
		Details: map[string]string{
			"operation": operation,
			"error":     err.Error(),
		},
	}
}
