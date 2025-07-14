// 常量配置
export const PAGE_SIZE = 13;
export const EMPTY_TEXT = '暂无数据';
export const LOADING_TEXT = '加载中...';
export const ERROR_TEXT = '加载失败：';
export const SEARCH_PLACEHOLDER = 'Search Resources...';

// 菜单配置
export const MENU_LIST = [
  { key: 'overview', label: 'Overview', icon: 'FaChartPie' },
  { key: 'nodes', label: 'Nodes', icon: 'FaDesktop' },
  { key: 'pods', label: 'Pod', icon: 'FaCubes' },
  { key: 'deployments', label: 'Deployment', icon: 'FaLayerGroup' },
  { key: 'statefulsets', label: 'StatefulSet', icon: 'FaDatabase' },
  { key: 'daemonsets', label: 'DaemonSet', icon: 'FaShieldAlt' },
  { key: 'cronjobs', label: 'CronJob', icon: 'FaRegClock' },
  { key: 'jobs', label: 'Job', icon: 'FaTasks' },
  { key: 'ingress', label: 'Ingress', icon: 'FaNetworkWired' },
  { key: 'services', label: 'Service', icon: 'FaProjectDiagram' },
  { key: 'events', label: 'Events', icon: 'FaBell' },
];

// API 映射
export const API_MAP = {
  overview: '/api/overview',
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
}; 