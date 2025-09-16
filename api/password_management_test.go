package api

import (
	"testing"

	"github.com/nick0323/K8sVision/model"
)

// TestPasswordManager_HashPassword 测试密码哈希
func TestPasswordManager_HashPassword(t *testing.T) {
	pm := NewPasswordManager()
	password := "TestPassword123!"

	hashed, err := pm.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hashed == "" {
		t.Error("HashPassword() returned empty string")
	}

	// 验证哈希格式：应该包含冒号分隔符
	if !contains(hashed, ":") {
		t.Error("HashPassword() should return format 'salt:hash'")
	}
}

// TestPasswordManager_VerifyPassword 测试密码验证
func TestPasswordManager_VerifyPassword(t *testing.T) {
	pm := NewPasswordManager()
	password := "TestPassword123!"

	hashed, err := pm.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	// 验证正确密码
	if !pm.VerifyPassword(password, hashed) {
		t.Error("VerifyPassword() should return true for correct password")
	}

	// 验证错误密码
	if pm.VerifyPassword("WrongPassword", hashed) {
		t.Error("VerifyPassword() should return false for wrong password")
	}
}

// TestPasswordManager_GeneratePassword 测试密码生成
func TestPasswordManager_GeneratePassword(t *testing.T) {
	pm := NewPasswordManager()

	tests := []struct {
		name     string
		length   int
		expected int
	}{
		{
			name:     "默认长度",
			length:   0,
			expected: model.DefaultPasswordLen,
		},
		{
			name:     "指定长度",
			length:   16,
			expected: 16,
		},
		{
			name:     "最大长度",
			length:   model.MaxPasswordLen,
			expected: model.MaxPasswordLen,
		},
		{
			name:     "超过最大长度",
			length:   model.MaxPasswordLen + 10,
			expected: model.MaxPasswordLen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := pm.GeneratePassword(tt.length)
			if err != nil {
				t.Errorf("GeneratePassword() error = %v", err)
			}

			if len(password) != tt.expected {
				t.Errorf("GeneratePassword() length = %d, want %d", len(password), tt.expected)
			}
		})
	}
}

// TestPasswordManager_ValidatePasswordStrength 测试密码强度验证
func TestPasswordManager_ValidatePasswordStrength(t *testing.T) {
	pm := NewPasswordManager()

	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{
			name:     "有效密码",
			password: "TestPassword123!",
			valid:    true,
		},
		{
			name:     "太短密码",
			password: "Test1!",
			valid:    false,
		},
		{
			name:     "太长密码",
			password: string(make([]byte, model.MaxPasswordLen+1)),
			valid:    false,
		},
		{
			name:     "缺少大写字母",
			password: "testpassword123!",
			valid:    false,
		},
		{
			name:     "缺少小写字母",
			password: "TESTPASSWORD123!",
			valid:    false,
		},
		{
			name:     "缺少数字",
			password: "TestPassword!",
			valid:    false,
		},
		{
			name:     "缺少特殊字符",
			password: "TestPassword123",
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, message := pm.ValidatePasswordStrength(tt.password)
			if valid != tt.valid {
				t.Errorf("ValidatePasswordStrength() = %v, want %v, message: %s", valid, tt.valid, message)
			}
		})
	}
}

// TestPasswordManager_VerifyPasswordWithDifferentSalts 测试不同盐的密码验证
func TestPasswordManager_VerifyPasswordWithDifferentSalts(t *testing.T) {
	pm := NewPasswordManager()
	password := "TestPassword123!"

	// 生成两个不同的哈希
	hashed1, err1 := pm.HashPassword(password)
	if err1 != nil {
		t.Fatalf("HashPassword() error = %v", err1)
	}

	hashed2, err2 := pm.HashPassword(password)
	if err2 != nil {
		t.Fatalf("HashPassword() error = %v", err2)
	}

	// 两个哈希应该不同（因为盐不同）
	if hashed1 == hashed2 {
		t.Error("HashPassword() should generate different hashes for same password")
	}

	// 但都应该能验证通过
	if !pm.VerifyPassword(password, hashed1) {
		t.Error("VerifyPassword() should work with first hash")
	}

	if !pm.VerifyPassword(password, hashed2) {
		t.Error("VerifyPassword() should work with second hash")
	}
}

// TestPasswordManager_InvalidHashFormat 测试无效哈希格式
func TestPasswordManager_InvalidHashFormat(t *testing.T) {
	pm := NewPasswordManager()
	password := "TestPassword123!"

	tests := []struct {
		name           string
		hashedPassword string
	}{
		{
			name:           "空字符串",
			hashedPassword: "",
		},
		{
			name:           "无冒号分隔符",
			hashedPassword: "invalidhash",
		},
		{
			name:           "多个冒号",
			hashedPassword: "salt:hash:extra",
		},
		{
			name:           "无效base64",
			hashedPassword: "invalidbase64:hash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if pm.VerifyPassword(password, tt.hashedPassword) {
				t.Errorf("VerifyPassword() should return false for invalid hash format: %s", tt.hashedPassword)
			}
		})
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && contains(s[1:], substr)
}
