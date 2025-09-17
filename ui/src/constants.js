// 常量配置
export const PAGE_SIZE = 20;
export const PAGE_SIZE_OPTIONS = [10, 20, 50, 100]; // 新增：可选的每页行数选项
export const EMPTY_TEXT = 'No data';
export const LOADING_TEXT = 'Loading...';
export const ERROR_TEXT = 'Load failed:';
export const SEARCH_PLACEHOLDER = 'Search Resources...';

// 菜单分组配置
export const MENU_LIST = [
  {
    group: 'Overview',
    items: [
      { key: 'overview', label: 'Overview', icon: 'FaChartPie' },
    ],
  },
  {
    group: 'Workloads',
    items: [
      { key: 'pods', label: 'Pod', icon: 'FaCubes' },
      { key: 'deployments', label: 'Deployment', icon: 'FaLayerGroup' },
      { key: 'statefulsets', label: 'StatefulSet', icon: 'FaDatabase' },
      { key: 'daemonsets', label: 'DaemonSet', icon: 'FaShieldAlt' },
      { key: 'jobs', label: 'Job', icon: 'FaTasks' },
      { key: 'cronjobs', label: 'CronJob', icon: 'FaRegClock' },
    ],
  },
  {
    group: 'Traffic',
    items: [
      { key: 'ingress', label: 'Ingress', icon: 'FaNetworkWired' },
      { key: 'services', label: 'Service', icon: 'FaProjectDiagram' },
    ],
  },
  {
    group: 'Storage',
    items: [
      { key: 'pvcs', label: 'PVC', icon: 'FaBoxOpen' },
      { key: 'pvs', label: 'PV', icon: 'FaHdd' },
      { key: 'storageclasses', label: 'StorageClass', icon: 'FaBoxes' },
    ],
  },
  {
    group: 'Config',
    items: [
      { key: 'configmaps', label: 'ConfigMap', icon: 'FaCog' },
      { key: 'secrets', label: 'Secret', icon: 'FaKey' },
    ],
  },
  {
    group: 'Others',
    items: [
      { key: 'namespaces', label: 'Namespace', icon: 'LuSquareDashed' },
      { key: 'nodes', label: 'Node', icon: 'FaDesktop' },
      { key: 'events', label: 'Event', icon: 'FaBell' },
    ],
  },
];

// API 映射
export const API_MAP = {
  overview: '/api/overview',
  namespaces: '/api/namespaces',
  nodes: '/api/nodes',
  pods: '/api/pods',
  deployments: '/api/deployments',
  statefulsets: '/api/statefulsets',
  daemonsets: '/api/daemonsets',
  cronjobs: '/api/cronjobs',
  jobs: '/api/jobs',
  ingress: '/api/ingress',
  services: '/api/services',
  events: '/api/events',
  // 存储资源API
  pvcs: '/api/pvcs',
  pvs: '/api/pvs',
  storageclasses: '/api/storageclasses',
  // 配置资源API
  configmaps: '/api/configmaps',
  secrets: '/api/secrets',
};

// 详情API映射
export const DETAIL_API_MAP = {
  namespaces: '/api/namespaces',
  pods: '/api/pods',
  nodes: '/api/nodes',
  services: '/api/services',
  ingress: '/api/ingress',
  // 存储资源详情API
  pvcs: '/api/pvcs',
  pvs: '/api/pvs',
  storageclasses: '/api/storageclasses',
  // 配置资源详情API
  configmaps: '/api/configmaps',
  secrets: '/api/secrets',
}; 