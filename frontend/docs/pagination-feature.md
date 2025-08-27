# 分页功能使用指南

## 概述

K8sVision 前端现在支持用户自定义选择每页显示的行数，提供了更灵活的数据浏览体验。

## 功能特性

- **动态每页行数选择**：用户可以选择每页显示 10、20、50 或 100 行数据
- **智能分页重置**：当改变每页行数时，自动重置到第一页
- **统一的分页Hook**：提供 `usePagination` Hook 简化分页逻辑
- **响应式设计**：分页组件支持固定定位和底部定位

## 使用方法

### 1. 使用 usePagination Hook（推荐）

```jsx
import { usePagination } from './hooks/usePagination';

export default function MyPage() {
  const {
    page,
    pageSize,
    handlePageChange,
    handlePageSizeChange,
    resetPagination,
    pageSizeOptions
  } = usePagination();

  // 在API调用中使用
  const fetchData = useCallback(() => {
    const params = new URLSearchParams({
      limit: pageSize.toString(),
      offset: ((page-1)*pageSize).toString(),
    });
    // ... API调用逻辑
  }, [page, pageSize]);

  return (
    <div>
      {/* 表格组件 */}
      <CommonTable
        data={data}
        pageSize={pageSize}
        currentPage={page}
        onPageChange={handlePageChange}
        // ... 其他属性
      />
      
      {/* 分页组件 */}
      <Pagination
        currentPage={page}
        total={total}
        pageSize={pageSize}
        onPageChange={handlePageChange}
        onPageSizeChange={handlePageSizeChange}
        pageSizeOptions={pageSizeOptions}
        fixedBottom={true}
      />
    </div>
  );
}
```

### 2. 手动管理分页状态

```jsx
import { PAGE_SIZE, PAGE_SIZE_OPTIONS } from './constants';

export default function MyPage() {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(PAGE_SIZE);

  const handlePageSizeChange = (newPageSize) => {
    setPageSize(newPageSize);
    setPage(1); // 重置到第一页
  };

  return (
    <Pagination
      currentPage={page}
      total={total}
      pageSize={pageSize}
      onPageChange={setPage}
      onPageSizeChange={handlePageSizeChange}
      pageSizeOptions={PAGE_SIZE_OPTIONS}
    />
  );
}
```

## 配置选项

### 每页行数选项

默认提供以下选项：
- 10 行/页
- 20 行/页（默认）
- 50 行/页
- 100 行/页

可以通过修改 `constants.js` 中的 `PAGE_SIZE_OPTIONS` 来自定义：

```js
export const PAGE_SIZE_OPTIONS = [5, 10, 20, 50, 100, 200];
```

### 分页组件属性

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `currentPage` | number | 1 | 当前页码 |
| `total` | number | 0 | 总数据量 |
| `pageSize` | number | 20 | 每页行数 |
| `onPageChange` | function | - | 页码变化回调 |
| `onPageSizeChange` | function | - | 每页行数变化回调 |
| `pageSizeOptions` | array | [10,20,50,100] | 可选的每页行数 |
| `fixed` | boolean | false | 是否固定在表格底部 |
| `fixedBottom` | boolean | false | 是否固定在页面底部 |

## 样式定制

分页组件使用以下CSS类名，可以通过CSS自定义样式：

```css
/* 分页容器 */
.table-pagination-area

/* 左侧信息区域 */
.pagination-total
.pagination-separator
.page-size-selector
.page-size-select

/* 右侧控制区域 */
.pagination-controls
.pagination-info
.pagination-btn
```

## 最佳实践

1. **使用Hook**：推荐使用 `usePagination` Hook 来管理分页状态
2. **重置分页**：当筛选条件改变时（如namespace切换），使用 `resetPagination()` 重置到第一页
3. **API集成**：确保后端API支持 `limit` 和 `offset` 参数
4. **性能考虑**：避免在每页行数变化时重复请求相同的数据

## 迁移指南

### 从旧版本迁移

1. 导入新的Hook：
   ```jsx
   import { usePagination } from './hooks/usePagination';
   ```

2. 替换分页状态管理：
   ```jsx
   // 旧版本
   const [page, setPage] = useState(1);
   const pageSize = PAGE_SIZE;
   
   // 新版本
   const { page, pageSize, handlePageChange } = usePagination();
   ```

3. 更新Pagination组件：
   ```jsx
   // 旧版本
   <Pagination pageSize={pageSize} onPageChange={setPage} />
   
   // 新版本
   <Pagination 
     pageSize={pageSize} 
     onPageChange={handlePageChange}
     onPageSizeChange={handlePageSizeChange}
     pageSizeOptions={pageSizeOptions}
   />
   ```

## 故障排除

### 常见问题

1. **每页行数选择器不显示**
   - 检查是否正确传递了 `onPageSizeChange` 和 `pageSizeOptions` 属性
   - 确保 `total` 大于 `pageSize`

2. **分页计算错误**
   - 验证 `pageSize` 和 `page` 状态是否正确更新
   - 检查API返回的 `total` 值是否正确

3. **样式问题**
   - 确保引入了相关的CSS样式
   - 检查CSS类名是否正确应用

## 更新日志

- **v1.0.0**: 初始版本，支持固定每页行数
- **v2.0.0**: 新增动态每页行数选择功能
- **v2.1.0**: 新增 `usePagination` Hook 简化使用
