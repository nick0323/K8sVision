package monitor

import (
	"go.uber.org/zap"
)

// Tracer 追踪器接口
type Tracer interface {
	// 可以在这里添加追踪相关的方法
}

// NewTracer 创建新的追踪器
func NewTracer(logger *zap.Logger) Tracer {
	logger.Info("初始化追踪器")
	// 这里可以添加实际的追踪器初始化逻辑
	return &tracerImpl{}
}

type tracerImpl struct{}

// InitTracing 初始化追踪
func InitTracing(logger *zap.Logger) {
	NewTracer(logger)
}
