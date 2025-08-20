# 配置说明

本文档详细介绍了 K8sVision 的配置选项和参数说明。

## 📋 配置概览

K8sVision 支持多种配置方式，优先级从高到低为：
1. 环境变量
2. 配置文件 (config.yaml)
3. 默认配置

## ⚙️ 配置文件

### 配置文件位置
- 默认路径: `./config.yaml`
- 自定义路径: 通过 `--config` 参数指定
- 容器内路径: `/app/config.yaml`

### 配置文件结构
```yaml
server:
  port: 8080
  host: "0.0.0.0"
  readTimeout: 30s
  writeTimeout: 30s
  idleTimeout: 60s
  maxHeaderBytes: 1048576

kubernetes:
  kubeconfig: ""
  context: ""
  timeout: 30s
  qps: 100
  burst: 200
  insecure: true
  caFile: ""
  certFile: ""
  keyFile: ""
  token: ""
  apiServer: ""

jwt:
  secret: "k8svision-secret-key"
  expiration: 24h
  issuer: "k8svision"
  audience: "k8svision-client"

log:
  level: "info"
  format: "json"
  output: "stdout"
  maxSize: 100
  maxBackups: 3
  maxAge: 28
  compress: true

auth:
  username: "admin"
  password: "admin"
  maxLoginFail: 5
  lockDuration: 10m
  sessionTimeout: 24h
  enableRateLimit: true
  rateLimit: 100

cache:
  enabled: true
  type: "memory"
  ttl: 5m
  maxSize: 1000
  cleanupInterval: 10m
```

## 🔧 配置参数详解

### 服务器配置 (server)

| 参数 | 类型 | 默认值 | 说明 | 环境变量 |
|------|------|--------|------|----------|
| `port` | string | "8080" | 服务监听端口 | `SERVER_PORT` |
| `host` | string | "0.0.0.0" | 服务监听地址 | `SERVER_HOST` |
| `readTimeout` | duration | 30s | 读取超时时间 | `SERVER_READ_TIMEOUT` |
| `writeTimeout` | duration | 30s | 写入超时时间 | `SERVER_WRITE_TIMEOUT` |
| `idleTimeout` | duration | 60s | 空闲超时时间 | `SERVER_IDLE_TIMEOUT` |
| `maxHeaderBytes` | int | 1048576 | 最大请求头大小 | `SERVER_MAX_HEADER_BYTES` |

**示例**:
```yaml
server:
  port: "9090"
  host: "127.0.0.1"
  readTimeout: 60s
  writeTimeout: 60s
  idleTimeout: 120s
  maxHeaderBytes: 2097152
```

### Kubernetes 配置 (kubernetes)

| 参数 | 类型 | 默认值 | 说明 | 环境变量 |
|------|------|--------|------|----------|
| `kubeconfig` | string | "" | kubeconfig 文件路径 | `KUBECONFIG` |
| `context` | string | "" | Kubernetes 上下文 | `KUBERNETES_CONTEXT` |
| `timeout` | duration | 30s | API 请求超时时间 | `KUBERNETES_TIMEOUT` |
| `qps` | float32 | 100 | 每秒查询数限制 | `KUBERNETES_QPS` |
| `burst` | int | 200 | 突发请求数限制 | `KUBERNETES_BURST` |
| `insecure` | bool | true | 是否跳过 TLS 验证 | `KUBERNETES_INSECURE` |
| `caFile` | string | "" | CA 证书文件路径 | `KUBERNETES_CA_FILE` |
| `certFile` | string | "" | 客户端证书文件路径 | `KUBERNETES_CERT_FILE` |
| `keyFile` | string | "" | 客户端密钥文件路径 | `KUBERNETES_KEY_FILE` |
| `token` | string | "" | 服务账户 Token | `KUBERNETES_TOKEN` |
| `apiServer` | string | "" | API 服务器地址 | `KUBERNETES_API_SERVER` |

**示例**:
```yaml
kubernetes:
  kubeconfig: "/path/to/kubeconfig"
  context: "production"
  timeout: 60s
  qps: 200
  burst: 400
  insecure: false
  caFile: "/path/to/ca.crt"
  certFile: "/path/to/client.crt"
  keyFile: "/path/to/client.key"
```

### JWT 配置 (jwt)

| 参数 | 类型 | 默认值 | 说明 | 环境变量 |
|------|------|--------|------|----------|
| `secret` | string | "k8svision-secret-key" | JWT 签名密钥 | `JWT_SECRET` |
| `expiration` | duration | 24h | Token 过期时间 | `JWT_EXPIRATION` |
| `issuer` | string | "k8svision" | Token 发行者 | `JWT_ISSUER` |
| `audience` | string | "k8svision-client" | Token 受众 | `JWT_AUDIENCE` |

**示例**:
```yaml
jwt:
  secret: "your-super-secret-key-here"
  expiration: 48h
  issuer: "k8svision-prod"
  audience: "k8svision-web-client"
```

