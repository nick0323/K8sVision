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
	healthScore := 100.0
	checks := make(map[string]interface{})

	// 检查系统指标是否可用
	if len(metrics) == 0 {
		healthScore -= 30
		checks["metrics"] = "unavailable"
	} else {
		checks["metrics"] = "ok"
	}

	// 检查监控系统
	if monitor.GetMetricsManager() != nil {
		checks["monitoring"] = "ok"
	} else {
		healthScore -= 20
		checks["monitoring"] = "degraded"
	}

	// 检查业务指标收集器
	if monitor.GetBusinessMetricsCollector() != nil {
		checks["business_metrics"] = "ok"
	} else {
		healthScore -= 15
		checks["business_metrics"] = "degraded"
	}

	// 确定整体状态
	var status string
	if healthScore >= 90 {
		status = "healthy"
	} else if healthScore >= 70 {
		status = "degraded"
	} else {
		status = "unhealthy"
	}

	return map[string]interface{}{
		"status": status,
		"score":  healthScore,
		"checks": checks,
	}
}

// extractKeyMetrics 提取关键指标
func extractKeyMetrics(metrics map[string]interface{}) map[string]interface{} {
	keyMetrics := make(map[string]interface{})

	// 从实际指标中提取关键信息
	if len(metrics) > 0 {
		// 尝试提取CPU相关指标
		if cpu, exists := metrics["cpu"]; exists {
			keyMetrics["cpu_usage"] = cpu
		} else {
			keyMetrics["cpu_usage"] = "N/A"
		}

		// 尝试提取内存相关指标
		if memory, exists := metrics["memory"]; exists {
			keyMetrics["memory_usage"] = memory
		} else {
			keyMetrics["memory_usage"] = "N/A"
		}

		// 尝试提取网络相关指标
		if network, exists := metrics["network"]; exists {
			keyMetrics["network_io"] = network
		} else {
			keyMetrics["network_io"] = map[string]string{
				"in":  "N/A",
				"out": "N/A",
			}
		}

		// 尝试提取连接数指标
		if connections, exists := metrics["connections"]; exists {
			keyMetrics["active_connections"] = connections
		} else {
			keyMetrics["active_connections"] = "N/A"
		}

		// 添加指标收集时间
		keyMetrics["collected_at"] = time.Now().Unix()
		keyMetrics["metrics_count"] = len(metrics)
	} else {
		// 没有指标数据时的默认值
		keyMetrics = map[string]interface{}{
			"cpu_usage":          "N/A",
			"memory_usage":       "N/A",
			"network_io":         map[string]string{"in": "N/A", "out": "N/A"},
			"active_connections": "N/A",
			"collected_at":       time.Now().Unix(),
			"metrics_count":      0,
			"status":             "no_data",
		}
	}

	return keyMetrics
}
