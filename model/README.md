# 模型层文档

模型层定义了K8sVision中使用的数据结构、常量和配置，为整个应用提供统一的数据模型。

## 📁 模块结构

```
model/
├── README.md                    # 模型层文档
├── types.go                     # 数据类型定义
├── config.go                    # 配置结构定义
└── consts.go                    # 常量定义
```

## 🔧 核心组件

### 1. 数据类型定义 (types.go)
定义了应用中使用的所有数据结构：

#### API响应结构
```go
type APIResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

#### 资源状态结构
为每个Kubernetes资源定义了标准化的状态结构：

**Pod状态**
```go
type PodStatus struct {
    Name        string            `json:"name"`
    Namespace   string            `json:"namespace"`
    Status      string            `json:"status"`
    Ready       string            `json:"ready"`
    Restarts    int32             `json:"restarts"`
    Age         string            `json:"age"`
    Node        string            `json:"node"`
    IP          string            `json:"ip"`
    Labels      map[string]string `json:"labels"`
    Annotations map[string]string `json:"annotations"`
    // ... 更多字段
}
```

**Deployment状态**
```go
type DeploymentStatus struct {
    Name              string            `json:"name"`
    Namespace         string            `json:"namespace"`
    ReadyReplicas     int32             `json:"readyReplicas"`
    DesiredReplicas   int32             `json:"desiredReplicas"`
    AvailableReplicas int32             `json:"availableReplicas"`
    Age               string            `json:"age"`
    Labels            map[string]string `json:"labels"`
    // ... 更多字段
}
```

#### 集群概览结构
```go
type OverviewStatus struct {
    ClusterInfo    ClusterInfo     `json:"clusterInfo"`
    NodeStats      NodeStats       `json:"nodeStats"`
    ResourceStats  ResourceStats   `json:"resourceStats"`
    HealthStatus   HealthStatus    `json:"healthStatus"`
    RecentEvents   []EventSummary  `json:"recentEvents"`
}