### 日志配置 (log)

| 参数 | 类型 | 默认值 | 说明 | 环境变量 |
|------|------|--------|------|----------|
| `level` | string | "info" | 日志级别 | `LOG_LEVEL` |
| `format` | string | "json" | 日志格式 | `LOG_FORMAT` |
| `output` | string | "stdout" | 日志输出 | `LOG_OUTPUT` |
| `maxSize` | int | 100 | 单个日志文件最大大小(MB) | `LOG_MAX_SIZE` |
| `maxBackups` | int | 3 | 保留的日志文件数量 | `LOG_MAX_BACKUPS` |
| `maxAge` | int | 28 | 日志文件保留天数 | `LOG_MAX_AGE` |
| `compress` | bool | true | 是否压缩日志文件 | `LOG_COMPRESS` |

**日志级别**: debug, info, warn, error

**示例**:
```yaml
log:
  level: "debug"
  format: "console"
  output: "/var/log/k8svision/app.log"
  maxSize: 200
  maxBackups: 5
  maxAge: 30
  compress: true
```

### 认证配置 (auth)

| 参数 | 类型 | 默认值 | 说明 | 环境变量 |
|------|------|--------|------|----------|
| `username` | string | "admin" | 默认用户名 | `LOGIN_USERNAME` |
| `password` | string | "admin" | 默认密码 | `LOGIN_PASSWORD` |
| `maxLoginFail` | int | 5 | 最大登录失败次数 | `LOGIN_MAX_FAIL` |
| `lockDuration` | duration | 10m | 账号锁定时间 | `LOGIN_LOCK_MINUTES` |
| `sessionTimeout` | duration | 24h | 会话超时时间 | `AUTH_SESSION_TIMEOUT` |
| `enableRateLimit` | bool | true | 是否启用频率限制 | `AUTH_ENABLE_RATE_LIMIT` |
| `rateLimit` | int | 100 | 频率限制阈值 | `AUTH_RATE_LIMIT` |

**示例**:
```yaml
auth:
  username: "admin"
  password: "secure-password-123"
  maxLoginFail: 3
  lockDuration: 15m
  sessionTimeout: 12h
  enableRateLimit: true
  rateLimit: 200
```

### 缓存配置 (cache)

| 参数 | 类型 | 默认值 | 说明 | 环境变量 |
|------|------|--------|------|----------|
| `enabled` | bool | true | 是否启用缓存 | `CACHE_ENABLED` |
| `type` | string | "memory" | 缓存类型 | `CACHE_TYPE` |
| `ttl` | duration | 5m | 缓存过期时间 | `CACHE_TTL` |
| `maxSize` | int | 1000 | 缓存最大条目数 | `CACHE_MAX_SIZE` |
| `cleanupInterval` | duration | 10m | 清理间隔 | `CACHE_CLEANUP_INTERVAL` |

**缓存类型**: memory, redis (计划支持)

**示例**:
```yaml
cache:
  enabled: true
  type: "memory"
  ttl: 10m
  maxSize: 2000
  cleanupInterval: 15m
```

## 🌍 环境变量

### 环境变量命名规则
- 前缀: `K8SVISION_`
- 分隔符: 下划线 `_`
- 大小写: 大写

### 常用环境变量

#### 基础配置
```bash
# 服务器配置
export K8SVISION_SERVER_PORT=9090
export K8SVISION_SERVER_HOST=0.0.0.0

# Kubernetes 配置
export K8SVISION_KUBERNETES_KUBECONFIG=/path/to/kubeconfig
export K8SVISION_KUBERNETES_CONTEXT=production
export K8SVISION_KUBERNETES_TIMEOUT=60s

# JWT 配置
export K8SVISION_JWT_SECRET=your-secret-key
export K8SVISION_JWT_EXPIRATION=48h

# 日志配置
export K8SVISION_LOG_LEVEL=debug
export K8SVISION_LOG_FORMAT=console

# 认证配置
export K8SVISION_AUTH_USERNAME=admin
export K8SVISION_AUTH_PASSWORD=secure-password
export K8SVISION_AUTH_MAX_LOGIN_FAIL=3

# 缓存配置
export K8SVISION_CACHE_ENABLED=true
export K8SVISION_CACHE_TTL=10m
```

#### 特殊环境变量
```bash
# 启用 Swagger 文档
export SWAGGER_ENABLE=true

# 设置 Gin 模式
export GIN_MODE=release

# 设置时区
export TZ=Asia/Shanghai
```

## 🔄 配置热重载

### 启用配置监听
```yaml
# 在配置文件中启用监听
config:
  watch: true
  watchInterval: 30s
```

### 配置变更处理
- 配置文件变更后自动重新加载
- 部分配置支持热重载
- 需要重启的配置会记录日志

## 🔒 安全配置

