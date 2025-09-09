# K8sVision 快速开始指南

## 🚀 5分钟快速部署

### 方式一：Docker Compose（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/your-org/K8sVision.git
cd K8sVision

# 2. 配置环境变量
export JWT_SECRET="your-jwt-secret-key"
export AUTH_USERNAME="admin"
export AUTH_PASSWORD="admin123"

# 3. 启动服务
docker-compose up -d

# 4. 访问应用
open http://localhost:8080
```

### 方式二：Kubernetes

```bash
# 1. 创建命名空间
kubectl create namespace k8svision

# 2. 部署应用
kubectl apply -f k8s/ -n k8svision

# 3. 查看状态
kubectl get pods -n k8svision

# 4. 访问应用
kubectl port-forward svc/k8svision-frontend 8080:80 -n k8svision
```

### 方式三：源码运行

```bash
# 1. 安装依赖
go mod download
cd frontend && npm install && cd ..

# 2. 构建前端
cd frontend && npm run build && cd ..

# 3. 运行后端
go run main.go

# 4. 访问应用
open http://localhost:8080
```

## 🔧 配置说明

### 基本配置

创建 `config.yaml` 文件：

```yaml
server:
  port: "8080"
  host: "0.0.0.0"

jwt:
  secret: "your-jwt-secret-key"
  expires: 24h

auth:
  username: "admin"
  password: "$2a$10$your-hashed-password"

kubernetes:
  kubeconfig: ""  # 使用默认 kubeconfig
```

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `K8SVISION_CONFIG` | `config.yaml` | 配置文件路径 |
| `K8SVISION_PORT` | `8080` | 服务端口 |
| `K8SVISION_JWT_SECRET` | - | JWT 密钥 |
| `K8SVISION_AUTH_USERNAME` | `admin` | 认证用户名 |
| `K8SVISION_AUTH_PASSWORD` | - | 认证密码 |

## 🔐 安全配置

### 生成安全密码

```bash
# 使用内置工具生成密码
curl -X POST http://localhost:8080/admin/password/generate \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{"length": 16}'
```

### 配置 HTTPS

```yaml
# 在 config.yaml 中添加
server:
  tls:
    enabled: true
    certFile: "/path/to/cert.pem"
    keyFile: "/path/to/key.pem"
```

## 📊 监控和健康检查

### 健康检查

```bash
# 检查服务状态
curl http://localhost:8080/health

# 查看指标
curl http://localhost:8080/metrics

# 查看缓存状态
curl http://localhost:8080/cache/stats
```

### 日志查看

```bash
# Docker 环境
docker logs -f k8svision

# Kubernetes 环境
kubectl logs -f deployment/k8svision-backend -n k8svision
```

## 🚨 故障排除

### 常见问题

1. **无法连接 Kubernetes**
   ```bash
   # 检查 kubeconfig
   kubectl config current-context
   
   # 测试连接
   kubectl get nodes
   ```

2. **认证失败**
   ```bash
   # 检查密码
   curl -X POST http://localhost:8080/admin/password/verify \
     -H "Content-Type: application/json" \
     -d '{"password":"your-password","hashedPassword":"$2a$10$..."}'
   ```

3. **端口被占用**
   ```bash
   # 检查端口
   netstat -tlnp | grep 8080
   
   # 修改端口
   export K8SVISION_PORT=8081
   ```

### 调试模式

```bash
# 启用调试日志
export K8SVISION_LOG_LEVEL=debug

# 运行应用
go run main.go
```

## 📈 性能优化

### 生产环境配置

```yaml
# 生产环境推荐配置
server:
  port: "8080"
  readTimeout: 60s
  writeTimeout: 60s
  idleTimeout: 300s

cache:
  type: "memory"
  ttl: 600s
  maxSize: 5000

kubernetes:
  qps: 100
  burst: 200
  timeout: 60s
```

### 资源限制

```yaml
# Kubernetes 资源限制
resources:
  requests:
    memory: "512Mi"
    cpu: "250m"
  limits:
    memory: "1Gi"
    cpu: "500m"
```

## 🔄 升级和回滚

### 升级应用

```bash
# 使用部署脚本
./scripts/deploy.sh prod upgrade v2.0.0

# 或手动升级
kubectl set image deployment/k8svision-backend k8svision=k8svision:v2.0.0 -n k8svision
```

### 回滚应用

```bash
# 使用部署脚本
./scripts/deploy.sh prod rollback

# 或手动回滚
kubectl rollout undo deployment/k8svision-backend -n k8svision
```

## 📚 更多信息

- [完整部署文档](DEPLOYMENT.md)
- [API 文档](http://localhost:8080/swagger/index.html)
- [项目文档](https://github.com/your-org/K8sVision/docs)
- [问题报告](https://github.com/your-org/K8sVision/issues)

## 🆘 获取帮助

- **文档**: [项目文档](https://github.com/your-org/K8sVision/docs)
- **问题**: [GitHub Issues](https://github.com/your-org/K8sVision/issues)
- **讨论**: [GitHub Discussions](https://github.com/your-org/K8sVision/discussions)
- **邮件**: support@k8svision.com

---

**注意**: 请根据您的实际环境调整配置参数。在生产环境中部署前，请务必进行充分的测试。
