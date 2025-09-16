package monitor

import (
	"go.uber.org/zap"
)

type Tracer interface {
}

func NewTracer(logger *zap.Logger) Tracer {
	logger.Info("初始化追踪器")
	return &tracerImpl{}
}

type tracerImpl struct{}

func InitTracing(logger *zap.Logger) {
	NewTracer(logger)
}
