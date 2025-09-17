# 前端模块文档

K8sVision前端是一个基于React的现代化Web应用，提供直观的Kubernetes集群管理界面。

## 📁 项目结构

```
frontend/
├── README.md                    # 前端文档
├── package.json                 # 依赖配置
├── vite.config.js               # Vite构建配置
├── index.html                   # HTML入口
├── nginx.conf                   # Nginx配置
├── Dockerfile                   # Docker构建文件
└── src/
    ├── App.jsx                  # 主应用组件
    ├── App.css                  # 主样式文件
    ├── main.jsx                 # 应用入口
    ├── index.css                # 全局样式
    ├── LoginPage.jsx            # 登录页面
    ├── LoginPage.css            # 登录页面样式
    ├── OverviewPage.jsx         # 概览页面
    ├── pages.jsx                # 页面组件定义
    ├── constants.js             # 常量定义
    ├── constants/
    │   └── pageConfigs.js       # 页面配置
    ├── components/              # 通用组件
    │   ├── ErrorBoundary.jsx    # 错误边界
    │   ├── LoadingSpinner.jsx   # 加载动画
    │   ├── LoadingSpinner.css   # 加载动画样式
    │   ├── NamespaceSelect.jsx  # 命名空间选择器
    │   ├── PageHeader.jsx       # 页面头部
    │   ├── PodDetailPage.jsx    # Pod详情页面
    │   ├── PodDetailPage.css    # Pod详情页面样式
    │   ├── README_PodDetail.md  # Pod详情页面文档
    │   ├── RefreshButton.jsx    # 刷新按钮
    │   ├── ResourceDetailComponents.jsx # 资源详情组件
    │   ├── ResourceDetailModal.css # 资源详情模态框样式
    │   ├── ResourceDetailModal.jsx # 资源详情模态框
    │   ├── ResourcePage.css     # 资源页面样式
    │   ├── ResourcePage.jsx     # 资源页面
    │   ├── ScrollbarHide.css    # 滚动条隐藏样式
    │   ├── SearchInput.jsx      # 搜索输入框
    │   └── StatusRenderer.jsx   # 状态渲染器
    ├── hooks/                   # 自定义Hooks
    │   ├── useGlobalSearch.js   # 全局搜索Hook
    │   ├── useOptimizedState.js # 优化状态Hook
    │   ├── usePagination.js     # 分页Hook
    │   └── useSimpleSearch.js   # 简单搜索Hook
    └── utils/                   # 工具函数
        ├── apiUtils.js          # API工具
        ├── authUtils.js         # 认证工具
        ├── commonRenderers.jsx  # 通用渲染器
        ├── dataUtils.js         # 数据处理工具
        ├── dateUtils.js         # 日期工具
        ├── hooks.js             # Hook工具
        ├── index.js             # 工具入口
        └── tableUtils.js        # 表格工具
```

## 🚀 技术栈

### 核心框架
- **React 18**: 用户界面库
- **Vite**: 构建工具和开发服务器
- **React Router**: 客户端路由

### UI组件
- **React Icons**: 图标库
- **自定义组件**: 专门为K8s管理设计的组件

### 状态管理
- **React Hooks**: 状态管理
- **Context API**: 全局状态
- **自定义Hooks**: 业务逻辑封装

### 样式方案
- **CSS Modules**: 组件样式隔离
- **CSS Variables**: 主题变量
- **响应式设计**: 移动端适配

## 🎨 设计系统

### 颜色方案
```css
:root {
  --primary-color: #2563eb;
  --secondary-color: #64748b;
  --success-color: #10b981;
  --warning-color: #f59e0b;
  --error-color: #ef4444;
  --info-color: #3b82f6;
  --background-color: #f8fafc;
  --surface-color: #ffffff;
  --text-primary: #1e293b;
  --text-secondary: #64748b;
  --border-color: #e2e8f0;
}
```

### 组件规范
- **一致性**: 统一的视觉风格
- **可访问性**: 支持键盘导航和屏幕阅读器
- **响应式**: 适配不同屏幕尺寸
- **性能**: 优化渲染性能

## 🔧 核心组件

### 1. 主应用组件 (App.jsx)
应用的主入口组件：
- 路由管理
- 全局状态管理
- 错误边界处理
- 认证状态管理

**主要功能：**
- 侧边栏导航
- 页面路由
- 用户认证
- 主题切换

### 2. 资源页面 (ResourcePage.jsx)
统一的资源列表页面：
- 数据获取和显示
- 分页和搜索
- 排序和过滤
- 操作按钮

