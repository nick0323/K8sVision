# 缓存模块文档

缓存模块提供高性能的内存缓存功能，用于优化K8sVision的数据访问性能。

## 📁 模块结构

```
cache/
├── README.md                    # 缓存模块文档
├── manager.go                   # 缓存管理器
└── memory.go                    # 内存缓存实现
```

## 🔧 核心组件

### 1. 缓存管理器 (manager.go)
提供统一的缓存管理接口：

```go
type Manager interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
    Clear() error
    Size() int
    Stats() CacheStats
    Close() error
}
```

**主要功能：**
- 缓存操作接口
- 统计信息收集
- 生命周期管理
- 错误处理

### 2. 内存缓存实现 (memory.go)
基于内存的缓存实现：

```go
type MemoryCache struct {
    data    map[string]*CacheItem
    mutex   sync.RWMutex
    stats   CacheStats
    maxSize int
    ttl     time.Duration
    stopCh  chan struct{}
}
```

**特性：**
- 线程安全
- TTL支持
- 容量限制
- 自动清理

## 🚀 主要功能

### 缓存操作

#### 获取缓存
```go
func (c *MemoryCache) Get(key string) (interface{}, bool) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    
    item, exists := c.data[key]
    if !exists {
        c.stats.Misses++
        return nil, false
    }
    
    if item.IsExpired() {
        delete(c.data, key)
        c.stats.Misses++
        return nil, false
    }
    
    c.stats.Hits++
    return item.Value, true
}
```

#### 设置缓存
```go
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    if c.maxSize > 0 && len(c.data) >= c.maxSize {
        c.evictOldest()
    }
    
    item := &CacheItem{
        Value:     value,
        ExpiresAt: time.Now().Add(ttl),
    }
    
    c.data[key] = item
    c.stats.Sets++
    return nil
}
```

#### 删除缓存
```go
func (c *MemoryCache) Delete(key string) error {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    if _, exists := c.data[key]; exists {
        delete(c.data, key)
        c.stats.Deletes++
    }
    
    return nil
}
```

### 统计信息

#### 缓存统计
```go
type CacheStats struct {
    Hits      int64     `json:"hits"`
    Misses    int64     `json:"misses"`
    Sets      int64     `json:"sets"`
    Deletes   int64     `json:"deletes"`
    Evictions int64     `json:"evictions"`
    Size      int       `json:"size"`
    HitRate   float64   `json:"hitRate"`
    LastReset time.Time `json:"lastReset"`
}
```

#### 统计计算
```go
func (c *MemoryCache) Stats() CacheStats {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    
    stats := c.stats
    stats.Size = len(c.data)
    
    total := stats.Hits + stats.Misses
    if total > 0 {
        stats.HitRate = float64(stats.Hits) / float64(total)
    }
    
    return stats
}
```

### 自动清理

#### 过期清理
```go
func (c *MemoryCache) cleanup() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    now := time.Now()
    for key, item := range c.data {
        if item.IsExpired() {
            delete(c.data, key)
            c.stats.Evictions++
        }
    }
}
```

#### 定期清理
```go
func (c *MemoryCache) startCleanup() {
    ticker := time.NewTicker(c.cleanupInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            c.cleanup()
        case <-c.stopCh:
            return
        }
    }
}
```

## ⚙️ 配置选项

### 缓存配置
```go
type CacheConfig struct {
    Enabled         bool          `yaml:"enabled" json:"enabled"`
    Type            string        `yaml:"type" json:"type"`
    TTL             time.Duration `yaml:"ttl" json:"ttl"`
    MaxSize         int           `yaml:"maxSize" json:"maxSize"`
    CleanupInterval time.Duration `yaml:"cleanupInterval" json:"cleanupInterval"`
}
```

### 配置示例
```yaml
cache:
  enabled: true
  type: "memory"
  ttl: "5m"
  maxSize: 1000
  cleanupInterval: "10m"
```

## 🔒 安全特性

### 线程安全
- 使用读写锁保护并发访问
- 原子操作更新统计信息
- 避免竞态条件

### 内存保护
- 容量限制防止内存溢出
- 自动清理过期数据
- 优雅关闭机制

### 数据隔离
- 键值对独立存储
- 类型安全的接口
- 避免数据污染

## 📊 性能优化

### 内存优化
- 使用指针减少内存复制
- 及时清理过期数据
- 限制缓存大小

### 并发优化
- 读写锁分离
- 减少锁持有时间
- 批量操作支持

### 算法优化
- LRU淘汰策略
- 哈希表快速查找
- 时间轮清理机制

## 🛠️ 使用指南

### 基本使用
```go
// 创建缓存管理器
cache := NewMemoryCache(1000, 5*time.Minute)

// 设置缓存
err := cache.Set("key1", "value1", 10*time.Minute)
if err != nil {
    log.Printf("Failed to set cache: %v", err)
}

// 获取缓存
value, exists := cache.Get("key1")
if exists {
    log.Printf("Cache hit: %v", value)
} else {
    log.Printf("Cache miss")
}

// 删除缓存
err = cache.Delete("key1")
if err != nil {
    log.Printf("Failed to delete cache: %v", err)
}
```

