package config

import (
	"testing"

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

func TestSetConfig(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)

	err := manager.Load("")
	assert.NoError(t, err)

	// 设置配置值
	manager.Set("server.port", "9090")
	assert.Equal(t, "9090", manager.GetString("server.port"))
}
