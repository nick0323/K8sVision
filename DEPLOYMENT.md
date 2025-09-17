# 部署和运维文档

本文档介绍K8sVision的部署、配置、监控和运维相关内容。

## 📋 目录

- [环境要求](#环境要求)
- [部署方式](#部署方式)
- [配置管理](#配置管理)
- [监控和日志](#监控和日志)
- [安全配置](#安全配置)
- [故障排查](#故障排查)
- [维护指南](#维护指南)

## 🔧 环境要求

### 最低要求
- **Kubernetes**: 1.20+
- **CPU**: 2核心
- **内存**: 4GB
- **存储**: 10GB

### 推荐配置
- **Kubernetes**: 1.24+
- **CPU**: 4核心
- **内存**: 8GB
- **存储**: 50GB

### 依赖组件
- **Ingress Controller**: Nginx/Traefik
- **Storage Class**: 支持动态卷供应
- **RBAC**: 集群管理员权限

## 🚀 部署方式

### 1. 使用Kubernetes YAML文件

#### 创建命名空间
```bash
kubectl create namespace k8svision
```

#### 创建RBAC配置
```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: k8svision
  namespace: k8svision
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: k8svision-role
rules:
- apiGroups: [""]
  resources: ["pods", "services", "nodes", "namespaces", "events", "persistentvolumes", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets", "daemonsets", "replicasets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["batch"]
  resources: ["jobs", "cronjobs"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["storage.k8s.io"]
  resources: ["storageclasses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: k8svision-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: k8svision-role
subjects:
- kind: ServiceAccount
  name: k8svision
  namespace: k8svision
EOF
```

#### 部署后端服务
```bash
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
      serviceAccountName: k8svision
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
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: k8svision-backend
  namespace: k8svision
spec:
  selector:
    app: k8svision-backend
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
EOF
```

#### 部署前端服务
```bash
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8svision-frontend
  namespace: k8svision
spec:
  replicas: 1
  selector:
    matchLabels:
      app: k8svision-frontend
  template:
    metadata:
      labels:
        app: k8svision-frontend
    spec:
      containers:
      - name: k8svision-frontend
        image: k8svision-frontend:latest
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
---
apiVersion: v1
kind: Service
metadata:
  name: k8svision-frontend
  namespace: k8svision
spec:
  selector:
    app: k8svision-frontend
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
EOF
```

#### 配置Ingress
```bash
kubectl apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: k8svision-ingress
  namespace: k8svision
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    nginx.ingress.kubernetes.io/ssl-redirect: "false"
spec:
  rules:
  - host: k8svision.local
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: k8svision-frontend
            port:
              number: 80
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: k8svision-backend
            port:
              number: 8080
EOF
```

#### 应用网络策略
```bash
kubectl apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: k8svision-network-policy
  namespace: k8svision
spec:
  podSelector:
    matchLabels:
      app: k8svision-backend
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
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          name: kube-system
    ports:
    - protocol: TCP
      port: 443
EOF
```

### 2. 使用Helm Chart

#### 添加Helm仓库
```bash
helm repo add k8svision https://charts.k8svision.com
helm repo update
```

#### 安装Chart
```bash
helm install k8svision k8svision/k8svision \
  --namespace k8svision \
  --create-namespace \
  --set config.jwt.secret="your-jwt-secret" \
  --set config.auth.username="admin" \
  --set config.auth.password="your-password"
```

### 3. 使用Docker Compose

#### 创建docker-compose.yml
```yaml
version: '3.8'
services:
  k8svision-backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - K8SVISION_JWT_SECRET=your-jwt-secret
      - K8SVISION_AUTH_USERNAME=admin
      - K8SVISION_AUTH_PASSWORD=your-password
      - K8SVISION_KUBECONFIG=/root/.kube/config
    volumes:
      - ~/.kube:/root/.kube:ro
    depends_on:
      - k8svision-frontend

  k8svision-frontend:
    build: ./frontend
    ports:
      - "3000:80"
    depends_on:
      - k8svision-backend
```

#### 启动服务
```bash
docker-compose up -d
```

## ⚙️ 配置管理

### 环境变量配置

#### 必需环境变量
```bash
# JWT密钥 (至少32位字符)
K8SVISION_JWT_SECRET="your-32-character-secret-key"

# 管理员用户名
K8SVISION_AUTH_USERNAME="admin"

# 管理员密码
K8SVISION_AUTH_PASSWORD="your-secure-password"
```

#### 可选环境变量
```bash
# Kubernetes配置
K8SVISION_KUBECONFIG="/path/to/kubeconfig"
K8SVISION_K8S_CONTEXT="production"

# 服务配置
K8SVISION_PORT="8080"
K8SVISION_HOST="0.0.0.0"

# 日志配置
K8SVISION_LOG_LEVEL="info"
K8SVISION_LOG_FORMAT="json"

# 缓存配置
K8SVISION_CACHE_ENABLED="true"
K8SVISION_CACHE_TTL="5m"
```

### ConfigMap配置

#### 创建ConfigMap
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: k8svision-config
  namespace: k8svision
data:
  config.yaml: |
    server:
      port: "8080"
      host: "0.0.0.0"
    
    kubernetes:
      context: "production"
      timeout: "30s"
    
    log:
      level: "info"
      format: "json"
    
    cache:
      enabled: true
      ttl: "5m"
```

### Secret配置

#### 创建Secret
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: k8svision-secrets
  namespace: k8svision
type: Opaque
data:
  jwt-secret: <base64-encoded-secret>
  auth-username: <base64-encoded-username>
  auth-password: <base64-encoded-password>
```

## 📊 监控和日志

### 指标监控

#### 应用指标
- 请求总数和成功率
- 响应时间分布
- 错误率统计
- 并发连接数

#### 系统指标
- CPU使用率
- 内存使用量
- 网络I/O
- 磁盘I/O

#### 业务指标
- 集群资源使用率
- 节点健康状态
- 资源创建/删除频率
- 用户活跃度

### 日志管理

#### 日志级别
- **DEBUG**: 详细调试信息
- **INFO**: 一般信息记录
- **WARN**: 警告信息
- **ERROR**: 错误信息

#### 日志格式
```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "level": "info",
  "message": "request completed",
  "traceId": "trace-123",
  "method": "GET",
  "path": "/api/pods",
  "statusCode": 200,
  "latency": "100ms"
}
```

#### 日志收集
- 使用Fluentd/Fluent Bit收集日志
- 发送到Elasticsearch或Loki
- 配置日志轮转和清理

### 健康检查

#### 存活探针
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

#### 就绪探针
```yaml
readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
```

## 🔒 安全配置

### 网络安全

#### 网络策略
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: k8svision-network-policy
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
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          name: kube-system
    ports:
    - protocol: TCP
      port: 443
```

#### TLS配置
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: k8svision-ingress
  namespace: k8svision
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - k8svision.example.com
    secretName: k8svision-tls
  rules:
  - host: k8svision.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: k8svision-frontend
            port:
              number: 80
```

### 访问控制

#### RBAC配置
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: k8svision-role
rules:
- apiGroups: [""]
  resources: ["pods", "services", "nodes", "namespaces"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets", "daemonsets"]
  verbs: ["get", "list", "watch"]
```

#### 认证配置
- 启用JWT认证
- 配置密码策略
- 设置会话超时
- 启用登录限制

## 🔍 故障排查

### 常见问题

#### 1. 服务无法启动
**症状**: Pod处于Pending或CrashLoopBackOff状态

**排查步骤**:
```bash
# 查看Pod状态
kubectl get pods -n k8svision

# 查看Pod详情
kubectl describe pod <pod-name> -n k8svision

# 查看Pod日志
kubectl logs <pod-name> -n k8svision

# 检查资源限制
kubectl top pods -n k8svision
```

**可能原因**:
- 资源不足
- 配置错误
- 依赖服务不可用
- 权限问题

#### 2. 无法连接Kubernetes API
**症状**: 后端服务报错"failed to connect to kubernetes"

**排查步骤**:
```bash
# 检查ServiceAccount
kubectl get serviceaccount k8svision -n k8svision

# 检查RBAC权限
kubectl auth can-i get pods --as=system:serviceaccount:k8svision:k8svision

# 测试API连接
kubectl get nodes
```

**解决方案**:
- 检查RBAC配置
- 验证ServiceAccount
- 确认API服务器地址

#### 3. 前端无法访问后端
**症状**: 前端页面显示网络错误

**排查步骤**:
```bash
# 检查Service配置
kubectl get svc -n k8svision

# 检查Ingress配置
kubectl get ingress -n k8svision

# 测试服务连通性
kubectl port-forward svc/k8svision-backend 8080:8080 -n k8svision
```

#### 4. 认证失败
**症状**: 登录时提示认证失败

**排查步骤**:
```bash
# 检查Secret配置
kubectl get secret k8svision-secrets -n k8svision -o yaml

# 验证JWT密钥
echo "your-jwt-secret" | base64

# 检查认证日志
kubectl logs <backend-pod> -n k8svision | grep auth
```

### 调试工具

#### 日志分析
```bash
# 实时查看日志
kubectl logs -f deployment/k8svision-backend -n k8svision

# 查看错误日志
kubectl logs deployment/k8svision-backend -n k8svision | grep ERROR

# 查看访问日志
kubectl logs deployment/k8svision-backend -n k8svision | grep "request completed"
```

#### 性能分析
```bash
# 查看资源使用情况
kubectl top pods -n k8svision

# 查看节点资源
kubectl top nodes

# 查看详细资源信息
kubectl describe pod <pod-name> -n k8svision
```

#### 网络诊断
```bash
# 测试服务连通性
kubectl exec -it <pod-name> -n k8svision -- curl http://localhost:8080/health

# 检查DNS解析
kubectl exec -it <pod-name> -n k8svision -- nslookup kubernetes.default

# 查看网络策略
kubectl get networkpolicy -n k8svision
```

## 🛠️ 维护指南

### 日常维护

#### 1. 监控检查
- 检查服务健康状态
- 查看资源使用情况
- 监控错误日志
- 检查性能指标

#### 2. 日志管理
- 定期清理旧日志
- 监控日志大小
- 检查日志轮转配置
- 分析错误模式

#### 3. 备份恢复
- 备份配置文件
- 备份数据库
- 测试恢复流程
- 文档化恢复步骤

### 版本升级

#### 1. 升级前准备
```bash
# 备份当前配置
kubectl get configmap k8svision-config -n k8svision -o yaml > config-backup.yaml

# 备份Secret
kubectl get secret k8svision-secrets -n k8svision -o yaml > secrets-backup.yaml

# 检查当前版本
kubectl get deployment k8svision-backend -n k8svision -o jsonpath='{.spec.template.spec.containers[0].image}'
```

#### 2. 执行升级
```bash
# 更新镜像
kubectl set image deployment/k8svision-backend k8svision-backend=k8svision:latest -n k8svision

# 等待升级完成
kubectl rollout status deployment/k8svision-backend -n k8svision

# 验证升级结果
kubectl get pods -n k8svision
```

#### 3. 回滚操作
```bash
# 查看升级历史
kubectl rollout history deployment/k8svision-backend -n k8svision

# 回滚到上一版本
kubectl rollout undo deployment/k8svision-backend -n k8svision

# 回滚到指定版本
kubectl rollout undo deployment/k8svision-backend --to-revision=2 -n k8svision
```

### 性能优化

#### 1. 资源调优
- 调整CPU和内存限制
- 优化JVM参数
- 配置缓存大小
- 调整并发数

#### 2. 网络优化
- 启用HTTP/2
- 配置连接池
- 优化超时设置
- 启用压缩

#### 3. 存储优化
- 使用SSD存储
- 配置存储类
- 优化I/O性能
- 定期清理数据

### 安全维护

#### 1. 定期更新
- 更新基础镜像
- 升级依赖包
- 应用安全补丁
- 更新证书

#### 2. 安全审计
- 检查访问日志
- 分析异常行为
- 验证权限配置
- 测试安全策略

#### 3. 备份安全
- 加密备份数据
- 安全存储备份
- 定期测试恢复
- 文档化流程

## 📞 支持联系

如有问题或需要支持，请：
1. 查看本文档的故障排查部分
2. 检查项目的Issues页面
3. 联系技术支持团队
4. 提交详细的错误报告

---

**K8sVision** - 让Kubernetes管理更简单、更可靠！
