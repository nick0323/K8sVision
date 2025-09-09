# K8sVision 部署文档

## 📋 目录

- [项目概述](#项目概述)
- [系统要求](#系统要求)
- [快速开始](#快速开始)
- [详细部署步骤](#详细部署步骤)
- [配置说明](#配置说明)
- [监控和运维](#监控和运维)
- [故障排除](#故障排除)
- [安全建议](#安全建议)
- [性能优化](#性能优化)

## 🚀 项目概述

K8sVision 是一个基于 Web 的 Kubernetes 集群管理平台，提供直观的界面来查看和管理 Kubernetes 资源。项目采用前后端分离架构：

- **后端**: Go + Gin + Kubernetes client-go
- **前端**: React + Vite + Ant Design
- **数据库**: 无状态设计，直接连接 Kubernetes API
- **部署**: Docker + Kubernetes

## 💻 系统要求

### 最低要求
- **CPU**: 2 核
- **内存**: 4GB RAM
- **存储**: 10GB 可用空间
- **网络**: 能够访问 Kubernetes 集群

### 推荐配置
- **CPU**: 4 核
- **内存**: 8GB RAM
- **存储**: 20GB 可用空间
- **网络**: 低延迟网络连接

### 软件依赖
- **Docker**: 20.10+ 
- **Kubernetes**: 1.20+
- **kubectl**: 1.20+
- **Go**: 1.24+ (仅开发环境)
- **Node.js**: 18+ (仅开发环境)

## 🚀 快速开始

### 1. 克隆项目
```bash
git clone https://github.com/your-org/K8sVision.git
cd K8sVision
```

### 2. 使用 Docker Compose 快速部署
```bash
# 复制配置文件
cp config.yaml.example config.yaml

# 编辑配置文件
vim config.yaml

# 启动服务
docker-compose up -d

# 查看服务状态
docker-compose ps
```

### 3. 访问应用
- **Web界面**: http://localhost:8080
- **API文档**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/health

## 📋 详细部署步骤

### 方式一：Docker 部署

#### 1. 构建镜像
```bash
# 构建完整镜像
docker build -t k8svision:latest .

# 或者分别构建前后端
docker build -f frontend/Dockerfile -t k8svision-frontend:latest ./frontend
docker build -f Dockerfile -t k8svision-backend:latest .
```

#### 2. 配置环境
```bash
# 创建配置目录
mkdir -p /opt/k8svision/config

# 复制配置文件
cp config.yaml /opt/k8svision/config/

# 设置权限
chmod 600 /opt/k8svision/config/config.yaml
```

#### 3. 运行容器
```bash
# 运行后端服务
docker run -d \
  --name k8svision-backend \
  -p 8080:8080 \
  -v /opt/k8svision/config:/app/config \
  -v ~/.kube:/root/.kube:ro \
  k8svision:latest

# 查看日志
docker logs -f k8svision-backend
```

### 方式二：Kubernetes 部署

#### 1. 准备 Kubernetes 配置
```bash
# 创建命名空间
kubectl create namespace k8svision

# 创建 ConfigMap
kubectl create configmap k8svision-config \
  --from-file=config.yaml \
  -n k8svision

# 创建 Secret (用于认证)
kubectl create secret generic k8svision-auth \
  --from-literal=username=admin \
  --from-literal=password='your-secure-password' \
  -n k8svision
```

#### 2. 部署应用
```bash
# 部署后端
kubectl apply -f k8s/backend-deployment.yaml

# 部署前端
kubectl apply -f k8s/frontend-deployment.yaml

# 部署 Ingress
kubectl apply -f k8s/ingress.yaml

# 查看部署状态
kubectl get pods -n k8svision
kubectl get svc -n k8svision
```

#### 3. 配置 RBAC
```yaml
# rbac.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: k8svision
  namespace: k8svision
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: k8svision
rules:
- apiGroups: [""]
  resources: ["*"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["*"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["extensions"]
  resources: ["*"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: k8svision
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: k8svision
subjects:
- kind: ServiceAccount
  name: k8svision
  namespace: k8svision
```

### 方式三：源码部署

#### 1. 后端部署
```bash
# 安装依赖
go mod download

# 构建应用
go build -o k8svision main.go

# 运行应用
./k8svision -config=config.yaml
```

#### 2. 前端部署
```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 构建前端
npm run build

# 使用 nginx 部署
sudo cp -r dist/* /var/www/html/
```

## ⚙️ 配置说明

### 主配置文件 (config.yaml)

```yaml
# 服务器配置
server:
  port: "8080"
  host: "0.0.0.0"
  readTimeout: 30s
  writeTimeout: 30s
  idleTimeout: 120s

# JWT 配置
jwt:
  secret: "your-jwt-secret-key"
  expires: 24h

# 认证配置
auth:
  username: "admin"
  password: "$2a$10$your-hashed-password"  # 使用 bcrypt 哈希
  maxLoginFail: 5
  lockDuration: 15m

# 日志配置
log:
  level: "info"
  format: "json"
  output: "stdout"

# 缓存配置
cache:
  type: "memory"
  ttl: 300s
  maxSize: 1000

# Kubernetes 配置
kubernetes:
  kubeconfig: ""  # 空值使用默认 kubeconfig
  qps: 50
  burst: 100
  timeout: 30s
```

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `K8SVISION_CONFIG` | `config.yaml` | 配置文件路径 |
| `K8SVISION_PORT` | `8080` | 服务端口 |
| `K8SVISION_LOG_LEVEL` | `info` | 日志级别 |
| `K8SVISION_JWT_SECRET` | - | JWT 密钥 |
| `K8SVISION_AUTH_USERNAME` | `admin` | 认证用户名 |
| `K8SVISION_AUTH_PASSWORD` | - | 认证密码 |
| `K8SVISION_KUBECONFIG` | - | Kubeconfig 路径 |

### 生产环境配置建议

```yaml
# 生产环境配置示例
server:
  port: "8080"
  host: "0.0.0.0"
  readTimeout: 60s
  writeTimeout: 60s
  idleTimeout: 300s

jwt:
  secret: "your-super-secure-jwt-secret-key-here"
  expires: 8h

auth:
  username: "admin"
  password: "$2a$10$your-production-hashed-password"
  maxLoginFail: 3
  lockDuration: 30m

log:
  level: "warn"
  format: "json"
  output: "stdout"

cache:
  type: "memory"
  ttl: 600s
  maxSize: 5000

kubernetes:
  kubeconfig: ""
  qps: 100
  burst: 200
  timeout: 60s
```

## 📊 监控和运维

### 健康检查

```bash
# 检查服务健康状态
curl http://localhost:8080/health

# 检查指标
curl http://localhost:8080/metrics

# 检查缓存状态
curl http://localhost:8080/cache/stats
```

### 日志管理

```bash
# 查看应用日志
docker logs -f k8svision-backend

# 在 Kubernetes 中查看日志
kubectl logs -f deployment/k8svision-backend -n k8svision

# 查看特定级别的日志
kubectl logs deployment/k8svision-backend -n k8svision | grep ERROR
```

### 性能监控

```bash
# 查看系统指标
curl http://localhost:8080/metrics/system

# 查看业务指标
curl http://localhost:8080/metrics/business

# 查看健康指标
curl http://localhost:8080/metrics/health
```

### 备份和恢复

```bash
# 备份配置
kubectl get configmap k8svision-config -n k8svision -o yaml > backup-config.yaml

# 备份 Secret
kubectl get secret k8svision-auth -n k8svision -o yaml > backup-secret.yaml

# 恢复配置
kubectl apply -f backup-config.yaml
kubectl apply -f backup-secret.yaml
```

## 🔧 故障排除

### 常见问题

#### 1. 服务无法启动
```bash
# 检查端口占用
netstat -tlnp | grep 8080

# 检查配置文件
./k8svision -config=config.yaml -check-config

# 查看详细日志
./k8svision -config=config.yaml -log-level=debug
```

#### 2. 无法连接 Kubernetes
```bash
# 检查 kubeconfig
kubectl config current-context

# 测试连接
kubectl get nodes

# 检查权限
kubectl auth can-i get pods --all-namespaces
```

#### 3. 认证失败
```bash
# 检查密码哈希
curl -X POST http://localhost:8080/admin/password/verify \
  -H "Content-Type: application/json" \
  -d '{"password":"your-password","hashedPassword":"$2a$10$..."}'

# 重置密码
curl -X POST http://localhost:8080/admin/password/change \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{"oldPassword":"old","newPassword":"new"}'
```

#### 4. 内存泄漏
```bash
# 查看内存使用
docker stats k8svision-backend

# 查看 Go 运行时信息
curl http://localhost:8080/debug/pprof/

# 生成内存分析文件
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

### 日志分析

```bash
# 查看错误日志
grep "ERROR" /var/log/k8svision.log

# 查看访问日志
grep "GET\|POST" /var/log/k8svision.log | tail -100

# 分析性能问题
grep "slow" /var/log/k8svision.log
```

## 🔒 安全建议

### 1. 认证安全
- 使用强密码策略
- 定期更换 JWT 密钥
- 启用登录失败锁定
- 使用 HTTPS 传输

### 2. 网络安全
```yaml
# 网络策略示例
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: k8svision-netpol
  namespace: k8svision
spec:
  podSelector:
    matchLabels:
      app: k8svision
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
```

### 3. 资源限制
```yaml
# 资源限制示例
resources:
  requests:
    memory: "512Mi"
    cpu: "250m"
  limits:
    memory: "1Gi"
    cpu: "500m"
```

### 4. 安全扫描
```bash
# 扫描镜像漏洞
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  aquasec/trivy image k8svision:latest

# 扫描代码漏洞
gosec ./...
```

## ⚡ 性能优化

### 1. 缓存优化
```yaml
# 调整缓存配置
cache:
  type: "memory"
  ttl: 600s
  maxSize: 10000
  cleanupInterval: 300s
```

### 2. 并发优化
```yaml
# 调整并发配置
kubernetes:
  qps: 100
  burst: 200
  timeout: 60s
```

### 3. 内存优化
```bash
# 设置 Go 运行时参数
export GOGC=100
export GOMEMLIMIT=1GiB

# 启用内存分析
export GODEBUG=madvdontneed=1
```

### 4. 网络优化
```yaml
# 调整网络超时
server:
  readTimeout: 30s
  writeTimeout: 30s
  idleTimeout: 120s
```

## 📈 扩展和升级

### 水平扩展
```bash
# 扩展副本数
kubectl scale deployment k8svision-backend --replicas=3 -n k8svision

# 使用 HPA
kubectl apply -f hpa.yaml
```

### 版本升级
```bash
# 滚动更新
kubectl set image deployment/k8svision-backend \
  k8svision=k8svision:v2.0.0 -n k8svision

# 回滚
kubectl rollout undo deployment/k8svision-backend -n k8svision
```

### 数据迁移
```bash
# 导出配置
kubectl get configmap k8svision-config -n k8svision -o yaml > config-backup.yaml

# 导入配置
kubectl apply -f config-backup.yaml
```

## 📞 支持和联系

- **文档**: [项目文档](https://github.com/your-org/K8sVision/docs)
- **问题报告**: [GitHub Issues](https://github.com/your-org/K8sVision/issues)
- **讨论**: [GitHub Discussions](https://github.com/your-org/K8sVision/discussions)
- **邮件**: support@k8svision.com

## 📄 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

---

**注意**: 请根据您的实际环境调整配置参数。在生产环境中部署前，请务必进行充分的测试。
