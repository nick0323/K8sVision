package config

import (
	"os"
	"testing"
	"time"

	"github.com/nick0323/K8sVision/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestConfigValidation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	_ = NewManager(logger)

	t.Run("有效配置", func(t *testing.T) {
		config := model.DefaultConfig()
		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("无效服务器端口", func(t *testing.T) {
		config := model.DefaultConfig()
		config.Server.Port = ""
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "服务器端口不能为空")
	})

	t.Run("无效JWT密钥长度", func(t *testing.T) {
		config := model.DefaultConfig()
		config.JWT.Secret = "short"
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "JWT密钥长度至少32位")
	})

	t.Run("无效日志级别", func(t *testing.T) {
		config := model.DefaultConfig()
		config.Log.Level = "invalid"
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的日志级别")
	})

	t.Run("无效日志格式", func(t *testing.T) {
		config := model.DefaultConfig()
		config.Log.Format = "invalid"
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的日志格式")
	})

	t.Run("无效认证配置", func(t *testing.T) {
		config := model.DefaultConfig()
		config.Auth.Username = ""
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "认证用户名不能为空")
	})

	t.Run("无效缓存配置", func(t *testing.T) {
		config := model.DefaultConfig()
		config.Cache.Enabled = true
		config.Cache.TTL = 0
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "缓存TTL必须大于0")
	})

	t.Run("无效Kubernetes配置", func(t *testing.T) {
		config := model.DefaultConfig()
		config.Kubernetes.QPS = 0
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Kubernetes QPS必须大于0")
	})
}

func TestConfigGetters(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)

	err := manager.Load("")
	assert.NoError(t, err)

	t.Run("GetJWTSecret", func(t *testing.T) {
		secret := manager.GetJWTSecret()
		assert.NotEmpty(t, secret)
		assert.Equal(t, "k8svision-secret-key", string(secret))
	})

	t.Run("GetAuthConfig", func(t *testing.T) {
		authConfig := manager.GetAuthConfig()
		assert.NotNil(t, authConfig)
		assert.Equal(t, "admin", authConfig.Username)
		assert.Equal(t, "admin", authConfig.Password)
		assert.Equal(t, 5, authConfig.MaxLoginFail)
	})

	t.Run("GetString", func(t *testing.T) {
		port := manager.GetString("server.port")
		assert.Equal(t, "8080", port)
	})

	t.Run("GetInt", func(t *testing.T) {
		qps := manager.GetInt("kubernetes.qps")
		assert.Equal(t, 100, qps)
	})

	t.Run("GetBool", func(t *testing.T) {
		insecure := manager.GetBool("kubernetes.insecure")
		assert.True(t, insecure)
	})

	t.Run("GetDuration", func(t *testing.T) {
		timeout := manager.GetDuration("kubernetes.timeout")
		assert.Equal(t, 30*time.Second, timeout)
	})
}

func TestConfigSetters(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)

	err := manager.Load("")
	assert.NoError(t, err)

	t.Run("Set配置值", func(t *testing.T) {
		manager.Set("server.port", "9090")
		port := manager.GetString("server.port")
		assert.Equal(t, "9090", port)
	})

	t.Run("Set后重新加载", func(t *testing.T) {
		manager.Set("log.level", "debug")
		level := manager.GetString("log.level")
		assert.Equal(t, "debug", level)
	})
}

func TestEnvironmentOverrides(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)

	// 设置环境变量
	os.Setenv("K8SVISION_SERVER_PORT", "9090")
	os.Setenv("K8SVISION_LOG_LEVEL", "debug")
	os.Setenv("JWT_SECRET", "test-secret-key-32-chars-long")
	os.Setenv("LOGIN_USERNAME", "testuser")
	os.Setenv("LOGIN_PASSWORD", "testpass")
	defer func() {
		os.Unsetenv("K8SVISION_SERVER_PORT")
		os.Unsetenv("K8SVISION_LOG_LEVEL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("LOGIN_USERNAME")
		os.Unsetenv("LOGIN_PASSWORD")
	}()

	err := manager.Load("")
	assert.NoError(t, err)

	config := manager.GetConfig()

	// 验证环境变量覆盖
	assert.Equal(t, "9090", config.Server.Port)
	assert.Equal(t, "debug", config.Log.Level)
	assert.Equal(t, "test-secret-key-32-chars-long", config.JWT.Secret)
	assert.Equal(t, "testuser", config.Auth.Username)
	assert.Equal(t, "testpass", config.Auth.Password)
}

func TestConfigReload(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)

	// 创建临时配置文件
	tempFile := "test_config_reload.yaml"
	content := `
server:
  port: 8080
log:
  level: info
jwt:
  secret: "test-secret-key-32-chars-long"
`
	err := os.WriteFile(tempFile, []byte(content), 0644)
	assert.NoError(t, err)
	defer os.Remove(tempFile)

	// 加载配置
	err = manager.Load(tempFile)
	assert.NoError(t, err)

	config := manager.GetConfig()
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "info", config.Log.Level)

	// 修改配置文件
	content = `
server:
  port: 9090
log:
  level: debug
jwt:
  secret: "test-secret-key-32-chars-long"
`
	err = os.WriteFile(tempFile, []byte(content), 0644)
	assert.NoError(t, err)

	// 手动重新加载
	err = manager.reload()
	assert.NoError(t, err)

	config = manager.GetConfig()
	assert.Equal(t, "9090", config.Server.Port)
	assert.Equal(t, "debug", config.Log.Level)
}
