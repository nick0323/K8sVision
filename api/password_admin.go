package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"
	"go.uber.org/zap"
)

// PasswordChangeRequest 密码修改请求
type PasswordChangeRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// PasswordGenerateRequest 密码生成请求
type PasswordGenerateRequest struct {
	Length int `json:"length,omitempty"`
}

// PasswordHashRequest 密码哈希请求
type PasswordHashRequest struct {
	Password string `json:"password" binding:"required"`
}

// RegisterPasswordAdmin 注册密码管理相关路由
func RegisterPasswordAdmin(r *gin.RouterGroup, logger *zap.Logger) {
	r.POST("/admin/password/change", changePassword(logger))
	r.POST("/admin/password/generate", generatePassword(logger))
	r.POST("/admin/password/hash", hashPassword(logger))
	r.POST("/admin/password/validate", validatePassword(logger))
}

// changePassword 修改密码
// @Summary 修改密码
// @Description 修改管理员密码
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body PasswordChangeRequest true "密码修改请求"
// @Success 200 {object} model.APIResponse
// @Router /admin/password/change [post]
func changePassword(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PasswordChangeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeBadRequest,
				Message: "请求参数错误",
				Details: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		// 验证新密码强度
		if err := PasswordUtil.ValidatePasswordStrength(req.NewPassword); err != nil {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeValidationFailed,
				Message: "密码强度不符合要求",
				Details: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		// 获取当前密码
		var currentPassword string
		if configManager != nil {
			authConfig := configManager.GetAuthConfig()
			currentPassword = authConfig.Password
		} else {
			currentPassword = "admin" // 默认密码
		}

		// 验证旧密码
		oldPasswordMatch := false
		if strings.Contains(currentPassword, ":") {
			oldPasswordMatch = PasswordUtil.VerifyPassword(req.OldPassword, currentPassword)
		} else {
			oldPasswordMatch = req.OldPassword == currentPassword
		}

		if !oldPasswordMatch {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeAuthError,
				Message: "旧密码错误",
			}, http.StatusUnauthorized)
			return
		}

		// 生成新密码哈希
		hashedPassword, err := PasswordUtil.HashPassword(req.NewPassword)
		if err != nil {
			logger.Error("密码哈希失败", zap.Error(err))
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeInternalServerError,
				Message: "密码处理失败",
			}, http.StatusInternalServerError)
			return
		}

		// 这里应该更新配置文件或数据库中的密码
		// 由于当前使用配置文件，这里只是返回哈希值供手动更新
		logger.Info("密码修改请求", zap.String("username", "admin"))

		middleware.ResponseSuccess(c, map[string]string{
			"message":        "密码修改成功，请将以下哈希值更新到配置文件中",
			"hashedPassword": hashedPassword,
		}, "密码修改成功", nil)
	}
}

// generatePassword 生成安全密码
// @Summary 生成安全密码
// @Description 生成指定长度的安全密码
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body PasswordGenerateRequest true "密码生成请求"
// @Success 200 {object} model.APIResponse
// @Router /admin/password/generate [post]
func generatePassword(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PasswordGenerateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			req.Length = 12 // 默认长度
		}

		password, err := PasswordUtil.GenerateSecurePassword(req.Length)
		if err != nil {
			logger.Error("密码生成失败", zap.Error(err))
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeInternalServerError,
				Message: "密码生成失败",
			}, http.StatusInternalServerError)
			return
		}

		// 同时生成哈希值
		hashedPassword, err := PasswordUtil.HashPassword(password)
		if err != nil {
			logger.Error("密码哈希失败", zap.Error(err))
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeInternalServerError,
				Message: "密码处理失败",
			}, http.StatusInternalServerError)
			return
		}

		middleware.ResponseSuccess(c, map[string]string{
			"password":       password,
			"hashedPassword": hashedPassword,
		}, "密码生成成功", nil)
	}
}

// hashPassword 哈希密码
// @Summary 哈希密码
// @Description 将明文密码转换为哈希值
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body PasswordHashRequest true "密码哈希请求"
// @Success 200 {object} model.APIResponse
// @Router /admin/password/hash [post]
func hashPassword(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PasswordHashRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeBadRequest,
				Message: "请求参数错误",
				Details: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		hashedPassword, err := PasswordUtil.HashPassword(req.Password)
		if err != nil {
			logger.Error("密码哈希失败", zap.Error(err))
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeInternalServerError,
				Message: "密码处理失败",
			}, http.StatusInternalServerError)
			return
		}

		middleware.ResponseSuccess(c, map[string]string{
			"hashedPassword": hashedPassword,
		}, "密码哈希成功", nil)
	}
}

// validatePassword 验证密码强度
// @Summary 验证密码强度
// @Description 验证密码是否符合强度要求
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body PasswordHashRequest true "密码验证请求"
// @Success 200 {object} model.APIResponse
// @Router /admin/password/validate [post]
func validatePassword(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PasswordHashRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeBadRequest,
				Message: "请求参数错误",
				Details: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		err := PasswordUtil.ValidatePasswordStrength(req.Password)
		if err != nil {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeValidationFailed,
				Message: "密码强度不符合要求",
				Details: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		middleware.ResponseSuccess(c, map[string]string{
			"message": "密码强度符合要求",
		}, "密码验证成功", nil)
	}
}
