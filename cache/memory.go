package cache

import (
	"context"
	"sync"
	"time"

	"github.com/nick0323/K8sVision/model"
	"go.uber.org/zap"
)

type CacheItem struct {
	Value      interface{}
	ExpireTime time.Time
	CreatedAt  time.Time
}

type MemoryCache struct {
	data            map[string]CacheItem
	mutex           sync.RWMutex
	maxSize         int
	ttl             time.Duration
	cleanupInterval time.Duration
	logger          *zap.Logger
	ctx             context.Context
	cancel          context.CancelFunc
}

// NewMemoryCache 创建新的内存缓存
func NewMemoryCache(config *model.CacheConfig, logger *zap.Logger) *MemoryCache {
	ctx, cancel := context.WithCancel(context.Background())

	cache := &MemoryCache{
		data:            make(map[string]CacheItem),
		maxSize:         config.MaxSize,
		ttl:             config.TTL,
		cleanupInterval: config.CleanupInterval,
		logger:          logger,
		ctx:             ctx,
		cancel:          cancel,
	}

	// 启动清理协程
	if config.Enabled {
		go cache.cleanupWorker()
	}

	return cache
}

// Set 设置缓存
func (c *MemoryCache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.ttl)
}

// SetWithTTL 设置缓存并指定TTL
func (c *MemoryCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 检查容量限制
	if len(c.data) >= c.maxSize {
		c.evictOldest()
	}

	now := time.Now()
	c.data[key] = CacheItem{
		Value:      value,
		ExpireTime: now.Add(ttl),
		CreatedAt:  now,
	}

	// 安全检查logger
	if c.logger != nil {
		c.logger.Debug("缓存设置", zap.String("key", key), zap.Duration("ttl", ttl))
	}
}

// Get 获取缓存
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// 检查是否过期
	if time.Now().After(item.ExpireTime) {
		// 异步删除过期项
		go c.delete(key)
		return nil, false
	}

	// 安全检查logger
	if c.logger != nil {
		c.logger.Debug("缓存命中", zap.String("key", key))
	}
	return item.Value, true
}

// Delete 删除缓存
func (c *MemoryCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)

	// 安全检查logger
	if c.logger != nil {
		c.logger.Debug("缓存删除", zap.String("key", key))
	}
}

// delete 内部删除方法（无锁）
func (c *MemoryCache) delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

// Clear 清空缓存
func (c *MemoryCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]CacheItem)

	// 安全检查logger
	if c.logger != nil {
		c.logger.Info("缓存已清空")
	}
}

// Size 获取缓存大小
func (c *MemoryCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.data)
}

// Keys 获取所有键
func (c *MemoryCache) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}

// evictOldest 淘汰最旧的项
func (c *MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, item := range c.data {
		if first || item.CreatedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.CreatedAt
			first = false
		}
	}

	if oldestKey != "" {
		delete(c.data, oldestKey)

		// 安全检查logger
		if c.logger != nil {
			c.logger.Debug("缓存淘汰", zap.String("key", oldestKey))
		}
	}
}

// cleanupWorker 清理工作协程
func (c *MemoryCache) cleanupWorker() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.ctx.Done():
			return
		}
	}
}

// cleanup 清理过期项
func (c *MemoryCache) cleanup() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	expiredCount := 0

	for key, item := range c.data {
		if now.After(item.ExpireTime) {
			delete(c.data, key)
			expiredCount++
		}
	}

	if expiredCount > 0 {
		// 安全检查logger
		if c.logger != nil {
			c.logger.Info("缓存清理完成", zap.Int("expiredCount", expiredCount))
		}
	}

	return expiredCount
}

// ClearExpired 手动清理过期项并返回清理数量
func (c *MemoryCache) ClearExpired() int {
	return c.cleanup()
}

// Close 关闭缓存
func (c *MemoryCache) Close() {
	c.cancel()
	c.Clear()

	// 安全检查logger
	if c.logger != nil {
		c.logger.Info("缓存已关闭")
	}
}

// GetStats 获取缓存统计信息
func (c *MemoryCache) GetStats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	now := time.Now()
	expiredCount := 0
	totalSize := 0
	hitRate := 0.0

	// 计算过期项数量
	for _, item := range c.data {
		if now.After(item.ExpireTime) {
			expiredCount++
		}
		// 简单估算内存使用
		totalSize += 100 // 假设每个项平均100字节
	}

	// 计算命中率（这里需要添加命中/未命中统计）
	// 在实际使用中，应该维护这些统计信息
	totalRequests := 0 // 需要从其他地方获取
	if totalRequests > 0 {
		hitRate = float64(len(c.data)) / float64(totalRequests) * 100
	}

	return map[string]interface{}{
		"size":         len(c.data),
		"maxSize":      c.maxSize,
		"expiredCount": expiredCount,
		"totalSize":    totalSize,
		"ttl":          c.ttl.String(),
		"hitRate":      hitRate,
		"utilization":  float64(len(c.data)) / float64(c.maxSize) * 100,
	}
}
