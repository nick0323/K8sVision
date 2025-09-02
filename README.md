# K8sVision

一个现代化的 Kubernetes 集群管理平台，提供直观的 Web 界面来管理和监控 Kubernetes 资源。

## 🚀 功能特性

### 📊 集群概览
- 实时集群状态监控
- 资源使用情况统计
- 节点健康状态检查

### 🔧 工作负载管理
- **Pod**: 容器实例管理，支持状态监控和日志查看
- **Deployment**: 应用部署管理，支持滚动更新
- **StatefulSet**: 有状态应用管理
- **DaemonSet**: 守护进程管理
- **Job**: 批处理任务管理
- **CronJob**: 定时任务管理

### 🌐 网络管理
- **Service**: 服务发现和负载均衡
- **Ingress**: 外部访问路由管理

### 💾 存储管理
- **PVC**: 持久卷声明管理
- **PV**: 持久卷管理
- **StorageClass**: 存储类管理

### ⚙️ 配置管理
- **ConfigMap**: 配置数据管理
- **Secret**: 敏感信息管理

### 🏷️ 集群管理
- **Namespace**: 命名空间管理
- **Node**: 节点管理
- **Event**: 事件监控

## 🏗️ 技术架构

### 后端技术栈
- **Go 1.24.4**: 主要编程语言
- **Gin**: Web 框架
- **Kubernetes Client-go**: K8s API 客户端
- **JWT**: 身份认证
- **Zap**: 结构化日志
- **Viper**: 配置管理
- **Swagger**: API 文档

### 前端技术栈
- **React 19.1.0**: UI 框架
- **Vite**: 构建工具
- **React Icons**: 图标库
- **ESLint**: 代码检查

### 核心特性
- 🔐 **JWT 身份认证**: 安全的用户认证机制
- 📈 **实时监控**: 集群资源实时状态监控
- 🔍 **智能搜索**: 支持多字段搜索和过滤
- 📄 **分页支持**: 高效的数据分页展示
- 🎨 **响应式设计**: 适配不同屏幕尺寸
- 🚀 **性能优化**: 懒加载和缓存机制
- 📊 **详情页面**: 丰富的资源详情展示

## 📦 安装部署

### 环境要求
- Go 1.24.4+
- Node.js 20+
- Kubernetes 集群访问权限

### 使用 Docker Compose 部署

1. **克隆项目**
```bash
git clone https://github.com/nick0323/K8sVision.git
cd K8sVision
```

2. **配置 Kubernetes 访问**
```bash
# 确保 kubeconfig 文件存在
cp ~/.kube/config ~/.kube/config.backup
```

3. **启动服务**
```bash
docker-compose up -d
```

4. **访问应用**
- 前端: http://localhost:80
- 后端: http://localhost:8080
- 默认登录: admin / 12345678

### 手动部署

#### 后端部署
```bash
# 构建后端
go build -o k8svision ./main.go

# 运行后端
./k8svision -config config.yaml
```

#### 前端部署
```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run dev

# 生产构建
npm run build
```

### Kubernetes 部署

```bash
# 部署到 Kubernetes 集群
kubectl apply -f k8s/
```

## ⚙️ 配置说明

### 后端配置 (config.yaml)

```yaml
server:
  port: 8080
  host: "0.0.0.0"

kubernetes:
  kubeconfig: ""  # 留空使用默认配置
  context: ""     # 留空使用当前上下文
  timeout: 30s
  qps: 100
  burst: 200

jwt:
  secret: "k8svision-secret-key"
  expiration: 24h

auth:
  username: "admin"
  password: "admin"
  maxLoginFail: 5
  lockDuration: 10m

cache:
  enabled: true
  type: "memory"
  ttl: 5m
  maxSize: 1000
```

### 环境变量配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `KUBECONFIG` | Kubernetes 配置文件路径 | `~/.kube/config` |
| `LOGIN_USERNAME` | 登录用户名 | `admin` |
| `LOGIN_PASSWORD` | 登录密码 | `admin` |
| `GIN_MODE` | Gin 运行模式 | `release` |
| `SWAGGER_ENABLE` | 启用 Swagger 文档 | `false` |

## 🔌 API 接口

