# 监控模块文档

监控模块提供K8sVision的指标收集、性能监控和链路追踪功能，确保系统的可观测性和稳定性。

## 📁 模块结构

```
monitor/
├── README.md                    # 监控模块文档
├── metrics.go                   # 系统指标收集
├── business_metrics.go          # 业务指标收集
└── tracing.go                   # 链路追踪
```

## 🔧 核心组件

### 1. 系统指标收集 (metrics.go)
负责收集系统级别的性能指标：

```go
type MetricsManager struct {
    metrics map[string]interface{}
    mutex   sync.RWMutex
    startTime time.Time
}
```

**主要功能：**
- CPU使用率监控
- 内存使用情况
- 网络I/O统计
- 连接数监控
- 响应时间统计

### 2. 业务指标收集 (business_metrics.go)
收集应用特定的业务指标：

```go
type BusinessMetricsCollector struct {
    metrics []BusinessMetric
    mutex   sync.RWMutex
}

type BusinessMetric struct {
    Name        string                 `json:"name"`
    Value       float64                `json:"value"`
    Unit        string                 `json:"unit"`
    MetricType  string                 `json:"metricType"`
    Labels      map[string]string      `json:"labels"`
    Timestamp   time.Time              `json:"timestamp"`
    Description string                 `json:"description"`
}
```

**业务指标类型：**
- 集群资源使用率
- API请求统计
- 用户活跃度
- 错误率统计
- 缓存命中率

### 3. 链路追踪 (tracing.go)
提供分布式链路追踪功能：

```go
type TraceContext struct {
    TraceID    string `json:"traceId"`
    SpanID     string `json:"spanId"`
    ParentID   string `json:"parentId,omitempty"`
    StartTime  time.Time `json:"startTime"`
    EndTime    time.Time `json:"endTime,omitempty"`
    Duration   time.Duration `json:"duration,omitempty"`
    Operation  string `json:"operation"`
    Tags       map[string]string `json:"tags"`
    Logs       []TraceLog `json:"logs,omitempty"`
}
```

## 🚀 主要功能

### 系统指标收集

#### CPU监控
```go
func (m *MetricsManager) collectCPU() {
    // 获取CPU使用率
    cpuUsage := getCPUUsage()
    m.setMetric("cpu", cpuUsage)
    
    // 获取CPU核心数
    cpuCores := getCPUCores()
    m.setMetric("cpu_cores", cpuCores)
    
    // 获取负载平均值
    loadAvg := getLoadAverage()
    m.setMetric("load_average", loadAvg)
}
```

#### 内存监控
```go
func (m *MetricsManager) collectMemory() {
    // 获取内存使用情况
    memStats := getMemoryStats()
    m.setMetric("memory", map[string]interface{}{
        "total":     memStats.Total,
        "used":      memStats.Used,
        "free":      memStats.Free,
        "available": memStats.Available,
        "usage_percent": memStats.UsagePercent,
    })
    
    // 获取交换空间使用情况
    swapStats := getSwapStats()
    m.setMetric("swap", swapStats)
}
```

#### 网络监控
```go
func (m *MetricsManager) collectNetwork() {
    // 获取网络接口统计
    netStats := getNetworkStats()
    m.setMetric("network", map[string]interface{}{
        "bytes_sent":     netStats.BytesSent,
        "bytes_recv":     netStats.BytesRecv,
        "packets_sent":   netStats.PacketsSent,
        "packets_recv":   netStats.PacketsRecv,
        "err_in":         netStats.ErrIn,
        "err_out":        netStats.ErrOut,
        "drop_in":        netStats.DropIn,
        "drop_out":       netStats.DropOut,
    })
}
```

### 业务指标收集

#### 集群资源指标
```go
func (c *BusinessMetricsCollector) collectClusterMetrics() {
    // 节点资源使用率
    nodeMetrics := c.getNodeResourceUsage()
    for _, metric := range nodeMetrics {
        c.addMetric(BusinessMetric{
            Name:        "node_cpu_usage",
            Value:       metric.CPUUsage,
            Unit:        "percent",
            MetricType:  "gauge",
            Labels:      map[string]string{"node": metric.NodeName},
            Timestamp:   time.Now(),
            Description: "Node CPU usage percentage",
        })
    }
    
    // Pod资源使用率
    podMetrics := c.getPodResourceUsage()
    for _, metric := range podMetrics {
        c.addMetric(BusinessMetric{
            Name:        "pod_memory_usage",
            Value:       metric.MemoryUsage,
            Unit:        "bytes",
            MetricType:  "gauge",
            Labels:      map[string]string{
                "pod":       metric.PodName,
                "namespace": metric.Namespace,
            },
            Timestamp:   time.Now(),
            Description: "Pod memory usage in bytes",
        })
    }
}
```

