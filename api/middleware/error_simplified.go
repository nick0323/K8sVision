package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nick0323/K8sVision/model"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
)

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		traceID := c.GetString("traceId")

		logger.Error("panic recovered",
			zap.String("traceId", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Any("error", recovered),
			zap.String("stack", string(debug.Stack())),
		)

		ResponseError(c, logger, &model.APIError{
			Code:    model.CodeInternalServerError,
			Message: "服务器内部错误",
			Details: fmt.Sprintf("panic: %v", recovered),
		}, http.StatusInternalServerError)
	})
}

func ResponseError(c *gin.Context, logger *zap.Logger, err error, httpCode int) {
	traceID := c.GetString("traceId")

	var apiError *model.APIError

	switch e := err.(type) {
	case *model.APIError:
		apiError = e
	case *errors.StatusError:
		apiError = ConvertK8sError(e)
	default:
		apiError = &model.APIError{
			Code:    model.CodeInternalServerError,
			Message: model.GetErrorMessage(model.CodeInternalServerError),
			Details: e.Error(),
		}
	}

	logFields := []zap.Field{
		zap.String("traceId", traceID),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("httpCode", httpCode),
		zap.Int("errorCode", apiError.Code),
		zap.String("errorMessage", apiError.Message),
	}

	if httpCode >= 500 {
		logger.Error("server error", append(logFields, zap.Error(err))...)
	} else {
		logger.Warn("client error", logFields...)
	}

	c.JSON(httpCode, model.APIResponse{
		Code:      apiError.Code,
		Message:   apiError.Message,
		Data:      apiError.Details,
		TraceID:   traceID,
		Timestamp: time.Now().Unix(),
	})
}

func ResponseSuccess(c *gin.Context, data interface{}, message string, page *model.PageMeta) {
	c.JSON(http.StatusOK, model.APIResponse{
		Code:      model.CodeSuccess,
		Message:   message,
		Data:      data,
		TraceID:   c.GetString("traceId"),
		Timestamp: time.Now().Unix(),
		Page:      page,
	})
}

func ConvertK8sError(k8sErr *errors.StatusError) *model.APIError {
	errorMappings := map[int32]int{
		http.StatusNotFound:           model.CodeResourceNotFound,
		http.StatusConflict:           model.CodeResourceExists,
		http.StatusForbidden:          model.CodePermissionDenied,
		http.StatusUnauthorized:       model.CodeUnauthorized,
		http.StatusBadRequest:         model.CodeBadRequest,
		http.StatusRequestTimeout:     model.CodeRequestTimeout,
		http.StatusServiceUnavailable: model.CodeServiceUnavailable,
	}

	code, exists := errorMappings[k8sErr.Status().Code]
	if !exists {
		code = model.CodeK8sAPIError
	}

	return &model.APIError{
		Code:    code,
		Message: model.GetErrorMessage(code),
		Details: map[string]interface{}{
			"k8sError":   k8sErr.Error(),
			"statusCode": k8sErr.Status().Code,
			"reason":     k8sErr.Status().Reason,
		},
	}
}

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

func MissingParameterError(param string) *model.APIError {
	return &model.APIError{
		Code:    model.CodeMissingParameter,
		Message: model.GetErrorMessage(model.CodeMissingParameter),
		Details: map[string]string{
			"parameter": param,
		},
	}
}

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
