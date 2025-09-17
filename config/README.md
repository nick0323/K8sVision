# 配置模块文档

配置模块负责K8sVision的配置管理，包括配置文件加载、环境变量处理、配置验证和热重载功能。

## 📁 模块结构

```
config/
├── README.md                    # 模块文档
└── manager.go                   # 配置管理器
```

## 🔧 核心组件

### 配置管理器 (manager.go)
提供统一的配置管理功能：
- 配置文件加载
- 环境变量覆盖
- 配置验证
- 热重载支持

**主要结构：**
```go
type Manager struct {
    config *model.Config
    viper  *viper.Viper
    mutex  sync.RWMutex
}
```

## ⚙️ 配置结构

### 服务器配置
```yaml
server:
  port: "8080"              # 服务端口
  host: "0.0.0.0"           # 绑定地址
  readTimeout: "30s"        # 读取超时
  writeTimeout: "30s"       # 写入超时
  idleTimeout: "60s"        # 空闲超时
  maxHeaderBytes: 1048576   # 最大请求头大小
```

### Kubernetes配置
```yaml
kubernetes:
  kubeconfig: ""            # Kubeconfig文件路径
  context: ""               # 集群上下文
  timeout: "30s"            # 请求超时
  qps: 100                  # 每秒查询数
  burst: 200                # 突发请求数
  insecure: true            # 是否跳过TLS验证
  caFile: ""                # CA证书文件
  certFile: ""              # 客户端证书
  keyFile: ""               # 客户端私钥
  token: ""                 # 访问令牌
  apiServer: ""             # API服务器地址
```

### JWT配置
```yaml
jwt:
  secret: ""                # JWT密钥
  expiration: "24h"         # 令牌过期时间
  issuer: "k8svision"       # 签发者
  audience: "k8svision-client" # 受众
```

### 日志配置
```yaml
log:
  level: "info"             # 日志级别
  format: "json"            # 日志格式
  output: "stdout"          # 输出目标
  maxSize: 100              # 最大文件大小(MB)
  maxBackups: 3             # 最大备份数
  maxAge: 28                # 最大保存天数
  compress: true            # 是否压缩
```

### 认证配置
```yaml
auth:
  username: ""              # 管理员用户名
  password: ""              # 管理员密码
  maxLoginFail: 5           # 最大登录失败次数
  lockDuration: "10m"       # 锁定持续时间
  sessionTimeout: "24h"     # 会话超时时间
  enableRateLimit: true     # 是否启用速率限制
  rateLimit: 100            # 速率限制(请求/分钟)
```

### 缓存配置
```yaml
cache:
  enabled: true             # 是否启用缓存
  type: "memory"            # 缓存类型
  ttl: "5m"                 # 缓存生存时间
  maxSize: 1000             # 最大缓存条目数
  cleanupInterval: "10m"    # 清理间隔
```

## 🚀 主要功能

### 配置加载
```go
func (m *Manager) LoadConfig(configFile string) error
```
- 从配置文件加载配置
- 支持YAML格式
- 自动处理环境变量覆盖

### 配置获取
```go
func (m *Manager) GetConfig() *model.Config
```
- 线程安全的配置获取
- 返回当前配置副本
- 支持并发访问

### 配置验证
```go
func (m *Manager) ValidateConfig() error
```
- 验证配置完整性
- 检查必需字段
- 验证配置值范围

### 热重载
```go
func (m *Manager) WatchConfig() error
```
- 监控配置文件变化
- 自动重新加载配置
- 通知配置变更

## 🔒 安全特性

### 敏感信息保护
- 密码加密存储
- 密钥安全处理
- 配置信息脱敏

### 环境变量支持
- 支持环境变量覆盖
- 敏感信息通过环境变量设置
- 避免硬编码敏感信息

### 配置验证
- 输入验证
- 类型检查
- 范围验证

## 📊 环境变量

### 必需环境变量
```bash
# JWT密钥 (至少32位字符)
export K8SVISION_JWT_SECRET="your-32-character-secret-key"

# 管理员用户名
export K8SVISION_AUTH_USERNAME="admin"

# 管理员密码
export K8SVISION_AUTH_PASSWORD="your-password"
```

### 可选环境变量
```bash
# Kubeconfig文件路径
export K8SVISION_KUBECONFIG="/path/to/kubeconfig"

# 日志级别
export K8SVISION_LOG_LEVEL="info"

# 服务端口
export K8SVISION_PORT="8080"

# 数据库连接
export K8SVISION_DATABASE_URL="postgres://user:pass@localhost/db"
```

## 🛠️ 开发指南

### 添加新配置项
1. 在`model/types.go`中定义配置结构
2. 在配置文件中添加默认值
3. 在`manager.go`中添加验证逻辑
4. 更新文档

### 配置验证示例
```go
func (m *Manager) validateServerConfig(cfg *model.ServerConfig) error {
    if cfg.Port == "" {
        return errors.New("server port is required")
    }
    
    if port, err := strconv.Atoi(cfg.Port); err != nil || port < 1 || port > 65535 {
        return errors.New("invalid server port")
    }
    
    return nil
}
```

### 环境变量处理
```go
func (m *Manager) loadFromEnv() {
    if secret := os.Getenv("K8SVISION_JWT_SECRET"); secret != "" {
        m.config.JWT.Secret = secret
    }
    
    if username := os.Getenv("K8SVISION_AUTH_USERNAME"); username != "" {
        m.config.Auth.Username = username
    }
}
```

## 📝 最佳实践

1. **配置设计**
   - 提供合理的默认值
   - 支持环境变量覆盖
   - 保持配置结构清晰

2. **安全考虑**
   - 敏感信息通过环境变量设置
   - 避免在配置文件中硬编码密码
   - 验证配置输入

3. **性能优化**
   - 使用读写锁保护配置
   - 避免频繁的配置重新加载
   - 缓存配置值

4. **可维护性**
   - 保持配置结构简单
   - 提供清晰的配置文档
   - 支持配置验证

## 🔍 故障排查

### 常见问题
1. **配置文件加载失败**
   - 检查文件路径
   - 验证文件格式
   - 查看文件权限

2. **环境变量未生效**
   - 检查环境变量名称
   - 确认环境变量值
   - 重启应用程序

3. **配置验证失败**
   - 检查必需字段
   - 验证配置值格式
   - 查看错误日志

### 调试工具
- 配置打印功能
- 验证错误日志
- 环境变量检查
- 配置文件语法检查

## 📚 配置示例

### 开发环境配置
```yaml
server:
  port: "8080"
  host: "localhost"

kubernetes:
  kubeconfig: "~/.kube/config"
  context: "minikube"

log:
  level: "debug"
  format: "console"

auth:
  username: "admin"
  password: "admin123"
```

### 生产环境配置
```yaml
server:
  port: "8080"
  host: "0.0.0.0"

kubernetes:
  kubeconfig: "/etc/k8s/kubeconfig"
  context: "production"

log:
  level: "info"
  format: "json"

auth:
  username: ""  # 通过环境变量设置
  password: ""  # 通过环境变量设置
```