### 生产环境建议
```yaml
# 服务器配置
server:
  port: "8080"
  host: "0.0.0.0"
  readTimeout: 30s
  writeTimeout: 30s

# JWT 配置
jwt:
  secret: "your-very-long-and-secure-secret-key"
  expiration: 24h
  issuer: "k8svision-prod"
  audience: "k8svision-web"

# 认证配置
auth:
  username: "admin"
  password: "complex-password-with-special-chars"
  maxLoginFail: 3
  lockDuration: 15m
  enableRateLimit: true
  rateLimit: 100

# 日志配置
log:
  level: "info"
  format: "json"
  output: "/var/log/k8svision/app.log"
  maxSize: 100
  maxBackups: 5
  maxAge: 30
  compress: true

# Kubernetes 配置
kubernetes:
  kubeconfig: "/path/to/secure/kubeconfig"
  context: "production"
  timeout: 60s
  qps: 100
  burst: 200
  insecure: false
  caFile: "/path/to/ca.crt"
  certFile: "/path/to/client.crt"
  keyFile: "/path/to/client.key"
```

### 安全最佳实践
1. **使用强密码**: 密码长度至少 12 位，包含大小写字母、数字和特殊字符
2. **定期更换密钥**: JWT 密钥应定期更换
3. **限制访问**: 使用防火墙限制访问端口
4. **启用 HTTPS**: 在生产环境中使用 HTTPS
5. **日志审计**: 启用详细的日志记录
6. **权限最小化**: 使用最小权限原则配置 Kubernetes 访问

## 📊 性能调优

### 高并发配置
```yaml
# 服务器配置
server:
  readTimeout: 60s
  writeTimeout: 60s
  idleTimeout: 120s
  maxHeaderBytes: 2097152

# Kubernetes 配置
kubernetes:
  timeout: 120s
  qps: 500
  burst: 1000

# 缓存配置
cache:
  enabled: true
  ttl: 10m
  maxSize: 5000
  cleanupInterval: 15m

# 认证配置
auth:
  enableRateLimit: true
  rateLimit: 1000
```

### 内存优化
```yaml
# 缓存配置
cache:
  maxSize: 1000  # 减少缓存大小
  ttl: 5m        # 减少缓存时间

# 日志配置
log:
  maxSize: 50    # 减少日志文件大小
  maxBackups: 3  # 减少日志文件数量
```

## 🔍 配置验证

### 配置检查
```bash
# 验证配置文件语法
./k8svision --config config.yaml --validate

# 检查配置加载
./k8svision --config config.yaml --dry-run
```

### 配置测试
```bash
# 测试 Kubernetes 连接
kubectl cluster-info

# 测试配置参数
curl -X GET http://localhost:8080/healthz
```

## 📝 配置示例

### 开发环境配置
```yaml
server:
  port: "8080"
  host: "0.0.0.0"

kubernetes:
  kubeconfig: "~/.kube/config"
  context: "minikube"
  insecure: true

jwt:
  secret: "dev-secret-key"
  expiration: 24h

log:
  level: "debug"
  format: "console"

auth:
  username: "admin"
  password: "123456"
  maxLoginFail: 10

cache:
  enabled: true
  ttl: 1m
```

### 生产环境配置
```yaml
server:
  port: "8080"
  host: "0.0.0.0"
  readTimeout: 30s
  writeTimeout: 30s

kubernetes:
  kubeconfig: "/etc/k8svision/kubeconfig"
  context: "production"
  timeout: 60s
  qps: 200
  burst: 400
  insecure: false
  caFile: "/etc/k8svision/ca.crt"
  certFile: "/etc/k8svision/client.crt"
  keyFile: "/etc/k8svision/client.key"

jwt:
  secret: "production-secret-key-very-long-and-secure"
  expiration: 24h
  issuer: "k8svision-prod"
  audience: "k8svision-web"

log:
  level: "info"
  format: "json"
  output: "/var/log/k8svision/app.log"
  maxSize: 100
  maxBackups: 5
  maxAge: 30
  compress: true

auth:
  username: "admin"
  password: "complex-production-password-123"
  maxLoginFail: 3
  lockDuration: 15m
  enableRateLimit: true
  rateLimit: 100

cache:
  enabled: true
  type: "memory"
  ttl: 10m
  maxSize: 2000
  cleanupInterval: 15m
```

## 📞 获取帮助

如果配置过程中遇到问题，请：

1. 查看 [常见问题](../troubleshooting/faq.md)
2. 检查 [错误代码](../troubleshooting/error-codes.md)
3. 提交 [GitHub Issue](https://github.com/nick0323/K8sVision/issues)

## 📚 相关文档

- [快速安装](../quickstart.md)
- [项目架构](../development/architecture.md)
- [API 文档](../api/README.md)
- [部署指南](../deployment/README.md)

---

**配置版本**: v1.0.0  
**最后更新**: 2024年12月 