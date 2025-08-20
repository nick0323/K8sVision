package cache

import (
	"testing"
	"time"

	"github.com/nick0323/K8sVision/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewMemoryCache(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &model.CacheConfig{
		Enabled:         true,
		Type:           "memory",
		TTL:            5 * time.Minute,
		MaxSize:        100,
		CleanupInterval: 1 * time.Minute,
	}
	
	cache := NewMemoryCache(config, logger)
	assert.NotNil(t, cache)
	assert.Equal(t, 100, cache.maxSize)
	assert.Equal(t, 5*time.Minute, cache.ttl)
}

func TestMemoryCacheSetAndGet(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &model.CacheConfig{
		Enabled:         true,
		Type:           "memory",
		TTL:            1 * time.Minute,
		MaxSize:        10,
		CleanupInterval: 1 * time.Minute,
	}
	
	cache := NewMemoryCache(config, logger)
	defer cache.Close()
	
	// 测试设置和获取
	cache.Set("key1", "value1")
	
	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)
	
	// 测试不存在的键
	value, exists = cache.Get("key2")
	assert.False(t, exists)
	assert.Nil(t, value)
}

func TestMemoryCacheTTL(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &model.CacheConfig{
		Enabled:         true,
		Type:           "memory",
		TTL:            100 * time.Millisecond,
		MaxSize:        10,
		CleanupInterval: 50 * time.Millisecond,
	}
	
	cache := NewMemoryCache(config, logger)
	defer cache.Close()
	
	// 设置短期缓存
	cache.Set("key1", "value1")
	
	// 立即获取应该存在
	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)
	
	// 等待过期
	time.Sleep(150 * time.Millisecond)
	
	// 过期后应该不存在
	value, exists = cache.Get("key1")
	assert.False(t, exists)
	assert.Nil(t, value)
}

func TestMemoryCacheMaxSize(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &model.CacheConfig{
		Enabled:         true,
		Type:           "memory",
		TTL:            1 * time.Minute,
		MaxSize:        2,
		CleanupInterval: 1 * time.Minute,
	}
	
	cache := NewMemoryCache(config, logger)
	defer cache.Close()
	
	// 添加两个项目
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	
	assert.Equal(t, 2, cache.Size())
	
	// 添加第三个项目，应该触发淘汰
	cache.Set("key3", "value3")
	
	// 应该只有2个项目
	assert.Equal(t, 2, cache.Size())
	
	// key1应该被淘汰（最旧的）
	_, exists := cache.Get("key1")
	assert.False(t, exists)
	
	// key2和key3应该存在
	_, exists = cache.Get("key2")
	assert.True(t, exists)
	_, exists = cache.Get("key3")
	assert.True(t, exists)
}

func TestMemoryCacheDelete(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &model.CacheConfig{
		Enabled:         true,
		Type:           "memory",
		TTL:            1 * time.Minute,
		MaxSize:        10,
		CleanupInterval: 1 * time.Minute,
	}
	
	cache := NewMemoryCache(config, logger)
	defer cache.Close()
	
	cache.Set("key1", "value1")
	
	// 删除前应该存在
	_, exists := cache.Get("key1")
	assert.True(t, exists)
	
	// 删除
	cache.Delete("key1")
	
	// 删除后应该不存在
	_, exists = cache.Get("key1")
	assert.False(t, exists)
}

func TestMemoryCacheClear(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &model.CacheConfig{
		Enabled:         true,
		Type:           "memory",
		TTL:            1 * time.Minute,
		MaxSize:        10,
		CleanupInterval: 1 * time.Minute,
	}
	
	cache := NewMemoryCache(config, logger)
	defer cache.Close()
	
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	
	assert.Equal(t, 2, cache.Size())
	
	cache.Clear()
	
	assert.Equal(t, 0, cache.Size())
	
	// 所有键都应该不存在
	_, exists := cache.Get("key1")
	assert.False(t, exists)
	_, exists = cache.Get("key2")
	assert.False(t, exists)
}

func TestMemoryCacheGetStats(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &model.CacheConfig{
		Enabled:         true,
		Type:           "memory",
		TTL:            1 * time.Minute,
		MaxSize:        10,
		CleanupInterval: 1 * time.Minute,
	}
	
	cache := NewMemoryCache(config, logger)
	defer cache.Close()
	
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	
	stats := cache.GetStats()
	
	assert.Equal(t, 2, stats["size"])
	assert.Equal(t, 10, stats["maxSize"])
	assert.Equal(t, "1m0s", stats["ttl"])
} 