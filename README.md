# K8sVision

K8sVision 是一款现代化的 Kubernetes 可视化监控与资源管理平台，前后端分离，支持一键本地部署与企业级 K8s 集群部署，界面美观、体验流畅、易于二次开发。

---

## 项目亮点

- **全资源可视化**：支持 Node、Pod、Deployment、StatefulSet、DaemonSet、Service、Ingress、Job、CronJob、Event、Namespace 等核心资源的实时监控与管理。
- **极致工程规范**：后端 Go + Gin + zap，前端 React + Vite，代码结构清晰，接口高度统一，注释与文档齐全。
- **多环境部署**：支持本地 docker-compose 一键部署，支持 K8s 集群生产级部署，环境变量灵活注入。
- **API 文档自动生成**：内置 Swagger，接口自描述，便于前后端协作与二次开发。
- **高可维护性**：分页、错误处理、日志、依赖注入、配置管理等全部工程化，易于扩展。

---

## 目录结构

```
K8sVision/
├── backend/                # 后端 Go 服务
│   ├── api/                # API 路由与 handler
│   ├── service/            # 业务逻辑与 K8s 采集
│   ├── model/              # 数据结构与常量
│   ├── docs/               # Swagger 文档
│   ├── main.go             # 入口
│   ├── Dockerfile          # 后端镜像构建
│   └── README.md
├── frontend/               # 前端 React 应用
│   ├── src/                # 页面、组件、样式
│   ├── public/             # 静态资源（如 favicon）
│   ├── Dockerfile          # 前端镜像构建
│   ├── nginx.conf          # Nginx 配置
│   └── README.md
├── k8s/                    # K8s 部署 YAML
│   ├── backend-deployment.yaml
│   ├── frontend-deployment.yaml
│   └── ingress.yaml
├── docker-compose.yml      # 一键本地部署
├── go.mod / go.sum         # Go 依赖
├── config.yaml             # 后端配置
└── README.md               # 项目说明（本文件）
```

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

## 环境变量与配置

### 后端（backend）

| 变量名         | 说明                   | 默认值         |
| -------------- | ---------------------- | -------------- |
| SERVER_PORT    | 服务监听端口           | 8080           |
| CONFIG_PATH    | 配置文件路径           | /app/config.yaml |

- 也可通过 `config.yaml` 配置 K8s 连接、日志等参数

### 前端（frontend）

| 变量名         | 说明                   | 默认值         |
| -------------- | ---------------------- | -------------- |
| VITE_API_URL   | 后端 API 基础路径      | /api           |

- 通过 nginx.conf 代理 `/api/` 到后端服务

---

## 主要功能

- **集群总览**：节点、Pod、命名空间、服务等核心指标一览
- **资源列表与详情**：所有 K8s 资源的分页、检索、状态展示
- **事件追踪**：实时展示集群 Event，支持筛选与详情
- **资源用量可视化**：CPU/内存请求、限制、分配、剩余等图表
- **健康检查**：内置 `/healthz` 路由，便于 K8s/运维监控
- **API 文档**：Swagger 自动生成，接口自描述

---

## 前端开发

```bash
cd frontend
npm install
npm run dev
# 默认端口 5173，支持热更新
```

## 后端开发

```bash
cd backend
go mod tidy
go run main.go
# 默认端口 8080
```

---

## K8s 部署说明

- **Service 名称即为 DNS 域名**，nginx.conf 代理可直接用 `k8svision-backend:8080`
- **Ingress** 支持前后端分流，推荐用独立域名或路径前缀
- **健康检查** 已内置，K8s 自动探活

---

## 常见问题

- **Q: 前端访问 404？**
  - A: 请确保 nginx.conf 配置了 `try_files $uri $uri/ /index.html;`，支持 history 路由。
- **Q: API 代理失败？**
  - A: 请检查 nginx.conf 的 `proxy_pass` 是否指向正确的 K8s Service 名称和端口。
- **Q: 如何自定义配置？**
  - A: 后端支持通过环境变量或 config.yaml 注入，前端支持 VITE_API_URL 环境变量。

---

## 贡献与二次开发

- 欢迎 PR、Issue、建议与二次开发！
- 代码结构清晰，注释齐全，接口自描述，易于团队协作。

---

## License

MIT
