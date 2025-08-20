# K8sVision API 文档

欢迎使用 K8sVision API 文档！本文档提供了完整的 RESTful API 接口说明。

## 📋 目录

- [认证说明](./authentication.md) - JWT 认证和权限控制
- [错误处理](./errors.md) - 错误代码和响应格式
- [通用接口](./common.md) - 通用响应格式和分页
- [资源接口](./resources/README.md) - 各资源类型的 API
- [监控接口](./monitoring.md) - 监控和指标接口

## 🚀 快速开始

### 基础信息
- **Base URL**: `http://localhost:8080/api`
- **协议**: HTTP/HTTPS
- **数据格式**: JSON
- **认证方式**: JWT Bearer Token

### 认证流程
1. 调用登录接口获取 Token
2. 在请求头中携带 Token
3. 访问受保护的资源接口

### 示例请求
```bash
# 1. 登录获取 Token
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"12345678"}'

# 2. 使用 Token 访问资源
curl -X GET http://localhost:8080/api/nodes \
  -H "Authorization: Bearer <your-token>"
```

## 📊 API 概览

### 认证接口
| 接口 | 方法 | 说明 | 认证 |
|------|------|------|------|
| `/api/login` | POST | 用户登录 | 否 |

### 集群管理接口
| 接口 | 方法 | 说明 | 认证 |
|------|------|------|------|
| `/api/overview` | GET | 集群概览 | 是 |
| `/api/events` | GET | 事件列表 | 是 |

### 计算资源接口
| 接口 | 方法 | 说明 | 认证 |
|------|------|------|------|
| `/api/nodes` | GET | 节点列表 | 是 |
| `/api/nodes/{name}` | GET | 节点详情 | 是 |
| `/api/pods` | GET | Pod 列表 | 是 |
| `/api/pods/{namespace}/{name}` | GET | Pod 详情 | 是 |

### 工作负载接口
| 接口 | 方法 | 说明 | 认证 |
|------|------|------|------|
| `/api/deployments` | GET | Deployment 列表 | 是 |
| `/api/deployments/{namespace}/{name}` | GET | Deployment 详情 | 是 |
| `/api/statefulsets` | GET | StatefulSet 列表 | 是 |
| `/api/statefulsets/{namespace}/{name}` | GET | StatefulSet 详情 | 是 |
| `/api/daemonsets` | GET | DaemonSet 列表 | 是 |
| `/api/daemonsets/{namespace}/{name}` | GET | DaemonSet 详情 | 是 |

### 网络资源接口
| 接口 | 方法 | 说明 | 认证 |
|------|------|------|------|
| `/api/services` | GET | Service 列表 | 是 |
| `/api/services/{namespace}/{name}` | GET | Service 详情 | 是 |
| `/api/ingresses` | GET | Ingress 列表 | 是 |
| `/api/ingresses/{namespace}/{name}` | GET | Ingress 详情 | 是 |

### 存储资源接口
| 接口 | 方法 | 说明 | 认证 |
|------|------|------|------|
| `/api/pvcs` | GET | PVC 列表 | 是 |
| `/api/pvcs/{namespace}/{name}` | GET | PVC 详情 | 是 |
| `/api/pvs` | GET | PV 列表 | 是 |
| `/api/pvs/{name}` | GET | PV 详情 | 是 |
| `/api/storageclasses` | GET | StorageClass 列表 | 是 |

### 配置资源接口
| 接口 | 方法 | 说明 | 认证 |
|------|------|------|------|
| `/api/configmaps` | GET | ConfigMap 列表 | 是 |
| `/api/configmaps/{namespace}/{name}` | GET | ConfigMap 详情 | 是 |
| `/api/secrets` | GET | Secret 列表 | 是 |
| `/api/secrets/{namespace}/{name}` | GET | Secret 详情 | 是 |