### 高级使用
```go
// 批量操作
func BatchSet(cache Manager, items map[string]interface{}, ttl time.Duration) error {
    for key, value := range items {
        if err := cache.Set(key, value, ttl); err != nil {
            return err
        }
    }
    return nil
}

// 条件获取
func GetWithFallback(cache Manager, key string, fallback func() (interface{}, error)) (interface{}, error) {
    if value, exists := cache.Get(key); exists {
        return value, nil
    }
    
    value, err := fallback()
    if err != nil {
        return nil, err
    }
    
    cache.Set(key, value, 5*time.Minute)
    return value, nil
}
```

### 统计监控
```go
// 获取缓存统计
stats := cache.Stats()
log.Printf("Cache Stats: Hits=%d, Misses=%d, HitRate=%.2f", 
    stats.Hits, stats.Misses, stats.HitRate)

// 重置统计
cache.ResetStats()
```

## 🧪 测试

### 单元测试
```go
func TestMemoryCache(t *testing.T) {
    cache := NewMemoryCache(100, time.Minute)
    
    // 测试设置和获取
    err := cache.Set("key1", "value1", time.Minute)
    assert.NoError(t, err)
    
    value, exists := cache.Get("key1")
    assert.True(t, exists)
    assert.Equal(t, "value1", value)
    
    // 测试过期
    time.Sleep(2 * time.Minute)
    _, exists = cache.Get("key1")
    assert.False(t, exists)
}
```

### 性能测试
```go
func BenchmarkMemoryCache(b *testing.B) {
    cache := NewMemoryCache(10000, time.Hour)
    
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            key := fmt.Sprintf("key%d", i%1000)
            cache.Set(key, fmt.Sprintf("value%d", i), time.Hour)
            cache.Get(key)
            i++
        }
    })
}
```

## 📝 最佳实践

1. **容量规划**
   - 根据内存大小设置最大容量
   - 监控缓存命中率
   - 调整TTL策略

2. **键设计**
   - 使用有意义的键名
   - 避免键冲突
   - 考虑键的长度

3. **错误处理**
   - 检查操作返回值
   - 记录错误日志
   - 提供降级方案

4. **监控告警**
   - 监控命中率
   - 设置容量告警
   - 跟踪性能指标

## 🔍 故障排查

### 常见问题

#### 1. 内存使用过高
**症状**: 系统内存使用率持续增长

**排查步骤**:
```go
// 检查缓存大小
stats := cache.Stats()
log.Printf("Cache size: %d", stats.Size)

// 检查配置
log.Printf("Max size: %d", cache.maxSize)
log.Printf("TTL: %v", cache.ttl)
```

**解决方案**:
- 减少最大容量
- 缩短TTL时间
- 增加清理频率

#### 2. 命中率过低
**症状**: 缓存命中率低于预期

**排查步骤**:
```go
// 检查统计信息
stats := cache.Stats()
log.Printf("Hit rate: %.2f", stats.HitRate)
log.Printf("Hits: %d, Misses: %d", stats.Hits, stats.Misses)
```

**解决方案**:
- 延长TTL时间
- 优化键设计
- 检查数据访问模式

#### 3. 并发问题
**症状**: 出现竞态条件或死锁

**排查步骤**:
```go
// 使用race detector
go run -race main.go

// 检查锁使用
// 确保锁的正确获取和释放
```

**解决方案**:
- 检查锁的获取顺序
- 减少锁持有时间
- 使用原子操作

### 调试工具

#### 统计信息
```go
func PrintCacheStats(cache Manager) {
    stats := cache.Stats()
    fmt.Printf("Cache Statistics:\n")
    fmt.Printf("  Size: %d\n", stats.Size)
    fmt.Printf("  Hits: %d\n", stats.Hits)
    fmt.Printf("  Misses: %d\n", stats.Misses)
    fmt.Printf("  Hit Rate: %.2f%%\n", stats.HitRate*100)
    fmt.Printf("  Sets: %d\n", stats.Sets)
    fmt.Printf("  Deletes: %d\n", stats.Deletes)
    fmt.Printf("  Evictions: %d\n", stats.Evictions)
}
```

#### 性能分析
```go
func ProfileCache(cache Manager, duration time.Duration) {
    start := time.Now()
    operations := 0
    
    for time.Since(start) < duration {
        key := fmt.Sprintf("key%d", operations%1000)
        cache.Set(key, fmt.Sprintf("value%d", operations), time.Hour)
        cache.Get(key)
        operations++
    }
    
    elapsed := time.Since(start)
    opsPerSec := float64(operations) / elapsed.Seconds()
    
    fmt.Printf("Operations: %d\n", operations)
    fmt.Printf("Duration: %v\n", elapsed)
    fmt.Printf("Ops/sec: %.2f\n", opsPerSec)
}
```

## 🔄 扩展功能

### 持久化缓存
可以扩展支持持久化存储：
- Redis缓存
- 文件系统缓存
- 数据库缓存

### 分布式缓存
支持多节点缓存：
- 一致性哈希
- 节点发现
- 数据同步

### 高级特性
- 缓存预热
- 缓存穿透保护
- 缓存雪崩防护
- 热点数据识别