type ClusterInfo struct {
    Name        string `json:"name"`
    Version     string `json:"version"`
    NodeCount   int    `json:"nodeCount"`
    NamespaceCount int `json:"namespaceCount"`
}
```

#### 认证相关结构
```go
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token   string    `json:"token"`
    Expires time.Time `json:"expires"`
}
```

### 2. 配置结构定义 (config.go)
定义了应用的配置结构：

#### 服务器配置
```go
type ServerConfig struct {
    Port         string `yaml:"port" json:"port"`
    Host         string `yaml:"host" json:"host"`
    ReadTimeout  string `yaml:"readTimeout" json:"readTimeout"`
    WriteTimeout string `yaml:"writeTimeout" json:"writeTimeout"`
    IdleTimeout  string `yaml:"idleTimeout" json:"idleTimeout"`
    MaxHeaderBytes int  `yaml:"maxHeaderBytes" json:"maxHeaderBytes"`
}
```

#### Kubernetes配置
```go
type KubernetesConfig struct {
    Kubeconfig string `yaml:"kubeconfig" json:"kubeconfig"`
    Context    string `yaml:"context" json:"context"`
    Timeout    string `yaml:"timeout" json:"timeout"`
    QPS        int    `yaml:"qps" json:"qps"`
    Burst      int    `yaml:"burst" json:"burst"`
    Insecure   bool   `yaml:"insecure" json:"insecure"`
    CAFile     string `yaml:"caFile" json:"caFile"`
    CertFile   string `yaml:"certFile" json:"certFile"`
    KeyFile    string `yaml:"keyFile" json:"keyFile"`
    Token      string `yaml:"token" json:"token"`
    APIServer  string `yaml:"apiServer" json:"apiServer"`
}
```

#### JWT配置
```go
type JWTConfig struct {
    Secret     string `yaml:"secret" json:"secret"`
    Expiration string `yaml:"expiration" json:"expiration"`
    Issuer     string `yaml:"issuer" json:"issuer"`
    Audience   string `yaml:"audience" json:"audience"`
}
```

#### 日志配置
```go
type LogConfig struct {
    Level      string `yaml:"level" json:"level"`
    Format     string `yaml:"format" json:"format"`
    Output     string `yaml:"output" json:"output"`
    MaxSize    int    `yaml:"maxSize" json:"maxSize"`
    MaxBackups int    `yaml:"maxBackups" json:"maxBackups"`
    MaxAge     int    `yaml:"maxAge" json:"maxAge"`
    Compress   bool   `yaml:"compress" json:"compress"`
}
```

#### 认证配置
```go
type AuthConfig struct {
    Username        string `yaml:"username" json:"username"`
    Password        string `yaml:"password" json:"password"`
    MaxLoginFail    int    `yaml:"maxLoginFail" json:"maxLoginFail"`
    LockDuration    string `yaml:"lockDuration" json:"lockDuration"`
    SessionTimeout  string `yaml:"sessionTimeout" json:"sessionTimeout"`
    EnableRateLimit bool   `yaml:"enableRateLimit" json:"enableRateLimit"`
    RateLimit       int    `yaml:"rateLimit" json:"rateLimit"`
}
```

#### 缓存配置
```go
type CacheConfig struct {
    Enabled         bool   `yaml:"enabled" json:"enabled"`
    Type            string `yaml:"type" json:"type"`
    TTL             string `yaml:"ttl" json:"ttl"`
    MaxSize         int    `yaml:"maxSize" json:"maxSize"`
    CleanupInterval string `yaml:"cleanupInterval" json:"cleanupInterval"`
}
```

#### 主配置结构
```go
type Config struct {
    Server     ServerConfig     `yaml:"server" json:"server"`
    Kubernetes KubernetesConfig `yaml:"kubernetes" json:"kubernetes"`
    JWT        JWTConfig        `yaml:"jwt" json:"jwt"`
    Log        LogConfig        `yaml:"log" json:"log"`
    Auth       AuthConfig       `yaml:"auth" json:"auth"`
    Cache      CacheConfig      `yaml:"cache" json:"cache"`
}
```

### 3. 常量定义 (consts.go)
定义了应用中使用的常量：

#### 响应状态码
```go
const (
    CodeSuccess      = 200
    CodeBadRequest   = 400
    CodeUnauthorized = 401
    CodeForbidden    = 403
    CodeNotFound     = 404
    CodeConflict     = 409
    CodeTooManyRequests = 429
    CodeInternalError = 500
    CodeServiceUnavailable = 503
)
```

#### 响应消息
```go
const (
    SuccessMessage         = "success"
    ListSuccessMessage    = "列表获取成功"
    DetailSuccessMessage  = "详情获取成功"
    CreateSuccessMessage  = "创建成功"
    UpdateSuccessMessage  = "更新成功"
    DeleteSuccessMessage  = "删除成功"
    
    ErrorMessageBadRequest    = "请求参数错误"
    ErrorMessageUnauthorized  = "未授权访问"
    ErrorMessageForbidden     = "禁止访问"
    ErrorMessageNotFound      = "资源不存在"
    ErrorMessageInternalError = "内部服务器错误"
)
```

#### 资源状态
```go
const (
    StatusRunning    = "Running"
    StatusPending    = "Pending"
    StatusSucceeded  = "Succeeded"
    StatusFailed     = "Failed"
    StatusUnknown    = "Unknown"
    StatusTerminating = "Terminating"
)
```

#### 时间格式
```go
const (
    TimeFormatDefault = "2006-01-02 15:04:05"
    TimeFormatISO8601 = "2006-01-02T15:04:05Z"
    TimeFormatRFC3339 = "2006-01-02T15:04:05Z07:00"
)
```

## 🎯 设计原则

### 1. 一致性
- 所有资源状态结构遵循相同的命名规范
- 统一的字段类型和格式
- 一致的JSON标签定义

### 2. 可扩展性
- 使用接口定义通用行为
- 支持动态字段扩展
- 版本兼容性考虑

### 3. 类型安全
- 使用强类型定义
- 避免使用`interface{}`除非必要
- 提供类型转换方法

### 4. 性能优化
- 避免深层嵌套结构
- 使用指针减少内存占用
- 提供序列化优化

## 🛠️ 使用指南

### 创建资源状态
```go
func CreatePodStatus(pod *v1.Pod) *PodStatus {
    return &PodStatus{
        Name:        pod.Name,
        Namespace:   pod.Namespace,
        Status:      determinePodStatus(pod),
        Ready:       calculateReadyStatus(pod),
        Restarts:    calculateRestarts(pod),
        Age:         calculateAge(pod.CreationTimestamp),
        Node:        pod.Spec.NodeName,
        IP:          pod.Status.PodIP,
        Labels:      pod.Labels,
        Annotations: pod.Annotations,
    }
}
```

### 处理API响应
```go
func HandleAPIResponse(data interface{}, message string) *APIResponse {
    return &APIResponse{
        Code:    CodeSuccess,
        Message: message,
        Data:    data,
    }
}

