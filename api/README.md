# API模块文档

API模块是K8sVision的后端接口层，负责处理HTTP请求、认证授权、数据验证和响应格式化。

## 📁 模块结构

```
api/
├── README.md                    # 模块文档
├── common.go                    # 通用API工具
├── auth_manager.go              # 认证管理器
├── login.go                     # 登录接口
├── password_management.go       # 密码管理
├── metrics.go                   # 监控指标接口
├── overview.go                  # 集群概览接口
├── namespace.go                 # 命名空间接口
├── node.go                      # 节点接口
├── pod.go                       # Pod接口
├── deployment.go                # Deployment接口
├── statefulset.go               # StatefulSet接口
├── daemonset.go                 # DaemonSet接口
├── service.go                   # Service接口
├── ingress.go                   # Ingress接口
├── job.go                       # Job接口
├── cronjob.go                   # CronJob接口
├── event.go                     # Event接口
├── pvc.go                       # PVC接口
├── pv.go                        # PV接口
├── storageclass.go              # StorageClass接口
├── configmap.go                 # ConfigMap接口
├── secret.go                    # Secret接口
└── middleware/                  # 中间件
    ├── auth.go                  # JWT认证中间件
    ├── error.go                  # 错误处理中间件
    ├── logging.go               # 日志中间件
    └── metrics.go               # 监控中间件
```

## 🔧 核心组件

### 1. 通用工具 (common.go)
提供API层的通用功能：
- 分页处理
- 搜索过滤
- 响应格式化
- 错误处理

**主要函数：**
- `HandleListWithPagination`: 处理带分页的列表请求
- `HandleDetailWithK8s`: 处理资源详情请求
- `ResponseSuccess`: 成功响应格式化
- `ResponseError`: 错误响应格式化

### 2. 认证管理 (auth_manager.go)
负责用户认证和会话管理：
- 登录尝试记录
- 密码验证
- 会话管理
- 登录限制

**主要结构：**
```go
type AuthManager struct {
    logger     *zap.Logger
    configMgr  *config.Manager
    attempts   map[string]*LoginAttempt
    mutex      sync.RWMutex
}
```

### 3. 中间件系统 (middleware/)
提供请求处理中间件：

#### 认证中间件 (auth.go)
- JWT令牌验证
- 用户身份确认
- 权限检查

#### 错误处理中间件 (error.go)
- 统一错误响应格式
- Panic恢复
- 错误日志记录

#### 日志中间件 (logging.go)
- 请求日志记录
- 响应时间统计
- 链路追踪

#### 监控中间件 (metrics.go)
- 性能指标收集
- 缓存统计
- 并发控制

## 🚀 API接口

### 认证接口

#### POST /api/login
用户登录接口
```json
{
  "username": "admin",
  "password": "password"
}
```

**响应：**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "jwt-token",
    "expires": "2024-01-01T00:00:00Z"
  }
}
```

### 资源管理接口

#### GET /api/overview
获取集群概览信息
**参数：** 无
**响应：** 集群状态、资源统计、健康指标

#### GET /api/pods
获取Pod列表
**参数：**
- `namespace` (可选): 命名空间过滤
- `limit` (可选): 每页数量
- `offset` (可选): 偏移量
- `search` (可选): 搜索关键词

#### GET /api/pods/{namespace}/{name}
获取Pod详情
**参数：**
- `namespace`: 命名空间
- `name`: Pod名称

### 监控接口

#### GET /api/metrics
获取系统指标
**响应：** 系统性能指标、业务指标

#### GET /api/metrics/health
获取健康状态
**响应：** 系统健康状态、关键指标

## 🔒 安全特性

### 认证机制
- JWT令牌认证
- 密码哈希存储
- 会话超时管理

### 安全措施
- 登录失败限制
- 速率限制
- 请求验证
- 错误信息脱敏

### 权限控制
- 基于角色的访问控制
- API接口权限验证
- 资源访问限制

## 📊 性能优化

### 缓存策略
- 内存缓存
- 缓存过期管理
- 缓存命中率统计

### 并发控制
- 请求并发限制
- 资源池管理
- 优雅降级

### 响应优化
- 分页查询
- 字段过滤
- 数据压缩

## 🛠️ 开发指南

### 添加新接口
1. 在对应的资源文件中添加处理函数
2. 实现业务逻辑
3. 添加错误处理
4. 注册路由

### 错误处理
使用统一的错误响应格式：
```go
middleware.ResponseError(c, logger, &model.APIError{
    Code:    model.CodeBadRequest,
    Message: "错误描述",
    Details: "详细错误信息",
}, http.StatusBadRequest)
```

### 日志记录
使用结构化日志：
```go
logger.Info("操作描述",
    zap.String("key", "value"),
    zap.Int("count", 123),
)
```

## 📝 最佳实践

1. **接口设计**
   - 遵循RESTful设计原则
   - 使用统一的响应格式
   - 提供清晰的错误信息

2. **性能考虑**
   - 实现分页查询
   - 使用缓存减少重复请求
   - 避免N+1查询问题

3. **安全实践**
   - 验证所有输入参数
   - 使用参数化查询
   - 记录安全相关事件

4. **可维护性**
   - 保持函数简洁
   - 添加必要的注释
   - 编写单元测试

## 🔍 故障排查

### 常见问题
1. **认证失败**
   - 检查JWT密钥配置
   - 验证用户名密码
   - 查看登录限制状态

2. **权限错误**
   - 检查Kubernetes RBAC配置
   - 验证集群连接
   - 查看错误日志

3. **性能问题**
   - 检查缓存配置
   - 监控并发数
   - 分析慢查询

### 调试工具
- 日志级别调整
- 指标监控
- 链路追踪
- 健康检查接口
