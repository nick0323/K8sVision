package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// PasswordUtils 密码工具
type PasswordUtils struct{}

// HashPassword 哈希密码
func (p *PasswordUtils) HashPassword(password string) (string, error) {
	// 生成随机盐
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("生成盐失败: %w", err)
	}

	// 将密码和盐组合
	passwordWithSalt := password + base64.URLEncoding.EncodeToString(salt)

	// 使用bcrypt哈希
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码哈希失败: %w", err)
	}

	// 将盐和哈希值组合存储
	return base64.URLEncoding.EncodeToString(salt) + ":" + string(hashedBytes), nil
}

// VerifyPassword 验证密码
func (p *PasswordUtils) VerifyPassword(password, hashedPassword string) bool {
	// 分离盐和哈希值
	parts := strings.Split(hashedPassword, ":")
	if len(parts) != 2 {
		return false
	}

	salt := parts[0]
	hash := parts[1]

	// 将密码和盐组合
	passwordWithSalt := password + salt

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordWithSalt))
	return err == nil
}

// ValidatePasswordStrength 验证密码强度
func (p *PasswordUtils) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("密码长度至少8位")
	}

	if len(password) > 128 {
		return fmt.Errorf("密码长度不能超过128位")
	}

	// 检查是否包含至少一个字母
	hasLetter := false
	for _, char := range password {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			hasLetter = true
			break
		}
	}
	if !hasLetter {
		return fmt.Errorf("密码必须包含至少一个字母")
	}

	// 检查是否包含至少一个数字
	hasDigit := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return fmt.Errorf("密码必须包含至少一个数字")
	}

	return nil
}

// GenerateSecurePassword 生成安全密码
func (p *PasswordUtils) GenerateSecurePassword(length int) (string, error) {
	if length < 8 {
		length = 12
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)

	for i := range password {
		randomBytes := make([]byte, 1)
		if _, err := rand.Read(randomBytes); err != nil {
			return "", fmt.Errorf("生成随机密码失败: %w", err)
		}
		password[i] = charset[randomBytes[0]%byte(len(charset))]
	}

	return string(password), nil
}

// 全局密码工具实例
var PasswordUtil = &PasswordUtils{}
