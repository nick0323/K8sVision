package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsRecorder 性能指标记录接口
type MetricsRecorder interface {
	RecordRequest(success bool, responseTime time.Duration)
	RecordConnection()
	RecordDisconnection()
	RecordError(err string)
	RecordCacheHit()
	RecordCacheMiss()
}

// MetricsMiddleware 性能监控中间件
func MetricsMiddleware(recorder MetricsRecorder) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 记录连接
		recorder.RecordConnection()
		defer recorder.RecordDisconnection()
		
		// 处理请求
		c.Next()
		
		// 计算响应时间
		responseTime := time.Since(start)
		
		// 判断请求是否成功
		success := c.Writer.Status() < 400
		
		// 记录请求统计
		recorder.RecordRequest(success, responseTime)
		
		// 记录错误
		if !success {
			recorder.RecordError(c.Errors.String())
		}
	}
}

// CacheMiddleware 缓存中间件
func CacheMiddleware(cacheManager interface{}, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只对GET请求进行缓存
		if c.Request.Method != "GET" {
			c.Next()
			return
		}
		
		// 生成缓存键
		cacheKey := generateCacheKey(c)
		
		// 尝试从缓存获取
		if cache, ok := cacheManager.(interface{ Get(key string) (interface{}, bool) }); ok {
			if value, exists := cache.Get(cacheKey); exists {
				// 缓存命中
				if recorder, ok := c.Get("metrics"); ok {
					if r, ok := recorder.(MetricsRecorder); ok {
						r.RecordCacheHit()
					}
				}
				
				// 返回缓存数据
				c.JSON(200, value)
				c.Abort()
				return
			}
		}
		
		// 缓存未命中
		if recorder, ok := c.Get("metrics"); ok {
			if r, ok := recorder.(MetricsRecorder); ok {
				r.RecordCacheMiss()
			}
		}
		
		// 继续处理请求
		c.Next()
		
		// 如果请求成功，缓存结果
		if c.Writer.Status() == 200 {
			// 这里需要获取响应数据，但gin的响应已经写入，无法直接获取
			// 实际使用时需要在handler中手动设置缓存
		}
	}
}

// generateCacheKey 生成缓存键
func generateCacheKey(c *gin.Context) string {
	// 使用请求路径和查询参数生成缓存键
	key := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		key += "?" + c.Request.URL.RawQuery
	}
	return key
}

// ConcurrencyMiddleware 并发控制中间件
func ConcurrencyMiddleware(maxConcurrency int) gin.HandlerFunc {
	semaphore := make(chan struct{}, maxConcurrency)
	
	return func(c *gin.Context) {
		select {
		case semaphore <- struct{}{}:
			defer func() { <-semaphore }()
			c.Next()
		default:
			// 并发数超限，返回503错误
			c.JSON(503, gin.H{
				"code":    503,
				"message": "服务繁忙，请稍后重试",
			})
			c.Abort()
		}
	}
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	// 简单的内存速率限制器
	// 实际生产环境建议使用Redis等外部存储
	limiters := make(map[string]*rateLimiter)
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		limiter, exists := limiters[clientIP]
		if !exists {
			limiter = newRateLimiter(limit, window)
			limiters[clientIP] = limiter
		}
		
		if !limiter.Allow() {
			c.JSON(429, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后重试",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// rateLimiter 简单的速率限制器
type rateLimiter struct {
	limit   int
	window  time.Duration
	tokens  chan struct{}
	lastRefill time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		limit:   limit,
		window:  window,
		tokens:  make(chan struct{}, limit),
		lastRefill: time.Now(),
	}
	
	// 初始化令牌
	for i := 0; i < limit; i++ {
		rl.tokens <- struct{}{}
	}
	
	// 启动令牌补充协程
	go rl.refillTokens()
	
	return rl
}

func (rl *rateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func (rl *rateLimiter) refillTokens() {
	ticker := time.NewTicker(rl.window / time.Duration(rl.limit))
	defer ticker.Stop()
	
	for range ticker.C {
		select {
		case rl.tokens <- struct{}{}:
			// 令牌补充成功
		default:
			// 令牌桶已满
		}
	}
} 