### 工作负载接口
| 接口 | 方法 | 说明 | 认证 |
|------|------|------|------|
| `/api/jobs` | GET | Job 列表 | 是 |
| `/api/jobs/{namespace}/{name}` | GET | Job 详情 | 是 |
| `/api/cronjobs` | GET | CronJob 列表 | 是 |
| `/api/cronjobs/{namespace}/{name}` | GET | CronJob 详情 | 是 |

### 监控接口
| 接口 | 方法 | 说明 | 认证 |
|------|------|------|------|
| `/metrics` | GET | 性能指标 | 否 |
| `/cache/stats` | GET | 缓存统计 | 否 |

## 🔐 认证说明

### JWT Token 格式
```
Authorization: Bearer <jwt-token>
```

### Token 结构
```json
{
  "username": "admin",
  "exp": 1703123456,
  "iat": 1703037056
}
```

### 登录接口
```bash
POST /api/login
Content-Type: application/json

{
  "username": "admin",
  "password": "12345678"
}
```

### 响应格式
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "timestamp": "2024-12-20T10:30:56Z"
}
```

## 📝 通用响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    // 具体数据
  },
  "page": {
    "total": 100,
    "limit": 10,
    "offset": 0
  },
  "timestamp": "2024-12-20T10:30:56Z"
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "参数错误",
  "details": "用户名和密码不能为空",
  "timestamp": "2024-12-20T10:30:56Z"
}
```

## 🔍 查询参数

### 分页参数
| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `limit` | int | 每页数量 | 10 |
| `offset` | int | 偏移量 | 0 |

### 过滤参数
| 参数 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `namespace` | string | 命名空间过滤 | `default` |
| `labelSelector` | string | 标签选择器 | `app=nginx` |
| `fieldSelector` | string | 字段选择器 | `status.phase=Running` |

### 排序参数
| 参数 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `sortBy` | string | 排序字段 | `name` |
| `sortOrder` | string | 排序方向 | `asc` 或 `desc` |

## 🚨 错误处理

### HTTP 状态码
- `200` - 成功
- `400` - 请求参数错误
- `401` - 未认证
- `403` - 权限不足
- `404` - 资源不存在
- `429` - 请求频率限制
- `500` - 服务器内部错误

### 错误代码
| 代码 | 说明 | HTTP 状态码 |
|------|------|-------------|
| `0` | 成功 | 200 |
| `400` | 参数错误 | 400 |
| `401` | 未认证 | 401 |
| `403` | 权限不足 | 403 |
| `404` | 资源不存在 | 404 |
| `429` | 请求频率限制 | 429 |
| `500` | 服务器内部错误 | 500 |

## 📊 数据模型

### 通用字段
所有资源都包含以下通用字段：
- `name` - 资源名称
- `namespace` - 命名空间（集群级资源除外）
- `labels` - 标签
- `annotations` - 注解
- `creationTimestamp` - 创建时间
- `status` - 状态

### 状态字段
- `Running` - 运行中
- `Pending` - 等待中
- `Failed` - 失败
- `Succeeded` - 成功
- `Unknown` - 未知

## 🔧 开发工具

### Swagger 文档
访问 http://localhost:8080/swagger/index.html 查看交互式 API 文档。

### 健康检查
```bash
curl http://localhost:8080/healthz
```

### 性能指标
```bash
curl http://localhost:8080/metrics
```

## 📞 支持

如果您在使用 API 时遇到问题，请：

1. 查看 [错误处理](./errors.md) 文档
2. 检查 [常见问题](../troubleshooting/faq.md)
3. 提交 [GitHub Issue](https://github.com/nick0323/K8sVision/issues)

## 📚 相关文档

- [认证说明](./authentication.md)
- [错误处理](./errors.md)
- [通用接口](./common.md)
- [资源接口](./resources/README.md)
- [监控接口](./monitoring.md)

---

**API 版本**: v1.0.0  
**最后更新**: 2024年12月 