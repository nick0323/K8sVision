# 常见问题解答 (FAQ)

本文档收集了 K8sVision 使用过程中的常见问题和解决方案。

## 🔐 认证相关

### Q1: 登录失败，提示"用户名或密码错误"
**问题描述**: 无法使用默认用户名密码登录系统

**可能原因**:
1. 环境变量配置错误
2. 配置文件中的用户名密码不正确
3. 登录失败次数过多，账号被锁定

**解决方案**:
```bash
# 1. 检查环境变量
echo $LOGIN_USERNAME
echo $LOGIN_PASSWORD

# 2. 检查配置文件
cat config.yaml | grep -A 5 auth

# 3. 重置登录失败计数
docker-compose restart backend

# 4. 使用默认凭据
用户名: admin
密码: 12345678
```

### Q2: JWT Token 过期
**问题描述**: 使用一段时间后提示 Token 过期

**解决方案**:
```bash
# 重新登录获取新的 Token
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"12345678"}'
```

**预防措施**:
- Token 默认 24 小时过期
- 建议在 Token 过期前主动刷新
- 可以调整 `JWT_EXPIRATION` 环境变量

## 🌐 网络连接

### Q3: 无法连接 Kubernetes 集群
**问题描述**: 后端无法连接到 Kubernetes 集群

**可能原因**:
1. kubeconfig 配置错误
2. 集群访问权限不足
3. 网络连接问题

**解决方案**:
```bash
# 1. 检查 kubectl 配置
kubectl cluster-info

# 2. 检查 kubeconfig 文件权限
ls -la ~/.kube/config
chmod 600 ~/.kube/config

# 3. 测试集群连接
kubectl get nodes

# 4. 检查 Docker 挂载
docker-compose logs backend | grep kubeconfig
```

### Q4: 前端无法访问后端 API
**问题描述**: 前端显示 API 连接错误

**可能原因**:
1. 后端服务未启动
2. 端口配置错误
3. CORS 配置问题

**解决方案**:
```bash
# 1. 检查后端服务状态
docker-compose ps
curl http://localhost:8080/healthz

# 2. 检查端口占用
netstat -tlnp | grep :8080

# 3. 查看后端日志
docker-compose logs backend

# 4. 检查前端配置
# 确认前端 API 地址配置正确
```

## 🐳 Docker 相关

### Q5: Docker 容器启动失败
**问题描述**: `docker-compose up` 失败

**可能原因**:
1. 端口被占用
2. 镜像拉取失败
3. 权限问题

**解决方案**:
```bash
# 1. 检查端口占用
netstat -tlnp | grep :8080
netstat -tlnp | grep :80

# 2. 清理容器和镜像
docker-compose down
docker system prune -f

# 3. 重新拉取镜像
docker-compose pull

# 4. 重新启动
docker-compose up -d

# 5. 查看详细错误
docker-compose logs
```

### Q6: 容器内存不足
**问题描述**: 容器因内存不足被杀死

**解决方案**:
```bash
# 1. 增加 Docker 内存限制
# 在 Docker Desktop 设置中增加内存限制

# 2. 优化容器资源使用
docker-compose down
docker system prune -f

# 3. 检查系统内存
free -h

# 4. 调整容器资源限制
# 在 docker-compose.yml 中添加资源限制
```

## ☸️ Kubernetes 相关

### Q7: 无法查看某些资源
**问题描述**: 某些 Kubernetes 资源无法在界面中显示

**可能原因**:
1. 权限不足
2. 资源不存在
3. API 版本不兼容

**解决方案**:
```bash
# 1. 检查权限
kubectl auth can-i list pods --all-namespaces
kubectl auth can-i list nodes

# 2. 检查资源是否存在
kubectl get pods --all-namespaces
kubectl get nodes

# 3. 检查 API 版本
kubectl api-resources

# 4. 查看后端日志
docker-compose logs backend | grep -i error
```

### Q8: 资源数据不更新
**问题描述**: 界面显示的资源数据不是最新的

**可能原因**:
1. 缓存问题
2. 集群数据同步延迟
3. 前端缓存问题

**解决方案**:
```bash
# 1. 清除后端缓存
curl -X POST http://localhost:8080/api/cache/clear

# 2. 刷新浏览器缓存
# 按 Ctrl+F5 强制刷新

# 3. 检查缓存配置
cat config.yaml | grep -A 10 cache

# 4. 重启服务
docker-compose restart
```

## 📊 性能问题

### Q9: 页面加载缓慢
**问题描述**: 界面加载速度很慢

**可能原因**:
1. 集群资源过多
2. 网络延迟
3. 缓存配置不当

