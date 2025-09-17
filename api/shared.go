package api

import "github.com/nick0323/K8sVision/config"

// 共享的配置管理器
var configManager *config.Manager

// SetConfigManager 设置配置管理器
func SetConfigManager(cm *config.Manager) {
	configManager = cm
}

// GetConfigManager 获取配置管理器
func GetConfigManager() *config.Manager {
	return configManager
}
