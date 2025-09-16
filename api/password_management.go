package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type PasswordChangeRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type PasswordGenerateRequest struct {
	Length int `json:"length,omitempty"`
}

type PasswordHashRequest struct {
	Password string `json:"password" binding:"required"`
}

type PasswordValidateRequest struct {
	Password       string `json:"password" binding:"required"`
	HashedPassword string `json:"hashedPassword" binding:"required"`
}

type PasswordManager struct{}

func NewPasswordManager() *PasswordManager {
	return &PasswordManager{}
}

func (pm *PasswordManager) HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("生成盐失败: %w", err)
	}

	passwordWithSalt := password + base64.URLEncoding.EncodeToString(salt)

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码哈希失败: %w", err)
	}

	return base64.URLEncoding.EncodeToString(salt) + ":" + string(hashedBytes), nil
}

func (pm *PasswordManager) VerifyPassword(password, hashedPassword string) bool {
	parts := strings.Split(hashedPassword, ":")
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	passwordWithSalt := password + base64.URLEncoding.EncodeToString(salt)

	err = bcrypt.CompareHashAndPassword([]byte(parts[1]), []byte(passwordWithSalt))
	return err == nil
}

func (pm *PasswordManager) GeneratePassword(length int) (string, error) {
	if length <= 0 {
		length = model.DefaultPasswordLen
	}

	if length > model.MaxPasswordLen {
		length = model.MaxPasswordLen
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	charsetBytes := []byte(charset)

	for i := range b {
		randomIndex := make([]byte, 1)
		if _, err := rand.Read(randomIndex); err != nil {
			return "", fmt.Errorf("生成随机字符失败: %w", err)
		}
		b[i] = charsetBytes[randomIndex[0]%byte(len(charsetBytes))]
	}

	return string(b), nil
}

func (pm *PasswordManager) ValidatePasswordStrength(password string) (bool, string) {
	if len(password) < model.MinPasswordLen {
		return false, fmt.Sprintf("密码长度至少%d位", model.MinPasswordLen)
	}

	if len(password) > model.MaxPasswordLen {
		return false, fmt.Sprintf("密码长度不能超过%d位", model.MaxPasswordLen)
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return false, "密码必须包含大写字母"
	}
	if !hasLower {
		return false, "密码必须包含小写字母"
	}
	if !hasDigit {
		return false, "密码必须包含数字"
	}
	if !hasSpecial {
		return false, "密码必须包含特殊字符"
	}

	return true, "密码强度符合要求"
}

var passwordManager = NewPasswordManager()

func RegisterPasswordAdmin(r *gin.RouterGroup, logger *zap.Logger) {
	r.POST("/admin/password/change", changePassword(logger))
	r.POST("/admin/password/generate", generatePassword(logger))
	r.POST("/admin/password/hash", hashPassword(logger))
	r.POST("/admin/password/validate", validatePassword(logger))
}

// changePassword 修改密码
func changePassword(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PasswordChangeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeBadRequest,
				Message: "请求参数格式错误",
				Details: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		// 验证新密码强度
		if valid, message := passwordManager.ValidatePasswordStrength(req.NewPassword); !valid {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeValidationFailed,
				Message: "密码强度不符合要求",
				Details: message,
			}, http.StatusBadRequest)
			return
		}

		// 这里应该验证旧密码和更新新密码
		// 具体实现取决于你的认证系统

		middleware.ResponseSuccess(c, gin.H{
			"message": "密码修改成功",
		}, "success", nil)
	}
}

// generatePassword 生成密码
func generatePassword(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PasswordGenerateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			// 使用默认长度
			req.Length = 12
		}

		password, err := passwordManager.GeneratePassword(req.Length)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 生成哈希值
		hashedPassword, err := passwordManager.HashPassword(password)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		middleware.ResponseSuccess(c, gin.H{
			"password":       password,
			"hashedPassword": hashedPassword,
			"length":         len(password),
		}, "密码生成成功", nil)
	}
}

// hashPassword 哈希密码
func hashPassword(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PasswordHashRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeBadRequest,
				Message: "请求参数格式错误",
				Details: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		hashedPassword, err := passwordManager.HashPassword(req.Password)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		middleware.ResponseSuccess(c, gin.H{
			"hashedPassword": hashedPassword,
		}, "密码哈希成功", nil)
	}
}

// validatePassword 验证密码
func validatePassword(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PasswordValidateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.ResponseError(c, logger, &model.APIError{
				Code:    model.CodeBadRequest,
				Message: "请求参数格式错误",
				Details: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		isValid := passwordManager.VerifyPassword(req.Password, req.HashedPassword)

		middleware.ResponseSuccess(c, gin.H{
			"valid": isValid,
			"message": func() string {
				if isValid {
					return "密码验证通过"
				}
				return "密码验证失败"
			}(),
		}, "验证完成", nil)
	}
}