#### API请求指标
```go
func (c *BusinessMetricsCollector) collectAPIMetrics() {
    // 请求总数
    c.addMetric(BusinessMetric{
        Name:        "api_requests_total",
        Value:       float64(c.getTotalRequests()),
        Unit:        "count",
        MetricType:  "counter",
        Labels:      map[string]string{},
        Timestamp:   time.Now(),
        Description: "Total number of API requests",
    })
    
    // 错误率
    errorRate := c.calculateErrorRate()
    c.addMetric(BusinessMetric{
        Name:        "api_error_rate",
        Value:       errorRate,
        Unit:        "percent",
        MetricType:  "gauge",
        Labels:      map[string]string{},
        Timestamp:   time.Now(),
        Description: "API error rate percentage",
    })
}
```

### 链路追踪

#### 创建追踪上下文
```go
func NewTraceContext(operation string) *TraceContext {
    return &TraceContext{
        TraceID:   generateTraceID(),
        SpanID:    generateSpanID(),
        StartTime: time.Now(),
        Operation: operation,
        Tags:      make(map[string]string),
        Logs:      make([]TraceLog, 0),
    }
}
```

#### 添加追踪日志
```go
func (t *TraceContext) AddLog(level string, message string, fields map[string]interface{}) {
    log := TraceLog{
        Timestamp: time.Now(),
        Level:     level,
        Message:   message,
        Fields:    fields,
    }
    
    t.mutex.Lock()
    defer t.mutex.Unlock()
    t.Logs = append(t.Logs, log)
}
```

#### 完成追踪
```go
func (t *TraceContext) Finish() {
    t.EndTime = time.Now()
    t.Duration = t.EndTime.Sub(t.StartTime)
    
    // 发送到追踪系统
    t.sendToTracingSystem()
}
```

## 📊 指标类型

### 系统指标
- **CPU指标**: 使用率、核心数、负载平均值
- **内存指标**: 总量、使用量、可用量、使用率
- **网络指标**: 发送/接收字节数、包数、错误数
- **磁盘指标**: 使用量、IOPS、读写速度
- **进程指标**: 进程数、线程数、文件描述符

### 业务指标
- **集群指标**: 节点数、Pod数、资源使用率
- **API指标**: 请求数、响应时间、错误率
- **用户指标**: 活跃用户数、会话数
- **缓存指标**: 命中率、大小、操作数
- **数据库指标**: 连接数、查询数、响应时间

### 自定义指标
- **应用指标**: 业务逻辑相关指标
- **性能指标**: 关键路径性能
- **质量指标**: 错误率、可用性
- **容量指标**: 资源容量规划

## 🔒 安全特性

### 数据脱敏
- 敏感信息过滤
- 用户数据保护
- 配置信息隐藏

### 访问控制
- 指标访问权限
- 追踪数据保护
- 审计日志记录

### 数据加密
- 传输加密
- 存储加密
- 密钥管理

## 📈 性能优化

### 指标收集优化
- 异步收集
- 批量处理
- 采样策略

### 存储优化
- 数据压缩
- 过期清理
- 分片存储

### 查询优化
- 索引优化
- 缓存策略
- 聚合计算

## 🛠️ 使用指南

### 基本使用
```go
// 获取指标管理器
metricsManager := GetMetricsManager()

// 收集系统指标
err := metricsManager.CollectMetrics()
if err != nil {
    log.Printf("Failed to collect metrics: %v", err)
}

// 获取所有指标
allMetrics := metricsManager.GetAllMetrics()
fmt.Printf("System metrics: %+v\n", allMetrics)
```

### 业务指标使用
```go
// 获取业务指标收集器
collector := GetBusinessMetricsCollector()

// 收集业务指标
err := collector.CollectMetrics()
if err != nil {
    log.Printf("Failed to collect business metrics: %v", err)
}

// 获取指标列表
metrics := collector.GetMetrics()
for _, metric := range metrics {
    fmt.Printf("Metric: %s = %.2f %s\n", 
        metric.Name, metric.Value, metric.Unit)
}
```

