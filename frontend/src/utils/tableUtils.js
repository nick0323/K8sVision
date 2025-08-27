/**
 * 通用表格工具函数
 * 消除重复的表格配置和逻辑
 */

import { 
  createNameRenderer, 
  createNamespaceRenderer, 
  createStatusRenderer, 
  createLabelsRenderer,
  createTimeRenderer,
  createNumberRenderer,
  createBooleanRenderer,
  createArrayRenderer,
  createUsageRenderer
} from './commonRenderers.jsx';

// 通用列类型定义
export const COLUMN_TYPES = {
  NAME: 'name',
  NAMESPACE: 'namespace',
  STATUS: 'status',
  LABELS: 'labels',
  TIME: 'time',
  NUMBER: 'number',
  BOOLEAN: 'boolean',
  ARRAY: 'array',
  USAGE: 'usage',
  TEXT: 'text'
};

// 列配置生成器
export const createColumn = (title, dataIndex, type = COLUMN_TYPES.TEXT, options = {}) => {
  const {
    width = null,
    sortable = false,
    filterable = false,
    render = null,
    ...rest
  } = options;

  // 根据类型自动生成渲染器
  let autoRender = render;
  if (!render) {
    switch (type) {
      case COLUMN_TYPES.NAME:
        // 名称列需要外部传入点击处理器
        autoRender = null;
        break;
      case COLUMN_TYPES.NAMESPACE:
        autoRender = createNamespaceRenderer();
        break;
      case COLUMN_TYPES.STATUS:
        autoRender = createStatusRenderer(options.statusMap);
        break;
      case COLUMN_TYPES.LABELS:
        autoRender = createLabelsRenderer();
        break;
      case COLUMN_TYPES.TIME:
        autoRender = createTimeRenderer();
        break;
      case COLUMN_TYPES.NUMBER:
        autoRender = createNumberRenderer(options.suffix, options.formatter);
        break;
      case COLUMN_TYPES.BOOLEAN:
        autoRender = createBooleanRenderer();
        break;
      case COLUMN_TYPES.ARRAY:
        autoRender = createArrayRenderer(options.separator, options.maxItems);
        break;
      case COLUMN_TYPES.USAGE:
        autoRender = createUsageRenderer();
        break;
      default:
        autoRender = null;
    }
  }

  return {
    title,
    dataIndex,
    render: autoRender,
    width,
    sortable,
    filterable,
    ...rest
  };
};

