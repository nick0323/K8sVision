# K8sVision

K8sVision 是一款现代化的 Kubernetes 可视化监控与资源管理平台，采用前后端分离架构，支持一键本地部署与企业级 K8s 集群生产部署。界面美观、体验流畅、工程规范，适合二次开发与团队协作。

---

## 目录结构

```
K8sVision/
├── api/                  # Go 后端 API 路由与 handler
├── service/              # 业务逻辑与 K8s 采集
├── model/                # 数据结构与常量
├── docs/                 # Swagger 文档
├── main.go               # Go 入口
├── config.yaml           # 后端配置
├── go.mod / go.sum       # Go 依赖管理
├── Dockerfile            # 后端镜像构建
├── frontend/             # React 前端应用
│   ├── src/              # 页面、组件、样式
│   ├── public/           # 静态资源
│   ├── Dockerfile        # 镜像构建
│   ├── nginx.conf        # Nginx 配置
│   └── package.json      # 前端依赖
├── k8s/                  # K8s 部署 YAML
│   ├── backend-deployment.yaml
│   ├── frontend-deployment.yaml
│   └── ingress.yaml
├── docker-compose.yml    # 一键本地部署
└── README.md             # 项目说明
```

---

## 功能特性

- **全资源可视化**：支持 Node、Pod、Deployment、StatefulSet、DaemonSet、Service、Ingress、Job、CronJob、Event、Namespace 等核心资源的实时监控与管理。
- **集群总览**：节点、Pod、命名空间、服务等核心指标一览。
- **资源列表与详情**：所有 K8s 资源的分页、检索、状态展示。
- **事件追踪**：实时展示集群 Event，支持筛选与详情。
- **资源用量可视化**：CPU/内存请求、限制、分配、剩余等图表。
- **健康检查**：内置 `/healthz` 路由，便于 K8s/运维监控。
- **API 文档**：Swagger 自动生成，接口自描述。
- **极致工程规范**：后端 Go + Gin + zap，前端 React + Vite，代码结构清晰，接口高度统一，注释与文档齐全。
- **多环境部署**：支持本地 docker-compose 一键部署，支持 K8s 集群生产级部署，环境变量灵活注入。
- **高可维护性**：分页、错误处理、日志、依赖注入、配置管理等全部工程化，易于扩展。

---

## 快速开始

### 1. 本地一键部署（推荐开发体验）

需安装 [Docker](https://www.docker.com/) 和 [docker-compose](https://docs.docker.com/compose/)

```bash
git clone https://github.com/nick0323/K8sVision.git
cd K8sVision
docker-compose up -d
```

- 访问前端：http://localhost
- 后端 API：http://localhost:8080/api
- Swagger 文档：http://localhost:8080/swagger/index.html

### 2. Kubernetes 集群部署

1. 构建并推送镜像（或用 CI 自动推送到镜像仓库）
2. 修改 `k8s/*.yaml` 中镜像名（如需）
3. 部署：

```bash
kubectl apply -f k8s/
```

- 访问方式：通过 Ingress 域名（如 `k8svision.local`，可本地 hosts 解析或用实际域名）

---

## 配置说明

### 后端（Go 服务）

- 配置文件：`config.yaml`（位于项目根目录）
- 依赖管理：`go.mod`、`go.sum`（位于项目根目录）
- 主要源码目录：`api/`、`service/`、`model/`、`docs/`、`main.go`

| 变量名           | 说明                   | 默认值         |
| ---------------- | ---------------------- | -------------- |
| SERVER_PORT      | 服务监听端口           | 8080           |
| CONFIG_PATH      | 配置文件路径           | ./config.yaml  |
| LOGIN_USERNAME   | 登录用户名             | admin          |
| LOGIN_PASSWORD   | 登录密码               | 123456         |
| JWT_SECRET       | JWT 签发密钥           | k8svision-secret-key |
| KUBECONFIG       | kubeconfig 路径        | ""             |

- 推荐通过环境变量安全注入账号密码和密钥，避免硬编码。
- 单集群配置示例：

```yaml
server:
  port: 8080
  logLevel: info

kubernetes:
  kubeconfig: ""

jwt:
  secret: ""
```

### 前端（frontend）

- 配置文件：`.env` 或环境变量
- 主要变量：

| 变量名         | 说明                   | 默认值         |
| -------------- | ---------------------- | -------------- |
| VITE_API_URL   | 后端 API 基础路径      | /api           |

- 通过 nginx.conf 代理 `/api/` 到后端服务

---

## 前端开发

```bash
cd frontend
npm install
npm run dev
# 默认端口 5173，支持热更新
```

- 主要依赖：React 19、Vite 7、react-icons
- 代码结构清晰，支持组件化开发与二次扩展

---

## 后端开发

```bash
# 以下命令均在项目根目录下执行
go mod tidy
go run main.go
# 默认端口 8080
```

- 主要依赖：Gin、zap、viper、swaggo
- 支持 JWT 认证、Swagger 文档、健康检查

---

## K8s 部署说明

- **Service 名称即为 DNS 域名**，nginx.conf 代理可直接用 `k8svision-backend:8080`
- **Ingress** 支持前后端分流，推荐用独立域名或路径前缀
- **健康检查** 已内置，K8s 自动探活
- 参考 `k8s/` 目录下 YAML 文件，按需修改镜像名和资源规格

---

## 主要接口与安全

- 登录接口：`/api/login`，JWT 认证，防爆破（连续失败5次，10分钟内禁止）
- 健康检查：`/healthz`，GET，返回 200 OK
- 资源接口：全部以 `/api/` 开头，需携带 Bearer Token
- Swagger 文档：`/swagger/index.html`，自动生成

---

## 常见问题

- **前端访问 404？**
  - 请确保 nginx.conf 配置了 `try_files $uri $uri/ /index.html;`，支持 history 路由。
- **API 代理失败？**
  - 请检查 nginx.conf 的 `proxy_pass` 是否指向正确的 K8s Service 名称和端口。
- **如何自定义配置？**
  - 后端支持通过环境变量或 config.yaml 注入，前端支持 VITE_API_URL 环境变量。

---

## 贡献与二次开发

- 欢迎 PR、Issue、建议与二次开发！
- 代码结构清晰，注释齐全，接口自描述，易于团队协作。

---

## License

MIT

---

如需更详细的接口说明、二次开发指引或架构解读，请查阅源码注释与 Swagger 文档。此 README 已覆盖项目架构、功能、部署、配置、开发、常见问题等核心内容，适合开源、企业或团队协作场景。
