# K8sVision API 文档

## 概述

K8sVision 是一个 Kubernetes 集群管理平台，提供完整的集群资源监控和管理功能。

## API 特性

### 认证
- 基于 JWT 的认证机制
- 支持登录限流保护（连续失败5次10分钟内禁止尝试）

### 分页支持
所有列表API都支持分页查询：
- `limit`: 每页数量（默认20）
- `offset`: 偏移量（默认0）

### 搜索功能
大部分列表API现在支持关键词搜索，支持以下字段搜索：

#### 工作负载资源
- **Deployments**: 名称、命名空间、状态等
- **StatefulSets**: 名称、命名空间、状态等  
- **DaemonSets**: 名称、命名空间、状态等
- **Jobs**: 名称、命名空间、状态等
- **CronJobs**: 名称、命名空间、状态等

#### 网络资源
- **Services**: 名称、命名空间、类型等
- **Ingress**: 名称、命名空间、主机等
- **Pods**: 名称、命名空间、状态、PodIP、节点等

#### 存储资源
- **PVCs**: 名称、命名空间、状态等
- **PVs**: 名称、状态、存储类等
- **StorageClasses**: 名称、供应者等

#### 配置资源
- **ConfigMaps**: 名称、命名空间等
- **Secrets**: 名称、命名空间、类型等

#### 集群资源
- **Nodes**: 名称、状态、角色等
- **Namespaces**: 名称等
- **Events**: 名称、命名空间、原因、消息等

## API 端点

### 认证
- `POST /api/login` - 用户登录

### 集群概览
- `GET /api/overview` - 获取集群资源总览

### 工作负载
- `GET /api/deployments` - 获取 Deployment 列表
- `GET /api/deployments/{namespace}/{name}` - 获取 Deployment 详情
- `GET /api/statefulsets` - 获取 StatefulSet 列表
- `GET /api/statefulsets/{namespace}/{name}` - 获取 StatefulSet 详情
- `GET /api/daemonsets` - 获取 DaemonSet 列表
- `GET /api/daemonsets/{namespace}/{name}` - 获取 DaemonSet 详情
- `GET /api/jobs` - 获取 Job 列表
- `GET /api/jobs/{namespace}/{name}` - 获取 Job 详情
- `GET /api/cronjobs` - 获取 CronJob 列表
- `GET /api/cronjobs/{namespace}/{name}` - 获取 CronJob 详情

### 网络
- `GET /api/pods` - 获取 Pod 列表
- `GET /api/pods/{namespace}/{name}` - 获取 Pod 详情
- `GET /api/services` - 获取 Service 列表
- `GET /api/services/{namespace}/{name}` - 获取 Service 详情
- `GET /api/ingress` - 获取 Ingress 列表
- `GET /api/ingress/{namespace}/{name}` - 获取 Ingress 详情

### 存储
- `GET /api/pvcs` - 获取 PVC 列表
- `GET /api/pvcs/{namespace}/{name}` - 获取 PVC 详情
- `GET /api/pvs` - 获取 PV 列表
- `GET /api/pvs/{name}` - 获取 PV 详情
- `GET /api/storageclasses` - 获取 StorageClass 列表
- `GET /api/storageclasses/{name}` - 获取 StorageClass 详情

### 配置
- `GET /api/configmaps` - 获取 ConfigMap 列表
- `GET /api/configmaps/{namespace}/{name}` - 获取 ConfigMap 详情
- `GET /api/secrets` - 获取 Secret 列表
- `GET /api/secrets/{namespace}/{name}` - 获取 Secret 详情

### 集群管理
- `GET /api/nodes` - 获取 Node 列表
- `GET /api/nodes/{name}` - 获取 Node 详情
- `GET /api/namespaces` - 获取 Namespace 列表
- `GET /api/namespaces/{name}` - 获取 Namespace 详情
- `GET /api/events` - 获取 Event 列表
- `GET /api/events/{namespace}/{name}` - 获取 Event 详情

## 响应格式

所有API都返回统一的响应格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "page": {
    "total": 100,
    "limit": 20,
    "offset": 0
  },
  "timestamp": 1640995200,
  "traceId": "abc123"
}
```

## 使用示例

### 搜索Deployments
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:8080/api/deployments?search=nginx&limit=10&offset=0"
```

### 分页查询Pods
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:8080/api/pods?namespace=default&limit=20&offset=40"
```

### 获取集群概览
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:8080/api/overview"
```

## 更新日志

### v1.0 (最新)
- 新增搜索功能支持，所有列表API现在支持关键词搜索
- 更新API描述，明确说明分页和搜索功能
- 优化swagger文档结构，提供更详细的参数说明

## 开发说明

### 本地开发
```bash
# 启用swagger文档
export SWAGGER_ENABLE=true

# 启动服务
go run main.go

# 访问swagger文档
http://localhost:8080/swagger/index.html
```

### 生成swagger文档
项目使用 `swag` 工具生成swagger文档，确保在修改API后重新生成文档。

## 支持

如有问题或建议，请提交Issue或联系开发团队。 