### 认证接口
- `POST /api/login` - 用户登录

### 资源管理接口
- `GET /api/overview` - 集群概览
- `GET /api/namespaces` - 命名空间列表
- `GET /api/nodes` - 节点列表
- `GET /api/pods` - Pod 列表
- `GET /api/deployments` - Deployment 列表
- `GET /api/services` - Service 列表
- `GET /api/events` - 事件列表
- 更多接口请参考 Swagger 文档

### 通用查询参数
- `namespace`: 命名空间过滤
- `limit`: 每页数量 (默认: 20)
- `offset`: 偏移量 (默认: 0)
- `search`: 搜索关键词

## 🎨 界面预览

### 主要功能
- 📊 **集群概览**: 实时显示集群整体状态
- 📋 **资源列表**: 分页展示各类 Kubernetes 资源
- 🔍 **智能搜索**: 支持多字段模糊搜索
- 📄 **详情查看**: 点击资源查看详细信息
- 🎛️ **命名空间切换**: 快速切换不同命名空间
- 🔄 **实时刷新**: 自动或手动刷新数据

### 响应式设计
- 💻 桌面端: 完整功能展示
- 📱 移动端: 适配小屏幕操作
- 🎨 现代化 UI: 清晰的信息层次

## 🛠️ 开发指南

### 项目结构
```
K8sVision/
├── api/                 # API 接口层
├── cache/              # 缓存管理
├── config/             # 配置管理
├── frontend/           # 前端代码
│   ├── src/
│   │   ├── components/ # React 组件
│   │   ├── hooks/      # 自定义 Hooks
│   │   ├── utils/      # 工具函数
│   │   └── constants/  # 常量配置
│   └── dist/           # 构建产物
├── k8s/                # Kubernetes 部署文件
├── model/              # 数据模型
├── monitor/            # 监控模块
├── service/            # 业务逻辑层
├── main.go             # 主程序入口
└── config.yaml         # 配置文件
```

### 开发环境设置
```bash
# 后端开发
go mod download
go run main.go

# 前端开发
cd frontend
npm install
npm run dev
```

### 代码规范
- 后端: 遵循 Go 官方代码规范
- 前端: 使用 ESLint 进行代码检查
- 提交: 使用有意义的提交信息

## 🧪 测试

### 后端测试
```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./api/middleware
```

### 前端测试
```bash
cd frontend
npm run lint
```

## 📈 性能优化

### 后端优化
- 🚀 **缓存机制**: 内存缓存减少 API 调用
- ⚡ **并发控制**: 限制并发请求数量
- 📊 **监控指标**: 实时性能监控

### 前端优化
- 🔄 **懒加载**: 按需加载页面组件
- 💾 **数据缓存**: 减少重复 API 请求
- 🎯 **虚拟滚动**: 处理大量数据展示

## 🔒 安全特性

- 🔐 **JWT 认证**: 安全的用户身份验证
- 🚫 **登录限制**: 防止暴力破解
- 🔒 **HTTPS 支持**: 加密数据传输
- 🛡️ **输入验证**: 防止注入攻击

## 🐛 故障排除

### 常见问题

1. **无法连接 Kubernetes 集群**
   - 检查 kubeconfig 文件路径
   - 确认集群访问权限
   - 验证网络连接

2. **登录失败**
   - 检查用户名密码配置
   - 确认 JWT 密钥设置
   - 查看后端日志

3. **前端无法访问后端**
   - 检查后端服务状态
   - 确认端口配置
   - 验证 CORS 设置

### 日志查看
```bash
# 查看后端日志
docker logs vision-backend

# 查看前端日志
docker logs vision-frontend
```

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 Apache 2.0 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Kubernetes](https://kubernetes.io/) - 容器编排平台
- [Gin](https://gin-gonic.com/) - Go Web 框架
- [React](https://reactjs.org/) - 前端框架
- [Vite](https://vitejs.dev/) - 构建工具

## 📞 联系方式

- 项目链接: [https://github.com/nick0323/K8sVision](https://github.com/nick0323/K8sVision)
- 问题反馈: [Issues](https://github.com/nick0323/K8sVision/issues)

---

⭐ 如果这个项目对您有帮助，请给它一个星标！
