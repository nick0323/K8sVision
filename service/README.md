# Service模块文档

Service模块是K8sVision的业务逻辑层，负责与Kubernetes API交互、数据处理和业务规则实现。

## 📁 模块结构

```
service/
├── README.md                    # 模块文档
├── common.go                    # 通用服务工具
├── k8s.go                       # Kubernetes客户端管理
├── overview.go                  # 集群概览服务
├── namespace.go                 # 命名空间服务
├── node.go                      # 节点服务
├── pod.go                       # Pod服务
├── deployment.go                # Deployment服务
├── statefulset.go               # StatefulSet服务
├── daemonset.go                 # DaemonSet服务
├── service.go                   # Service服务
├── ingress.go                   # Ingress服务
├── job.go                       # Job服务
├── cronjob.go                   # CronJob服务
├── event.go                     # Event服务
├── pvc.go                       # PVC服务
├── pv.go                        # PV服务
├── storageclass.go              # StorageClass服务
├── configmap.go                 # ConfigMap服务
├── secret.go                    # Secret服务
```

## 🔧 核心组件

### 1. Kubernetes客户端管理 (k8s.go)
负责Kubernetes客户端的创建和管理：
- 客户端连接管理
- 配置加载
- 连接缓存
- 错误处理

**主要函数：**
- `GetK8sClient`: 获取Kubernetes客户端
- `GetK8sConfig`: 获取Kubernetes配置
- `applyK8sConfig`: 应用Kubernetes配置

### 2. 通用服务工具 (common.go)
提供业务层的通用功能：
- 资源状态提取
- 数据转换
- 工具函数

**主要函数：**
- `GetJobStatus`: 获取Job状态
- `SafeInt32Ptr`: 安全获取Int32指针
- `SafeBoolPtr`: 安全获取Bool指针
- `FormatDuration`: 格式化持续时间

### 3. 资源服务模块
每个Kubernetes资源都有对应的服务模块：

#### Pod服务 (pod.go)
- 获取Pod列表和详情
- 处理Pod状态和指标
- 支持命名空间过滤

#### Deployment服务 (deployment.go)
- 管理Deployment资源
- 状态计算和转换
- 副本数统计

#### Service服务 (service.go)
- 处理Service资源
- 端口信息提取
- 类型转换

## 🚀 主要功能

### 资源列表获取
所有资源服务都提供列表获取功能：
```go
func ListPods(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]model.PodStatus, error)
```

**特性：**
- 支持命名空间过滤
- 分页查询
- 状态计算
- 错误处理

### 资源详情获取
提供单个资源的详细信息：
```go
func GetPodDetail(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) (*model.PodStatus, error)
```

**特性：**
- 完整资源信息
- 状态分析
- 关联资源查询

### 集群概览 (overview.go)
提供集群级别的统计信息：
- 节点统计
- 资源使用情况
- 健康状态
- 性能指标

## 📊 数据处理

### 状态计算
为每个资源计算标准化的状态信息：
```go
type PodStatus struct {
    Name        string    `json:"name"`
    Namespace   string    `json:"namespace"`
    Status      string    `json:"status"`
    Ready       string    `json:"ready"`
    Restarts    int32     `json:"restarts"`
    Age         string    `json:"age"`
    // ... 更多字段
}
```

### 数据转换
将Kubernetes原生对象转换为前端友好的格式：
- 时间格式化
- 状态标准化
- 资源单位转换
- 嵌套对象扁平化

### 错误处理
统一的错误处理机制：
- 连接错误处理
- 权限错误处理
- 资源不存在处理
- 超时处理

## 🔒 安全特性

### 权限控制
- 基于Kubernetes RBAC
- 最小权限原则
- 资源访问验证

### 数据安全
- 敏感信息过滤
- 密码脱敏
- 配置信息保护

### 连接安全
- TLS证书验证
- 令牌认证
- 连接超时控制

## 📈 性能优化

### 并发处理
- 并发安全的客户端管理
- 连接池复用
- 请求限流

### 缓存策略
- 客户端连接缓存
- 配置信息缓存
- 结果缓存

### 资源优化
- 内存使用优化
- 网络请求优化
- 数据处理优化

## 🛠️ 开发指南

### 添加新资源服务
1. 创建资源服务文件
2. 实现列表和详情函数
3. 定义资源状态结构
4. 添加错误处理
5. 编写单元测试

### 状态计算模式
```go
func calculateResourceStatus(resource *v1.Resource) model.ResourceStatus {
    return model.ResourceStatus{
        Name:      resource.Name,
        Namespace: resource.Namespace,
        Status:    determineStatus(resource),
        Age:       calculateAge(resource.CreationTimestamp),
        // ... 其他字段
    }
}
```

### 错误处理模式
```go
func ListResources(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]model.ResourceStatus, error) {
    // 获取资源列表
    list, err := clientset.CoreV1().Resources(namespace).List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to list resources: %w", err)
    }
    
    // 转换数据
    var resources []model.ResourceStatus
    for _, item := range list.Items {
        resources = append(resources, calculateResourceStatus(&item))
    }
    
    return resources, nil
}
```

## 🧪 测试

### 单元测试
每个服务模块都包含对应的测试文件：
- 功能测试
- 边界条件测试
- 错误情况测试

### 测试工具
- Mock Kubernetes客户端
- 测试数据生成
- 断言工具

### 测试覆盖
- 主要功能路径
- 错误处理路径
- 边界条件

## 📝 最佳实践

1. **资源管理**
   - 及时关闭资源
   - 避免内存泄漏
   - 合理使用缓存

2. **错误处理**
   - 提供有意义的错误信息
   - 记录详细的错误日志
   - 优雅降级

3. **性能考虑**
   - 避免重复查询
   - 使用并发处理
   - 优化数据结构

4. **代码质量**
   - 保持函数简洁
   - 添加必要注释
   - 遵循Go语言规范

## 🔍 故障排查

### 常见问题
1. **连接失败**
   - 检查Kubernetes配置
   - 验证网络连接
   - 确认权限设置

2. **权限错误**
   - 检查RBAC配置
   - 验证ServiceAccount
   - 查看错误日志

3. **性能问题**
   - 监控资源使用
   - 检查查询效率
   - 分析缓存命中率

### 调试工具
- 详细日志记录
- 性能指标监控
- 健康检查接口
- 连接状态检查