// 预定义的列配置
export const PREDEFINED_COLUMNS = {
  // 基础列
  name: (options = {}) => createColumn('Name', 'name', COLUMN_TYPES.NAME, options),
  namespace: (options = {}) => createColumn('Namespace', 'namespace', COLUMN_TYPES.NAMESPACE, options),
  status: (options = {}) => createColumn('Status', 'status', COLUMN_TYPES.STATUS, options),
  age: (options = {}) => createColumn('Age', 'age', COLUMN_TYPES.TIME, options),
  created: (options = {}) => createColumn('Created', 'created', COLUMN_TYPES.TIME, options),
  
  // 工作负载相关列
  replicas: (options = {}) => createColumn('Replicas', 'replicas', COLUMN_TYPES.NUMBER, options),
  availableReplicas: (options = {}) => createColumn('Available', 'availableReplicas', COLUMN_TYPES.NUMBER, options),
  desiredReplicas: (options = {}) => createColumn('Desired', 'desiredReplicas', COLUMN_TYPES.NUMBER, options),
  readyReplicas: (options = {}) => createColumn('Ready', 'readyReplicas', COLUMN_TYPES.NUMBER, options),
  
  // 资源使用列
  cpuUsage: (options = {}) => createColumn('CPU Usage', 'cpuUsage', COLUMN_TYPES.USAGE, options),
  memoryUsage: (options = {}) => createColumn('Memory Usage', 'memoryUsage', COLUMN_TYPES.USAGE, options),
  
  // 网络相关列
  podIP: (options = {}) => createColumn('Pod IP', 'podIP', COLUMN_TYPES.TEXT, options),
  clusterIP: (options = {}) => createColumn('Cluster IP', 'clusterIP', COLUMN_TYPES.TEXT, options),
  nodeName: (options = {}) => createColumn('Node', 'nodeName', COLUMN_TYPES.TEXT, options),
  
  // 存储相关列
  capacity: (options = {}) => createColumn('Capacity', 'capacity', COLUMN_TYPES.TEXT, options),
  accessMode: (options = {}) => createColumn('Access Mode', 'accessMode', COLUMN_TYPES.TEXT, options),
  storageClass: (options = {}) => createColumn('Storage Class', 'storageClass', COLUMN_TYPES.TEXT, options),
  
  // 配置相关列
  type: (options = {}) => createColumn('Type', 'type', COLUMN_TYPES.TEXT, options),
  data: (options = {}) => createColumn('Data', 'data', COLUMN_TYPES.LABELS, options),
  labels: (options = {}) => createColumn('Labels', 'labels', COLUMN_TYPES.LABELS, options),
  
  // 事件相关列
  reason: (options = {}) => createColumn('Reason', 'reason', COLUMN_TYPES.TEXT, options),
  message: (options = {}) => createColumn('Message', 'message', COLUMN_TYPES.TEXT, options),
  lastTimestamp: (options = {}) => createColumn('Last Timestamp', 'lastTimestamp', COLUMN_TYPES.TIME, options),
  
  // 其他列
  version: (options = {}) => createColumn('Version', 'version', COLUMN_TYPES.TEXT, options),
  internalIP: (options = {}) => createColumn('Internal IP', 'internalIP', COLUMN_TYPES.TEXT, options),
  roles: (options = {}) => createColumn('Roles', 'roles', COLUMN_TYPES.ARRAY, { separator: ', ', maxItems: 2, ...options }),
  ports: (options = {}) => createColumn('Ports', 'ports', COLUMN_TYPES.ARRAY, { separator: ', ', maxItems: 2, ...options }),
  hosts: (options = {}) => createColumn('Hosts', 'hosts', COLUMN_TYPES.ARRAY, { separator: ', ', maxItems: 2, ...options }),
  
  // 布尔值列
  allowVolumeExpansion: (options = {}) => createColumn('Allow Volume Expansion', 'allowVolumeExpansion', COLUMN_TYPES.BOOLEAN, options),
  
  // 自定义列
  custom: (title, dataIndex, options = {}) => createColumn(title, dataIndex, COLUMN_TYPES.TEXT, options)
};

// 表格配置生成器
export const createTableConfig = (columns, options = {}) => {
  const {
    sortable = false,
    filterable = false,
    selectable = false,
    expandable = false,
    ...rest
  } = options;

  return {
    columns,
    sortable,
    filterable,
    selectable,
    expandable,
    ...rest
  };
};

// 列排序工具
export const sortColumns = (columns, order = []) => {
  if (!order.length) return columns;
  
  const sortedColumns = [];
  const remainingColumns = [...columns];
  
  // 按照指定顺序排列列
  order.forEach(key => {
    const index = remainingColumns.findIndex(col => col.dataIndex === key);
    if (index !== -1) {
      sortedColumns.push(remainingColumns[index]);
      remainingColumns.splice(index, 1);
    }
  });
  
  // 添加剩余的列
  return [...sortedColumns, ...remainingColumns];
};

// 列过滤工具
export const filterColumns = (columns, visibleColumns = []) => {
  if (!visibleColumns.length) return columns;
  
  return columns.filter(col => visibleColumns.includes(col.dataIndex));
};

// 列宽度计算工具
export const calculateColumnWidths = (columns, containerWidth, minWidth = 100) => {
  const totalColumns = columns.length;
  const availableWidth = containerWidth - (totalColumns - 1) * 8; // 减去边框宽度
  
  if (availableWidth <= 0) {
    return columns.map(() => minWidth);
  }
  
  // 计算每列的基础宽度
  const baseWidth = Math.max(minWidth, Math.floor(availableWidth / totalColumns));
  
  // 为有特殊宽度要求的列分配空间
  return columns.map(col => {
    if (col.width) {
      return Math.max(minWidth, col.width);
    }
    return baseWidth;
  });
};