**解决方案**:
```bash
# 1. 检查集群资源数量
kubectl get pods --all-namespaces | wc -l
kubectl get nodes | wc -l

# 2. 调整缓存配置
# 增加缓存 TTL 和大小

# 3. 检查网络延迟
ping <kubernetes-api-server>

# 4. 查看性能指标
curl http://localhost:8080/metrics
```

### Q10: API 响应超时
**问题描述**: API 请求超时

**解决方案**:
```bash
# 1. 检查 Kubernetes API 响应时间
time kubectl get nodes

# 2. 调整超时配置
# 在 config.yaml 中增加 timeout 配置

# 3. 检查网络连接
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8080/healthz

# 4. 查看后端日志
docker-compose logs backend | grep timeout
```

## 🔧 配置问题

### Q11: 配置文件不生效
**问题描述**: 修改配置文件后不生效

**解决方案**:
```bash
# 1. 检查配置文件路径
docker-compose exec backend ls -la /app/config.yaml

# 2. 重启服务
docker-compose restart backend

# 3. 检查环境变量覆盖
docker-compose exec backend env | grep K8SVISION

# 4. 查看配置加载日志
docker-compose logs backend | grep config
```

### Q12: 环境变量不生效
**问题描述**: 设置的环境变量没有生效

**解决方案**:
```bash
# 1. 检查环境变量格式
# 确保变量名正确，如 LOGIN_USERNAME 而不是 LOGIN_USER

# 2. 重启容器
docker-compose down
docker-compose up -d

# 3. 验证环境变量
docker-compose exec backend env | grep LOGIN

# 4. 检查 docker-compose.yml 配置
cat docker-compose.yml | grep -A 5 environment
```

## 📝 日志问题

### Q13: 日志文件过大
**问题描述**: 日志文件占用磁盘空间过多

**解决方案**:
```bash
# 1. 查看日志大小
du -sh /var/lib/docker/containers/*/*.log

# 2. 清理日志
docker system prune -f

# 3. 配置日志轮转
# 在 docker-compose.yml 中配置日志驱动

# 4. 调整日志级别
# 在 config.yaml 中设置 log.level
```

### Q14: 无法查看详细日志
**问题描述**: 日志信息不够详细，难以排查问题

**解决方案**:
```bash
# 1. 调整日志级别为 debug
export LOG_LEVEL=debug
docker-compose restart backend

# 2. 查看完整日志
docker-compose logs -f --tail=100 backend

# 3. 启用 Swagger 文档
export SWAGGER_ENABLE=true
docker-compose restart backend

# 4. 查看性能指标
curl http://localhost:8080/metrics
```

## 🚀 部署问题

### Q15: 生产环境部署失败
**问题描述**: 在生产环境中部署失败

**可能原因**:
1. 资源限制
2. 网络配置
3. 权限问题

**解决方案**:
```bash
# 1. 检查资源限制
kubectl describe pod <pod-name> -n k8svision

# 2. 检查事件
kubectl get events -n k8svision --sort-by='.lastTimestamp'

# 3. 检查网络策略
kubectl get networkpolicies -n k8svision

# 4. 检查 RBAC 权限
kubectl auth can-i --list -n k8svision
```

### Q16: 升级版本后出现问题
**问题描述**: 升级 K8sVision 版本后功能异常

**解决方案**:
```bash
# 1. 备份配置
cp config.yaml config.yaml.backup

# 2. 查看升级日志
docker-compose logs backend | grep -i upgrade

# 3. 回滚到旧版本
docker-compose down
git checkout <previous-version>
docker-compose up -d

# 4. 检查兼容性
# 查看版本更新说明
```

## 📞 获取帮助

如果以上解决方案无法解决您的问题，请：

1. **查看详细日志**:
   ```bash
   docker-compose logs -f backend
   ```

2. **检查系统状态**:
   ```bash
   # 检查服务状态
   docker-compose ps
   
   # 检查资源使用
   docker stats
   
   # 检查网络连接
   curl -v http://localhost:8080/healthz
   ```

3. **收集诊断信息**:
   ```bash
   # 收集系统信息
   docker-compose exec backend uname -a
   docker-compose exec backend cat /etc/os-release
   
   # 收集配置信息
   docker-compose exec backend env
   docker-compose exec backend cat /app/config.yaml
   ```

4. **提交 Issue**:
   - 访问 [GitHub Issues](https://github.com/nick0323/K8sVision/issues)
   - 提供详细的错误信息和复现步骤
   - 附上相关的日志和配置信息

## 📚 相关文档

- [快速安装](../quickstart.md)
- [配置说明](../configuration.md)
- [错误代码](./error-codes.md)
- [日志分析](./logs.md)
- [API 文档](../api/README.md)

---

**最后更新**: 2024年12月 