**特性：**
- 懒加载
- 虚拟滚动
- 实时更新
- 批量操作

### 3. 资源详情页面 (PodDetailPage.jsx)
详细的资源信息展示：
- 分层信息架构
- 标签页式展示
- 实时状态更新
- 操作功能

### 4. 通用组件
#### LoadingSpinner
加载动画组件：
- 多种加载样式
- 可配置大小
- 自定义文本

#### StatusRenderer
状态渲染组件：
- 状态颜色编码
- 图标显示
- 工具提示

#### SearchInput
搜索输入组件：
- 实时搜索
- 搜索建议
- 清空功能

## 🎯 自定义Hooks

### useOptimizedState
优化的状态管理Hook：
```javascript
const [state, setState] = useOptimizedState(initialState);
```
- 防止不必要的重渲染
- 状态更新优化
- 内存使用优化

### usePagination
分页管理Hook：
```javascript
const {
  currentPage,
  pageSize,
  totalPages,
  goToPage,
  nextPage,
  prevPage
} = usePagination(data, initialPageSize);
```

### useGlobalSearch
全局搜索Hook：
```javascript
const {
  searchTerm,
  setSearchTerm,
  filteredData,
  isSearching
} = useGlobalSearch(data, searchFields);
```

### useSimpleSearch
简单搜索Hook：
```javascript
const {
  searchTerm,
  setSearchTerm,
  filteredData
} = useSimpleSearch(data, searchKey);
```

## 🛠️ 工具函数

### API工具 (apiUtils.js)
- HTTP请求封装
- 错误处理
- 请求拦截器
- 响应格式化

### 认证工具 (authUtils.js)
- 令牌管理
- 登录状态检查
- 自动登出
- 权限验证

### 数据处理工具 (dataUtils.js)
- 数据转换
- 排序和过滤
- 分页处理
- 搜索功能

### 日期工具 (dateUtils.js)
- 日期格式化
- 相对时间计算
- 时区处理
- 日期验证

## 📱 响应式设计

### 断点设置
```css
/* 移动端 */
@media (max-width: 768px) { }

/* 平板端 */
@media (min-width: 769px) and (max-width: 1024px) { }

/* 桌面端 */
@media (min-width: 1025px) { }
```

### 布局适配
- 侧边栏折叠
- 表格横向滚动
- 按钮组重排
- 字体大小调整

## 🚀 性能优化

### 代码分割
- 路由级别的代码分割
- 组件懒加载
- 动态导入

### 渲染优化
- React.memo
- useMemo
- useCallback
- 虚拟滚动

### 资源优化
- 图片懒加载
- 字体优化
- CSS优化
- 打包优化

## 🧪 测试

### 测试工具
- Jest: 单元测试框架
- React Testing Library: 组件测试
- Cypress: 端到端测试

### 测试覆盖
- 组件渲染测试
- 用户交互测试
- API集成测试
- 错误处理测试

## 📦 构建和部署

### 开发环境
```bash
npm install
npm run dev
```

### 生产构建
```bash
npm run build
```

### Docker部署
```dockerfile
FROM nginx:alpine
COPY dist/ /usr/share/nginx/html/
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
```

## 🔧 开发指南

### 添加新页面
1. 在`pages.jsx`中定义页面组件
2. 在`constants/pageConfigs.js`中添加页面配置
3. 在`App.jsx`中添加路由
4. 实现页面组件

### 添加新组件
1. 在`components/`目录下创建组件文件
2. 实现组件逻辑和样式
3. 添加必要的测试
4. 更新文档

### 样式规范
- 使用CSS Modules
- 遵循BEM命名规范
- 使用CSS变量
- 保持样式简洁

## 📝 最佳实践

1. **组件设计**
   - 单一职责原则
   - 可复用性
   - 可测试性
   - 可访问性

2. **状态管理**
   - 最小化状态
   - 状态提升
   - 避免不必要的重渲染
   - 使用适当的Hook

3. **性能考虑**
   - 懒加载
   - 虚拟化
   - 防抖和节流
   - 内存管理

4. **代码质量**
   - 类型检查
   - 代码规范
   - 注释文档
   - 单元测试

## 🔍 故障排查

### 常见问题
1. **页面空白**
   - 检查控制台错误
   - 验证API连接
   - 检查认证状态

2. **性能问题**
   - 使用React DevTools
   - 检查重渲染
   - 分析包大小

3. **样式问题**
   - 检查CSS加载
   - 验证选择器
   - 查看浏览器兼容性

### 调试工具
- React DevTools
- Chrome DevTools
- Vite DevTools
- 网络面板
