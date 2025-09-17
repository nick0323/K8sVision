package api

import (
	"encoding/base64"
	"strings"
)

// isHashedPassword 判断密码是否为哈希格式
// 哈希格式: base64(salt):bcrypt_hash
func isHashedPassword(password string) bool {
	// 检查是否包含冒号分隔符
	if !strings.Contains(password, ":") {
		return false
	}

	parts := strings.Split(password, ":")
	if len(parts) != 2 {
		return false
	}

	// 检查第一部分是否为有效的base64编码
	_, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	// 检查第二部分是否为有效的bcrypt哈希（以$2a$、$2b$或$2y$开头）
	hashPart := parts[1]
	if len(hashPart) < 60 { // bcrypt哈希长度至少60字符
		return false
	}

	return strings.HasPrefix(hashPart, "$2a$") ||
		strings.HasPrefix(hashPart, "$2b$") ||
		strings.HasPrefix(hashPart, "$2y$")
}
