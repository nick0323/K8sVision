# K8sVision

K8sVision 是一款现代化的 Kubernetes 可视化监控与资源管理平台，前后端分离，支持本地和生产环境一键部署，界面美观，易于二次开发。

---

## 快速开始

### 1. 本地开发

```bash
git clone https://github.com/nick0323/K8sVision.git
cd K8sVision
docker-compose up -d
```
- 前端：http://localhost
- 后端API：http://localhost:8080/api
- Swagger文档：http://localhost:8080/swagger/index.html

### 2. Kubernetes 集群部署

1. 构建并推送镜像
2. 修改 `k8s/*.yaml` 中镜像名
3. 部署：
   ```bash
   kubectl apply -f k8s/
   ```
- 通过 Ingress 域名访问（如 `k8svision.local`）

---

## 配置与环境变量

所有配置均可通过环境变量或 `config.yaml` 设置，**环境变量优先生效**。

| 变量名               | 说明                       | 默认值                  |
|----------------------|----------------------------|-------------------------|
| SERVER_PORT          | 服务监听端口               | 8080                    |
| LOGIN_USERNAME       | 登录用户名                 | admin                   |
| LOGIN_PASSWORD       | 登录密码                   | 123456                  |
| JWT_SECRET           | JWT 签发密钥               | k8svision-secret-key    |
| KUBECONFIG           | kubeconfig 路径            | ""                      |
| SWAGGER_ENABLE       | 是否启用Swagger文档        | false                   |
| LOGIN_MAX_FAIL       | 登录失败最大次数           | 5                       |
| LOGIN_LOCK_MINUTES   | 登录失败锁定时长（分钟）   | 10                      |

**补充说明：**
- `SWAGGER_ENABLE`：设置为`true`时启用Swagger接口文档。
- `LOGIN_MAX_FAIL`：连续失败达到该值后锁定账号。
- `LOGIN_LOCK_MINUTES`：账号锁定后解锁时间。

---

## 主要功能

- 支持 Node、Pod、Deployment、StatefulSet、DaemonSet、Service、Ingress、Job、CronJob、Event、Namespace 等核心资源的可视化管理
- 集群资源总览与健康检查
- 资源用量可视化（CPU/内存等）
- 实时事件追踪与检索
- 分页、检索、状态展示
- JWT 认证、Swagger 文档
- 多环境部署（本地、K8s）

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
├── Dockerfile            # 后端镜像构建
├── frontend/             # React 前端应用
├── k8s/                  # K8s 部署 YAML
├── docker-compose.yml    # 一键本地部署
└── README.md             # 项目说明
```

---

## 开发指南

### 后端开发

```bash
go mod tidy
go run main.go
```
- 默认端口 8080
- 依赖：Gin、zap、viper、swaggo

### 前端开发

```bash
cd frontend
npm install
npm run dev
```
- 默认端口 5173

---

## 接口与安全

- 登录接口：`/api/login`，JWT 认证，防爆破
- 健康检查：`/healthz`，GET
- 资源接口：全部以 `/api/` 开头，需 Bearer Token
- Swagger 文档：`/swagger/index.html`

---

## FAQ

**Q: 前端访问 404？**  
A: 检查 nginx.conf 配置 `try_files $uri $uri/ /index.html;`。

**Q: API 代理失败？**  
A: 检查 nginx.conf 的 `proxy_pass` 是否正确。

**Q: 如何自定义配置？**  
A: 支持通过环境变量或 config.yaml 注入。

---

## 贡献与 License

- 欢迎 PR、Issue、建议与二次开发！
- License: MIT

---

如需详细接口说明、二次开发指引或架构解读，请查阅源码注释与 Swagger 文档。