func HandleAPIError(err error, code int) *APIResponse {
    return &APIResponse{
        Code:    code,
        Message: ErrorMessageInternalError,
        Error: &APIError{
            Code:    code,
            Message: err.Error(),
        },
    }
}
```

### 配置验证
```go
func (c *Config) Validate() error {
    if c.Server.Port == "" {
        return errors.New("server port is required")
    }
    
    if c.JWT.Secret == "" {
        return errors.New("JWT secret is required")
    }
    
    if c.Auth.Username == "" {
        return errors.New("auth username is required")
    }
    
    return nil
}
```

## 📊 数据转换

### Kubernetes对象转换
提供从Kubernetes原生对象到应用模型的转换方法：

```go
func ConvertPodToStatus(pod *v1.Pod) *PodStatus {
    return &PodStatus{
        Name:        pod.Name,
        Namespace:   pod.Namespace,
        Status:      getPodStatus(pod),
        Ready:       getReadyStatus(pod),
        Restarts:    getRestartCount(pod),
        Age:         getAge(pod.CreationTimestamp),
        Node:        pod.Spec.NodeName,
        IP:          pod.Status.PodIP,
        Labels:      pod.Labels,
        Annotations: pod.Annotations,
    }
}
```

### 时间格式化
提供统一的时间格式化方法：

```go
func FormatTime(t *metav1.Time) string {
    if t == nil {
        return ""
    }
    return t.Format(TimeFormatDefault)
}

func FormatDuration(start, end time.Time) string {
    duration := end.Sub(start)
    if duration < time.Minute {
        return fmt.Sprintf("%.0fs", duration.Seconds())
    } else if duration < time.Hour {
        return fmt.Sprintf("%.0fm", duration.Minutes())
    } else {
        return fmt.Sprintf("%.0fh", duration.Hours())
    }
}
```

## 🔒 安全考虑

### 敏感数据处理
- 密码字段使用指针类型，避免意外序列化
- 提供脱敏方法
- 支持加密存储

### 输入验证
- 提供验证标签
- 实现验证方法
- 防止注入攻击

### 权限控制
- 定义权限级别常量
- 提供权限检查方法
- 支持细粒度控制

## 📝 最佳实践

1. **结构设计**
   - 保持结构简洁
   - 使用有意义的字段名
   - 避免过度嵌套

2. **类型选择**
   - 优先使用具体类型
   - 合理使用指针类型
   - 避免不必要的`interface{}`

3. **性能考虑**
   - 使用结构体标签优化序列化
   - 避免深层复制
   - 提供批量处理方法

4. **可维护性**
   - 添加必要的注释
   - 提供示例用法
   - 保持向后兼容

## 🧪 测试

### 单元测试
为每个结构体提供测试用例：
- 字段验证测试
- 转换方法测试
- 边界条件测试

### 集成测试
- 配置加载测试
- 数据序列化测试
- 性能基准测试

## 🔍 故障排查

### 常见问题
1. **序列化错误**
   - 检查JSON标签
   - 验证字段类型
   - 查看错误日志

2. **配置加载失败**
   - 检查配置文件格式
   - 验证必需字段
   - 查看环境变量

3. **类型转换错误**
   - 检查源数据类型
   - 验证转换逻辑
   - 添加错误处理

### 调试工具
- 结构体打印
- 类型检查
- 序列化测试
- 性能分析
