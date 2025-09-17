# K8sVision

一个现代化的Kubernetes集群可视化管理平台，提供直观的Web界面来管理和监控Kubernetes资源。

## 🚀 功能特性

### 核心功能
- **集群概览**: 实时显示集群状态、资源使用情况和关键指标
- **资源管理**: 支持Pod、Deployment、Service、Ingress等主要K8s资源的管理
- **实时监控**: 提供资源使用情况、性能指标和健康状态监控
- **多集群支持**: 支持管理多个Kubernetes集群
- **安全认证**: 基于JWT的身份认证和授权机制

### 支持的资源类型
- **工作负载**: Pod、Deployment、StatefulSet、DaemonSet、Job、CronJob
- **服务发现**: Service、Ingress
- **配置管理**: ConfigMap、Secret
- **存储管理**: PVC、PV、StorageClass
- **集群管理**: Node、Namespace、Event

### 技术特性
- **现代化UI**: 基于React的响应式用户界面
- **高性能**: 内存缓存和并发控制优化
- **可观测性**: 完整的日志记录、指标收集和链路追踪
- **安全性**: 密码加密、登录限制、速率限制
- **可扩展性**: 模块化架构，易于扩展新功能

## 🏗️ 架构设计

### 后端架构
```
K8sVision Backend
├── API Layer (api/)
│   ├── 资源管理接口
│   ├── 认证中间件
│   └── 错误处理
├── Service Layer (service/)
│   ├── Kubernetes客户端封装
│   ├── 业务逻辑处理
│   └── 数据转换
├── Model Layer (model/)
│   ├── 数据模型定义
│   ├── 常量定义
│   └── 配置结构
├── Cache Layer (cache/)
│   ├── 内存缓存管理
│   └── 缓存策略
├── Config Layer (config/)
│   └── 配置管理
└── Monitor Layer (monitor/)
    ├── 指标收集
    ├── 业务监控
    └── 链路追踪
```

### 前端架构
```
K8sVision Frontend
├── Components (components/)
│   ├── 通用组件
│   ├── 资源详情组件
│   └── 状态渲染组件
├── Pages (pages/)
│   ├── 资源列表页面
│   └── 概览页面
├── Hooks (hooks/)
│   ├── 状态管理
│   ├── 分页处理
│   └── 搜索功能
├── Utils (utils/)
│   ├── API工具
│   ├── 数据处理
│   └── 认证工具
└── Constants (constants/)
    ├── 页面配置
    └── 常量定义
```

## 🛠️ 技术栈

### 后端技术
- **语言**: Go 1.24+
- **Web框架**: Gin
- **Kubernetes客户端**: client-go
- **日志**: Zap
- **配置管理**: Viper
- **认证**: JWT
- **缓存**: 内存缓存
- **监控**: 自定义指标收集

### 前端技术
- **框架**: React 18
- **构建工具**: Vite
- **UI组件**: 自定义组件库
- **状态管理**: React Hooks
- **HTTP客户端**: Fetch API
- **图标**: React Icons

### 部署技术
- **容器化**: Docker
- **编排**: Kubernetes
- **反向代理**: Nginx
- **配置**: ConfigMap/Secret

## 📦 快速开始

### 环境要求
- Go 1.24+
- Node.js 18+
- Kubernetes集群访问权限
- Docker (可选)

### 本地开发

1. **克隆项目**
```bash
git clone <repository-url>
cd K8sVision
```

2. **配置环境变量**
```bash
export K8SVISION_JWT_SECRET="your-32-character-secret-key"
export K8SVISION_AUTH_USERNAME="admin"
export K8SVISION_AUTH_PASSWORD="your-password"
```

3. **启动后端服务**
```bash
# 安装依赖
go mod download

# 运行服务
go run main.go
```

4. **启动前端服务**
```bash
cd frontend
npm install
npm run dev
```

5. **访问应用**
- 前端: http://localhost:5173
- 后端API: http://localhost:8080

### Docker部署

1. **构建镜像**
```bash
# 构建后端镜像
docker build -t k8svision-backend .

# 构建前端镜像
cd frontend
docker build -t k8svision-frontend .
```

2. **使用Docker Compose**
```bash
docker-compose up -d
```

