package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"
	"github.com/nick0323/K8sVision/monitor"
	"go.uber.org/zap"
)

// MetricsResponse 指标响应
type MetricsResponse struct {
	Timestamp time.Time              `json:"timestamp"`
	Metrics   map[string]interface{} `json:"metrics"`
	Summary   map[string]interface{} `json:"summary"`
}

// RegisterMetrics 注册指标相关路由
func RegisterMetrics(r *gin.RouterGroup, logger *zap.Logger) {
	r.GET("/metrics", getMetrics(logger))
	r.GET("/metrics/business", getBusinessMetrics(logger))
	r.GET("/metrics/system", getSystemMetrics(logger))
	r.GET("/metrics/health", getHealthMetrics(logger))
}

// getMetrics 获取所有指标
// @Summary 获取系统指标
// @Description 获取系统性能指标和业务指标
// @Tags Metrics
// @Security BearerAuth
// @Produce json
// @Success 200 {object} MetricsResponse
// @Router /metrics [get]
func getMetrics(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取系统指标
		systemMetrics := monitor.GetMetricsManager().GetAllMetrics()

		// 获取业务指标
		businessCollector := monitor.GetBusinessMetricsCollector()
		businessMetrics := businessCollector.CollectMetrics()

		// 构建响应
		response := MetricsResponse{
			Timestamp: time.Now(),
			Metrics: map[string]interface{}{
				"system":   systemMetrics,
				"business": businessMetrics,
			},
			Summary: map[string]interface{}{
				"total_metrics":  len(systemMetrics) + len(businessMetrics),
				"system_count":   len(systemMetrics),
				"business_count": len(businessMetrics),
			},
		}

		middleware.ResponseSuccess(c, response, "指标获取成功", nil)
	}
}

// getBusinessMetrics 获取业务指标
// @Summary 获取业务指标
// @Description 获取K8s资源和API相关的业务指标
// @Tags Metrics
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.APIResponse
// @Router /metrics/business [get]
func getBusinessMetrics(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		collector := monitor.GetBusinessMetricsCollector()
		metrics := collector.CollectMetrics()

		// 按类型分组指标
		groupedMetrics := make(map[string][]monitor.BusinessMetric)
		for _, metric := range metrics {
			groupedMetrics[metric.MetricType] = append(groupedMetrics[metric.MetricType], metric)
		}

		middleware.ResponseSuccess(c, groupedMetrics, "业务指标获取成功", nil)
	}
}

// getSystemMetrics 获取系统指标
// @Summary 获取系统指标
// @Description 获取系统性能指标，包括CPU、内存、网络等
// @Tags Metrics
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.APIResponse
// @Router /metrics/system [get]
func getSystemMetrics(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricsManager := monitor.GetMetricsManager()
		systemMetrics := metricsManager.GetAllMetrics()

		// 计算系统健康状态
		healthStatus := calculateSystemHealth(systemMetrics)

		response := map[string]interface{}{
			"metrics": systemMetrics,
			"health":  healthStatus,
		}

		middleware.ResponseSuccess(c, response, "系统指标获取成功", nil)
	}
}

// getHealthMetrics 获取健康指标
// @Summary 获取健康指标
// @Description 获取系统健康状态和关键指标
// @Tags Metrics
// @Produce json
// @Success 200 {object} model.APIResponse
// @Router /metrics/health [get]
func getHealthMetrics(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricsManager := monitor.GetMetricsManager()
		systemMetrics := metricsManager.GetAllMetrics()

		// 计算健康状态
		healthStatus := calculateSystemHealth(systemMetrics)

		// 获取关键指标
		keyMetrics := extractKeyMetrics(systemMetrics)

		response := map[string]interface{}{
			"status":      healthStatus["status"],
			"score":       healthStatus["score"],
			"key_metrics": keyMetrics,
			"timestamp":   time.Now(),
		}

		// 根据健康状态设置HTTP状态码
		status := healthStatus["status"].(string)
		httpStatus := http.StatusOK
		if status == "unhealthy" {
			httpStatus = http.StatusServiceUnavailable
		} else if status == "degraded" {
			httpStatus = http.StatusOK
		}

		c.JSON(httpStatus, gin.H{
			"code":    model.CodeSuccess,
			"message": "健康指标获取成功",
			"data":    response,
		})
	}
}

// calculateSystemHealth 计算系统健康状态
func calculateSystemHealth(metrics map[string]interface{}) map[string]interface{} {
	// 这里应该根据实际指标计算健康状态
	// 目前返回模拟数据
	healthScore := 95.0
	status := "healthy"

	// 检查关键指标
	if healthScore < 50 {
		status = "unhealthy"
	} else if healthScore < 80 {
		status = "degraded"
	}

	return map[string]interface{}{
		"status": status,
		"score":  healthScore,
		"checks": map[string]interface{}{
			"cpu":     "ok",
			"memory":  "ok",
			"disk":    "ok",
			"network": "ok",
		},
	}
}

// extractKeyMetrics 提取关键指标
func extractKeyMetrics(metrics map[string]interface{}) map[string]interface{} {
	// 这里应该从实际指标中提取关键信息
	// 目前返回模拟数据
	return map[string]interface{}{
		"cpu_usage":    25.5,
		"memory_usage": 1024 * 1024 * 100, // 100MB
		"disk_usage":   45.2,
		"network_io": map[string]float64{
			"in":  1024 * 1024, // 1MB/s
			"out": 512 * 1024,  // 512KB/s
		},
		"active_connections": 150,
		"request_rate":       100.5, // requests/second
	}
}

