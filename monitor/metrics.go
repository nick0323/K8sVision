package monitor

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

// Metrics 性能指标
type Metrics struct {
	// 请求统计
	TotalRequests    int64         `json:"totalRequests"`
	SuccessfulRequests int64       `json:"successfulRequests"`
	FailedRequests   int64         `json:"failedRequests"`
	
	// 响应时间统计
	MinResponseTime  time.Duration `json:"minResponseTime"`
	MaxResponseTime  time.Duration `json:"maxResponseTime"`
	AvgResponseTime  time.Duration `json:"avgResponseTime"`
	
	// 缓存统计
	CacheHits        int64         `json:"cacheHits"`
	CacheMisses      int64         `json:"cacheMisses"`
	
	// 并发统计
	CurrentConnections int64       `json:"currentConnections"`
	MaxConnections    int64        `json:"maxConnections"`
	
	// 错误统计
	ErrorCount       int64         `json:"errorCount"`
	LastError        string        `json:"lastError"`
	LastErrorTime    time.Time     `json:"lastErrorTime"`
	
	// 系统统计
	StartTime        time.Time     `json:"startTime"`
	Uptime           time.Duration `json:"uptime"`
	
	// 内部字段
	mutex            sync.RWMutex
	totalResponseTime time.Duration
	requestCount     int64
	logger           *zap.Logger
}

// NewMetrics 创建新的性能指标收集器
func NewMetrics(logger *zap.Logger) *Metrics {
	return &Metrics{
		StartTime: time.Now(),
		logger:    logger,
	}
}

// RecordRequest 记录请求
func (m *Metrics) RecordRequest(success bool, responseTime time.Duration) {
	atomic.AddInt64(&m.TotalRequests, 1)
	
	if success {
		atomic.AddInt64(&m.SuccessfulRequests, 1)
	} else {
		atomic.AddInt64(&m.FailedRequests, 1)
	}
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	// 更新响应时间统计
	m.totalResponseTime += responseTime
	m.requestCount++
	
	// 更新最小响应时间
	if m.MinResponseTime == 0 || responseTime < m.MinResponseTime {
		m.MinResponseTime = responseTime
	}
	
	// 更新最大响应时间
	if responseTime > m.MaxResponseTime {
		m.MaxResponseTime = responseTime
	}
	
	// 计算平均响应时间
	if m.requestCount > 0 {
		m.AvgResponseTime = m.totalResponseTime / time.Duration(m.requestCount)
	}
}

// RecordCacheHit 记录缓存命中
func (m *Metrics) RecordCacheHit() {
	atomic.AddInt64(&m.CacheHits, 1)
}

// RecordCacheMiss 记录缓存未命中
func (m *Metrics) RecordCacheMiss() {
	atomic.AddInt64(&m.CacheMisses, 1)
}

// RecordConnection 记录连接
func (m *Metrics) RecordConnection() {
	current := atomic.AddInt64(&m.CurrentConnections, 1)
	
	// 更新最大连接数
	for {
		max := atomic.LoadInt64(&m.MaxConnections)
		if current <= max || atomic.CompareAndSwapInt64(&m.MaxConnections, max, current) {
			break
		}
	}
}

// RecordDisconnection 记录断开连接
func (m *Metrics) RecordDisconnection() {
	atomic.AddInt64(&m.CurrentConnections, -1)
}

// RecordError 记录错误
func (m *Metrics) RecordError(err string) {
	atomic.AddInt64(&m.ErrorCount, 1)
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.LastError = err
	m.LastErrorTime = time.Now()
}

// GetStats 获取统计信息
func (m *Metrics) GetStats() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	uptime := time.Since(m.StartTime)
	totalRequests := atomic.LoadInt64(&m.TotalRequests)
	successfulRequests := atomic.LoadInt64(&m.SuccessfulRequests)
	failedRequests := atomic.LoadInt64(&m.FailedRequests)
	cacheHits := atomic.LoadInt64(&m.CacheHits)
	cacheMisses := atomic.LoadInt64(&m.CacheMisses)
	currentConnections := atomic.LoadInt64(&m.CurrentConnections)
	maxConnections := atomic.LoadInt64(&m.MaxConnections)
	errorCount := atomic.LoadInt64(&m.ErrorCount)
	
	// 计算成功率
	successRate := 0.0
	if totalRequests > 0 {
		successRate = float64(successfulRequests) / float64(totalRequests) * 100
	}
	
	// 计算缓存命中率
	cacheHitRate := 0.0
	totalCacheRequests := cacheHits + cacheMisses
	if totalCacheRequests > 0 {
		cacheHitRate = float64(cacheHits) / float64(totalCacheRequests) * 100
	}
	
	return map[string]interface{}{
		"totalRequests":      totalRequests,
		"successfulRequests": successfulRequests,
		"failedRequests":     failedRequests,
		"successRate":        successRate,
		"minResponseTime":    m.MinResponseTime.String(),
		"maxResponseTime":    m.MaxResponseTime.String(),
		"avgResponseTime":    m.AvgResponseTime.String(),
		"cacheHits":          cacheHits,
		"cacheMisses":        cacheMisses,
		"cacheHitRate":       cacheHitRate,
		"currentConnections": currentConnections,
		"maxConnections":     maxConnections,
		"errorCount":         errorCount,
		"lastError":          m.LastError,
		"lastErrorTime":      m.LastErrorTime,
		"startTime":          m.StartTime,
		"uptime":             uptime.String(),
	}
}

// Reset 重置统计信息
func (m *Metrics) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	atomic.StoreInt64(&m.TotalRequests, 0)
	atomic.StoreInt64(&m.SuccessfulRequests, 0)
	atomic.StoreInt64(&m.FailedRequests, 0)
	atomic.StoreInt64(&m.CacheHits, 0)
	atomic.StoreInt64(&m.CacheMisses, 0)
	atomic.StoreInt64(&m.CurrentConnections, 0)
	atomic.StoreInt64(&m.MaxConnections, 0)
	atomic.StoreInt64(&m.ErrorCount, 0)
	
	m.MinResponseTime = 0
	m.MaxResponseTime = 0
	m.AvgResponseTime = 0
	m.totalResponseTime = 0
	m.requestCount = 0
	m.LastError = ""
	m.LastErrorTime = time.Time{}
	m.StartTime = time.Now()
	
	m.logger.Info("性能指标已重置")
}

// Monitor 性能监控器
type Monitor struct {
	metrics *Metrics
	logger  *zap.Logger
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewMonitor 创建新的性能监控器
func NewMonitor(logger *zap.Logger) *Monitor {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Monitor{
		metrics: NewMetrics(logger),
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// GetMetrics 获取性能指标
func (m *Monitor) GetMetrics() *Metrics {
	return m.metrics
}

// StartPeriodicLogging 启动定期日志记录
func (m *Monitor) StartPeriodicLogging(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				stats := m.metrics.GetStats()
				m.logger.Info("性能指标统计",
					zap.Int64("totalRequests", stats["totalRequests"].(int64)),
					zap.Float64("successRate", stats["successRate"].(float64)),
					zap.String("avgResponseTime", stats["avgResponseTime"].(string)),
					zap.Float64("cacheHitRate", stats["cacheHitRate"].(float64)),
					zap.Int64("currentConnections", stats["currentConnections"].(int64)),
				)
			case <-m.ctx.Done():
				return
			}
		}
	}()
}

// Close 关闭监控器
func (m *Monitor) Close() {
	m.cancel()
	m.logger.Info("性能监控器已关闭")
} 