### Kubernetes部署

1. **创建Kubernetes资源**
```bash
# 创建命名空间
kubectl create namespace k8svision

# 创建ConfigMap和Secret
kubectl create configmap k8svision-config --from-file=config.yaml -n k8svision
kubectl create secret generic k8svision-secrets \
  --from-literal=jwt-secret=your-jwt-secret \
  --from-literal=auth-username=admin \
  --from-literal=auth-password=your-password \
  -n k8svision

# 部署应用
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8svision-backend
  namespace: k8svision
spec:
  replicas: 1
  selector:
    matchLabels:
      app: k8svision-backend
  template:
    metadata:
      labels:
        app: k8svision-backend
    spec:
      containers:
      - name: k8svision-backend
        image: k8svision:latest
        ports:
        - containerPort: 8080
        env:
        - name: K8SVISION_JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: k8svision-secrets
              key: jwt-secret
        - name: K8SVISION_AUTH_USERNAME
          valueFrom:
            secretKeyRef:
              name: k8svision-secrets
              key: auth-username
        - name: K8SVISION_AUTH_PASSWORD
          valueFrom:
            secretKeyRef:
              name: k8svision-secrets
              key: auth-password
EOF
```

2. **检查部署状态**
```bash
kubectl get pods -n k8svision
```

## ⚙️ 配置说明

### 后端配置 (config.yaml)
```yaml
server:
  port: "8080"
  host: "0.0.0.0"

kubernetes:
  kubeconfig: ""  # Kubeconfig文件路径
  context: ""     # 集群上下文

jwt:
  secret: ""      # JWT密钥
  expiration: "24h"

auth:
  username: ""    # 管理员用户名
  password: ""    # 管理员密码

cache:
  enabled: true
  ttl: "5m"
```

### 环境变量
- `K8SVISION_JWT_SECRET`: JWT密钥 (必需)
- `K8SVISION_AUTH_USERNAME`: 管理员用户名 (必需)
- `K8SVISION_AUTH_PASSWORD`: 管理员密码 (必需)
- `K8SVISION_KUBECONFIG`: Kubeconfig文件路径 (可选)
- `K8SVISION_LOG_LEVEL`: 日志级别 (可选)

## 📚 API文档

### 认证接口
- `POST /api/login` - 用户登录

### 资源管理接口
- `GET /api/overview` - 集群概览
- `GET /api/pods` - Pod列表
- `GET /api/deployments` - Deployment列表
- `GET /api/services` - Service列表
- `GET /api/nodes` - Node列表
- `GET /api/namespaces` - Namespace列表

### 监控接口
- `GET /api/metrics` - 系统指标
- `GET /api/metrics/health` - 健康检查

## 🔧 开发指南

### 代码结构
```
K8sVision/
├── api/                 # API接口层
├── service/             # 业务逻辑层
├── model/               # 数据模型
├── cache/               # 缓存管理
├── config/              # 配置管理
├── monitor/             # 监控模块
├── frontend/            # 前端代码
├── main.go              # 主程序入口
└── config.yaml          # 配置文件
```

### 添加新资源类型
1. 在`model/types.go`中定义资源状态结构
2. 在`service/`中实现资源获取逻辑
3. 在`api/`中实现API接口
4. 在`main.go`中注册路由
5. 在前端添加对应的页面组件

### 代码规范
- 遵循Go语言官方代码规范
- 使用有意义的变量和函数命名
- 添加必要的注释和文档
- 编写单元测试
- 使用`go fmt`格式化代码

## 🚀 部署指南

### 生产环境部署
1. 配置环境变量
2. 设置Kubernetes集群访问权限
3. 配置反向代理和SSL证书
4. 设置监控和日志收集
5. 配置备份和恢复策略

### 安全建议
- 使用强密码和JWT密钥
- 启用HTTPS
- 配置网络策略
- 定期更新依赖
- 监控异常访问

## 🤝 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 📄 许可证

本项目采用MIT许可证 - 查看[LICENSE](LICENSE)文件了解详情。

## 🆘 支持

如有问题或建议，请：
1. 查看[Issues](../../issues)
2. 创建新的Issue
3. 联系维护团队

---

**K8sVision** - 让Kubernetes管理更简单、更直观！