### 链路追踪使用
```go
// 创建追踪上下文
trace := NewTraceContext("api_request")
defer trace.Finish()

// 添加标签
trace.AddTag("method", "GET")
trace.AddTag("path", "/api/pods")
trace.AddTag("user_id", "12345")

// 添加日志
trace.AddLog("info", "Processing request", map[string]interface{}{
    "request_id": "req-123",
    "duration":   "100ms",
})

// 执行业务逻辑
result, err := processRequest()
if err != nil {
    trace.AddTag("error", "true")
    trace.AddLog("error", "Request failed", map[string]interface{}{
        "error": err.Error(),
    })
}
```

## 🧪 测试

### 单元测试
```go
func TestMetricsManager(t *testing.T) {
    manager := NewMetricsManager()
    
    // 测试指标收集
    err := manager.CollectMetrics()
    assert.NoError(t, err)
    
    // 测试指标获取
    metrics := manager.GetAllMetrics()
    assert.NotEmpty(t, metrics)
    
    // 验证特定指标
    cpu, exists := metrics["cpu"]
    assert.True(t, exists)
    assert.IsType(t, float64(0), cpu)
}
```

### 性能测试
```go
func BenchmarkMetricsCollection(b *testing.B) {
    manager := NewMetricsManager()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        manager.CollectMetrics()
    }
}
```

## 📝 最佳实践

1. **指标设计**
   - 使用有意义的指标名称
   - 遵循命名规范
   - 提供清晰的描述

2. **收集策略**
   - 合理设置收集频率
   - 避免过度收集
   - 使用采样策略

3. **存储管理**
   - 设置数据保留期
   - 定期清理旧数据
   - 监控存储使用量

4. **告警配置**
   - 设置合理的阈值
   - 避免告警风暴
   - 提供降级策略

## 🔍 故障排查

### 常见问题

#### 1. 指标收集失败
**症状**: 指标数据不更新或为空

**排查步骤**:
```go
// 检查收集器状态
if !collector.IsRunning() {
    log.Printf("Collector is not running")
}

// 检查错误日志
logs := collector.GetErrorLogs()
for _, log := range logs {
    log.Printf("Collection error: %v", log)
}
```

**解决方案**:
- 检查系统权限
- 验证配置参数
- 重启收集器

#### 2. 性能问题
**症状**: 指标收集影响应用性能

**排查步骤**:
```go
// 检查收集频率
interval := collector.GetCollectionInterval()
log.Printf("Collection interval: %v", interval)

// 检查指标数量
count := collector.GetMetricsCount()
log.Printf("Metrics count: %d", count)
```

**解决方案**:
- 降低收集频率
- 减少指标数量
- 使用异步收集

#### 3. 存储空间不足
**症状**: 指标数据无法存储

**排查步骤**:
```go
// 检查存储使用量
usage := storage.GetUsage()
log.Printf("Storage usage: %d bytes", usage)

// 检查数据保留期
retention := storage.GetRetentionPeriod()
log.Printf("Retention period: %v", retention)
```

**解决方案**:
- 缩短数据保留期
- 增加存储容量
- 启用数据压缩

### 调试工具

#### 指标导出
```go
func ExportMetrics(manager *MetricsManager, format string) {
    metrics := manager.GetAllMetrics()
    
    switch format {
    case "json":
        data, _ := json.MarshalIndent(metrics, "", "  ")
        fmt.Println(string(data))
    case "prometheus":
        exportPrometheusFormat(metrics)
    case "influxdb":
        exportInfluxDBFormat(metrics)
    }
}
```

#### 性能分析
```go
func ProfileMetricsCollection(manager *MetricsManager, duration time.Duration) {
    start := time.Now()
    iterations := 0
    
    for time.Since(start) < duration {
        manager.CollectMetrics()
        iterations++
    }
    
    elapsed := time.Since(start)
    avgTime := elapsed / time.Duration(iterations)
    
    fmt.Printf("Iterations: %d\n", iterations)
    fmt.Printf("Total time: %v\n", elapsed)
    fmt.Printf("Average time: %v\n", avgTime)
    fmt.Printf("Iterations/sec: %.2f\n", float64(iterations)/elapsed.Seconds())
}
```

## 🔄 扩展功能

### 指标导出
- Prometheus格式
- InfluxDB格式
- Graphite格式
- 自定义格式

### 告警系统
- 阈值告警
- 趋势告警
- 异常检测
- 通知集成

### 可视化
- 仪表板
- 图表生成
- 实时监控
- 历史分析
