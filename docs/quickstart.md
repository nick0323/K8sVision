# 快速安装指南

本指南将帮助您在 5 分钟内完成 K8sVision 的安装和部署。

## 🎯 安装方式选择

| 安装方式 | 适用场景 | 复杂度 | 时间 |
|---------|---------|--------|------|
| Docker Compose | 本地开发、测试 | ⭐ | 5分钟 |
| Kubernetes | 生产环境 | ⭐⭐⭐ | 15分钟 |
| 源码编译 | 开发调试 | ⭐⭐ | 10分钟 |

## 🐳 Docker Compose 安装（推荐）

### 前置要求
- Docker 20.10+
- Docker Compose 2.0+
- 至少 2GB 可用内存

### 安装步骤

#### 1. 克隆项目
```bash
git clone https://github.com/nick0323/K8sVision.git
cd K8sVision
```

#### 2. 配置 Kubernetes 访问
确保您有可用的 Kubernetes 集群访问权限：

```bash
# 检查 kubectl 配置
kubectl cluster-info

# 或者使用 kubeconfig 文件
export KUBECONFIG=/path/to/your/kubeconfig
```

#### 3. 修改配置（可选）
编辑 `docker-compose.yml` 文件，根据需要调整配置：

```yaml
services:
  backend:
    environment:
      - LOGIN_USERNAME=admin          # 登录用户名
      - LOGIN_PASSWORD=12345678       # 登录密码
      - KUBECONFIG=/app/config.yaml   # kubeconfig 路径
    volumes:
      - ~/.kube/config:/app/config.yaml:ro  # 挂载 kubeconfig
```

#### 4. 启动服务
```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

#### 5. 访问应用
- **前端界面**: http://localhost
- **后端 API**: http://localhost:8080/api
- **Swagger 文档**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/healthz

#### 6. 登录系统
- 用户名: `admin`
- 密码: `12345678`（或您在配置中设置的密码）

## ☸️ Kubernetes 安装

### 前置要求
- Kubernetes 1.20+
- kubectl 已配置
- Helm 3.0+（可选）

### 安装步骤

#### 1. 克隆项目
```bash
git clone https://github.com/nick0323/K8sVision.git
cd K8sVision
```

#### 2. 创建命名空间
```bash
kubectl create namespace k8svision
```

#### 3. 配置 Secret
```bash
# 创建包含 kubeconfig 的 Secret
kubectl create secret generic k8svision-kubeconfig \
  --from-file=config=~/.kube/config \
  -n k8svision
```

#### 4. 修改部署配置
编辑 `k8s/deployment.yaml` 文件，根据需要调整配置。

#### 5. 部署应用
```bash
# 部署所有资源
kubectl apply -f k8s/

# 查看部署状态
kubectl get all -n k8svision

# 查看 Pod 日志
kubectl logs -f deployment/k8svision-backend -n k8svision
```

#### 6. 配置 Ingress（可选）
```bash
# 创建 Ingress
kubectl apply -f k8s/ingress.yaml

# 添加域名解析
echo "127.0.0.1 k8svision.local" >> /etc/hosts
```

#### 7. 访问应用
- **通过 Ingress**: https://k8svision.local
- **通过 NodePort**: http://<node-ip>:30080
- **通过 Port-Forward**: 
  ```bash
  kubectl port-forward svc/k8svision-frontend 80:80 -n k8svision
  ```

## 🔧 源码编译安装

### 前置要求
- Go 1.24+
- Node.js 18+
- npm 或 yarn

### 安装步骤

#### 1. 克隆项目
```bash
git clone https://github.com/nick0323/K8sVision.git
cd K8sVision
```

#### 2. 编译后端
```bash
# 安装依赖
go mod tidy

# 编译
go build -o k8svision main.go

# 运行
./k8svision
```

#### 3. 编译前端
```bash
cd frontend

# 安装依赖
npm install

# 开发模式运行
npm run dev

# 或构建生产版本
npm run build
```

#### 4. 配置环境变量
```bash
export LOGIN_USERNAME=admin
export LOGIN_PASSWORD=12345678
export KUBECONFIG=/path/to/your/kubeconfig
export SWAGGER_ENABLE=true
```

## ⚙️ 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 | 必需 |
|--------|------|--------|------|
| `LOGIN_USERNAME` | 登录用户名 | admin | 否 |
| `LOGIN_PASSWORD` | 登录密码 | 123456 | 否 |
| `JWT_SECRET` | JWT 密钥 | k8svision-secret-key | 否 |
| `KUBECONFIG` | kubeconfig 路径 | "" | 是 |
| `SWAGGER_ENABLE` | 启用 Swagger | false | 否 |
| `LOGIN_MAX_FAIL` | 最大登录失败次数 | 5 | 否 |
| `LOGIN_LOCK_MINUTES` | 锁定时间（分钟） | 10 | 否 |

### 配置文件
项目支持通过 `config.yaml` 文件进行配置，详细配置说明请参考 [配置文档](./configuration.md)。

## 🔍 验证安装

### 1. 检查服务状态
```bash
# Docker Compose
docker-compose ps

# Kubernetes
kubectl get pods -n k8svision
```

### 2. 检查健康状态
```bash
# 健康检查
curl http://localhost:8080/healthz

# 预期输出: ok
```

### 3. 检查 API 状态
```bash
# 获取集群概览
curl -H "Authorization: Bearer <your-token>" \
  http://localhost:8080/api/overview
```

### 4. 检查前端访问
在浏览器中访问 http://localhost，应该能看到登录界面。

## 🚨 常见问题

### 1. 无法连接 Kubernetes 集群
**问题**: 后端无法连接到 Kubernetes 集群
**解决方案**:
```bash
# 检查 kubeconfig 配置
kubectl cluster-info

# 确保 kubeconfig 文件权限正确
chmod 600 ~/.kube/config

# 检查集群连接
kubectl get nodes
```

### 2. 前端无法访问后端 API
**问题**: 前端显示 API 连接错误
**解决方案**:
```bash
# 检查后端服务状态
docker-compose logs backend

# 检查网络连接
curl http://localhost:8080/healthz

# 检查 CORS 配置
```

### 3. 登录失败
**问题**: 无法登录系统
**解决方案**:
```bash
# 检查用户名密码配置
echo $LOGIN_USERNAME
echo $LOGIN_PASSWORD

# 查看登录日志
docker-compose logs backend | grep login

# 重置登录失败计数
docker-compose restart backend
```

### 4. 端口冲突
**问题**: 端口被占用
**解决方案**:
```bash
# 查看端口占用
netstat -tlnp | grep :8080

# 修改端口配置
# 编辑 docker-compose.yml 或 k8s 配置文件
```

## 📞 获取帮助

如果遇到其他问题，请：

1. 查看 [常见问题](./troubleshooting/faq.md)
2. 检查 [错误代码](./troubleshooting/error-codes.md)
3. 查看 [日志分析](./troubleshooting/logs.md)
4. 提交 [GitHub Issue](https://github.com/nick0323/K8sVision/issues)

## 🎉 下一步

安装完成后，您可以：

1. 阅读 [用户手册](./user-guide/README.md) 了解如何使用
2. 查看 [API 文档](./api/README.md) 了解接口
3. 参考 [开发指南](./development/setup.md) 进行二次开发
4. 查看 [监控运维](./deployment/monitoring.md) 配置监控

---

**恭喜！您已成功安装 K8sVision！** 🎊 