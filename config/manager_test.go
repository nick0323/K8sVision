package config

import (
	"os"
	"testing"
	"time"

	"github.com/nick0323/K8sVision/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewManager(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)
	
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.GetConfig())
}

func TestDefaultConfig(t *testing.T) {
	config := model.DefaultConfig()
	
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, "info", config.Log.Level)
	assert.Equal(t, "k8svision-secret-key", config.JWT.Secret)
	assert.Equal(t, 5, config.Auth.MaxLoginFail)
}

func TestLoadConfig(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)
	
	// 测试加载默认配置
	err := manager.Load("")
	assert.NoError(t, err)
	
	config := manager.GetConfig()
	assert.Equal(t, "8080", config.Server.Port)
}

func TestEnvironmentOverrides(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)
	
	// 设置环境变量
	os.Setenv("K8SVISION_SERVER_PORT", "9090")
	os.Setenv("K8SVISION_LOG_LEVEL", "debug")
	os.Setenv("JWT_SECRET", "test-secret")
	defer func() {
		os.Unsetenv("K8SVISION_SERVER_PORT")
		os.Unsetenv("K8SVISION_LOG_LEVEL")
		os.Unsetenv("JWT_SECRET")
	}()
	
	err := manager.Load("")
	assert.NoError(t, err)
	
	config := manager.GetConfig()
	assert.Equal(t, "9090", config.Server.Port)
	assert.Equal(t, "debug", config.Log.Level)
	assert.Equal(t, "test-secret", config.JWT.Secret)
}

func TestConfigValidation(t *testing.T) {
	config := model.DefaultConfig()
	err := config.Validate()
	assert.NoError(t, err)
}

func TestGetServerAddress(t *testing.T) {
	config := model.DefaultConfig()
	address := config.GetServerAddress()
	assert.Equal(t, "0.0.0.0:8080", address)
}

func TestIsDevelopment(t *testing.T) {
	config := model.DefaultConfig()
	assert.False(t, config.IsDevelopment())
	
	config.Log.Level = "debug"
	assert.True(t, config.IsDevelopment())
}

func TestIsProduction(t *testing.T) {
	config := model.DefaultConfig()
	assert.True(t, config.IsProduction())
	
	config.Log.Level = "debug"
	assert.False(t, config.IsProduction())
}

func TestConfigGetters(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)
	
	err := manager.Load("")
	assert.NoError(t, err)
	
	// 测试配置结构体直接访问
	config := manager.GetConfig()
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, float32(100), config.Kubernetes.QPS)
	assert.Equal(t, true, config.Kubernetes.Insecure)
	assert.Equal(t, 30*time.Second, config.Kubernetes.Timeout)
}

func TestSetConfig(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)
	
	err := manager.Load("")
	assert.NoError(t, err)
	
	// 设置配置值
	manager.Set("server.port", "9090")
	assert.Equal(t, "9090", manager.GetString("server.port"))
}

func TestConfigReload(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)
	
	// 创建临时配置文件
	tempFile := "test_config.yaml"
	content := `
server:
  port: 8080
log:
  level: info
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