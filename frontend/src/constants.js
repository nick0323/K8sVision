// 常量配置
export const PAGE_SIZE = 20;
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
      { key: 'pods', label: 'Pods', icon: 'FaCubes' },
      { key: 'deployments', label: 'Deployments', icon: 'FaLayerGroup' },
      { key: 'statefulsets', label: 'StatefulSets', icon: 'FaDatabase' },
      { key: 'daemonsets', label: 'DaemonSets', icon: 'FaShieldAlt' },
      { key: 'jobs', label: 'Jobs', icon: 'FaTasks' },
      { key: 'cronjobs', label: 'CronJobs', icon: 'FaRegClock' },
    ],
  },
  {
    group: 'Traffic',
    items: [
      { key: 'ingress', label: 'Ingresses', icon: 'FaNetworkWired' },
      { key: 'services', label: 'Services', icon: 'FaProjectDiagram' },
    ],
  },
  {
    group: 'Storage',
    items: [
      { key: 'pvcs', label: 'PVCs', icon: 'FaBoxOpen' },
      { key: 'pvs', label: 'PVs', icon: 'FaHdd' },
      { key: 'storageclasses', label: 'StorageClasses', icon: 'FaBoxes' },
    ],
  },
  {
    group: 'Config',
    items: [
      { key: 'configmaps', label: 'ConfigMaps', icon: 'FaCog' },
      { key: 'secrets', label: 'Secrets', icon: 'FaKey' },
    ],
  },
  {
    group: 'Others',
    items: [
      { key: 'namespaces', label: 'Namespaces', icon: 'LuSquareDashed' },
      { key: 'nodes', label: 'Nodes', icon: 'FaDesktop' },
      { key: 'events', label: 'Events', icon: 'FaBell' },
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
  // 存储资源详情API
  pvcs: '/api/pvcs',
  pvs: '/api/pvs',
  storageclasses: '/api/storageclasses',
  // 配置资源详情API
  configmaps: '/api/configmaps',
  secrets: '/api/secrets',
}; 