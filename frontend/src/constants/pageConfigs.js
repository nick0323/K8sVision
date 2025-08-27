// 页面配置常量
// 使用通用表格工具消除重复代码
import { PREDEFINED_COLUMNS } from '../utils/tableUtils';

// 通用配置生成器
const createResourceConfig = (title, resourceType, columns, statusMap = {}) => ({
  title,
  apiEndpoint: `/api/${resourceType}`,
  resourceType,
  columns,
  statusMap
});

// Pods页面配置
export const PODS_CONFIG = createResourceConfig('Pods', 'pods', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.status(),
  PREDEFINED_COLUMNS.cpuUsage(),
  PREDEFINED_COLUMNS.memoryUsage(),
  PREDEFINED_COLUMNS.podIP(),
  PREDEFINED_COLUMNS.nodeName(),
]);

// Deployments页面配置
export const DEPLOYMENTS_CONFIG = createResourceConfig('Deployments', 'deployments', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.availableReplicas(),
  PREDEFINED_COLUMNS.desiredReplicas(),
  PREDEFINED_COLUMNS.status({ statusMap: {
    'Healthy': 'status-ready',
    'PartialAvailable': 'status-pending',
    'Abnormal': 'status-failed'
  }}),
]);

// StatefulSets页面配置
export const STATEFULSETS_CONFIG = createResourceConfig('StatefulSets', 'statefulsets', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.availableReplicas(),
  PREDEFINED_COLUMNS.desiredReplicas(),
  PREDEFINED_COLUMNS.status(),
]);

// DaemonSets页面配置
export const DAEMONSETS_CONFIG = createResourceConfig('DaemonSets', 'daemonsets', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.availableReplicas(),
  PREDEFINED_COLUMNS.desiredReplicas(),
  PREDEFINED_COLUMNS.status(),
]);

// CronJobs页面配置
export const CRONJOBS_CONFIG = createResourceConfig('CronJobs', 'cronjobs', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.custom('Schedule', 'schedule'),
  PREDEFINED_COLUMNS.custom('Last Run', 'lastRun'),
  PREDEFINED_COLUMNS.status(),
]);

// Jobs页面配置
export const JOBS_CONFIG = createResourceConfig('Jobs', 'jobs', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.custom('Completions', 'completions'),
  PREDEFINED_COLUMNS.custom('Parallelism', 'parallelism'),
  PREDEFINED_COLUMNS.status(),
]);

// Ingress页面配置
export const INGRESS_CONFIG = createResourceConfig('Ingresses', 'ingress', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.hosts(),
  PREDEFINED_COLUMNS.custom('Address', 'address'),
  PREDEFINED_COLUMNS.status(),
]);

// Services页面配置
export const SERVICES_CONFIG = createResourceConfig('Services', 'services', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.type(),
  PREDEFINED_COLUMNS.clusterIP(),
  PREDEFINED_COLUMNS.ports(),
]);

// Events页面配置
export const EVENTS_CONFIG = createResourceConfig('Events', 'events', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.type(),
  PREDEFINED_COLUMNS.reason(),
  PREDEFINED_COLUMNS.message(),
  PREDEFINED_COLUMNS.lastTimestamp(),
]);

// PVCs页面配置
export const PVCS_CONFIG = createResourceConfig('PVCs', 'pvcs', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.status(),
  PREDEFINED_COLUMNS.capacity(),
  PREDEFINED_COLUMNS.accessMode(),
  PREDEFINED_COLUMNS.storageClass(),
]);

// PVs页面配置
export const PVS_CONFIG = createResourceConfig('PVs', 'pvs', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.status(),
  PREDEFINED_COLUMNS.capacity(),
  PREDEFINED_COLUMNS.accessMode(),
  PREDEFINED_COLUMNS.storageClass(),
  PREDEFINED_COLUMNS.custom('Claim', 'claim'),
]);

// StorageClasses页面配置
export const STORAGECLASSES_CONFIG = createResourceConfig('StorageClasses', 'storageclasses', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.custom('Provisioner', 'provisioner'),
  PREDEFINED_COLUMNS.custom('Reclaim Policy', 'reclaimPolicy'),
  PREDEFINED_COLUMNS.custom('Volume Binding Mode', 'volumeBindingMode'),
  PREDEFINED_COLUMNS.allowVolumeExpansion(),
]);

// ConfigMaps页面配置
export const CONFIGMAPS_CONFIG = createResourceConfig('ConfigMaps', 'configmaps', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.data(),
  PREDEFINED_COLUMNS.created(),
]);

// Secrets页面配置
export const SECRETS_CONFIG = createResourceConfig('Secrets', 'secrets', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.type(),
  PREDEFINED_COLUMNS.data(),
  PREDEFINED_COLUMNS.created(),
]);

// Namespaces页面配置
export const NAMESPACES_CONFIG = createResourceConfig('Namespaces', 'namespaces', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.status(),
  PREDEFINED_COLUMNS.age(),
  PREDEFINED_COLUMNS.labels(),
]);

// Nodes页面配置
export const NODES_CONFIG = createResourceConfig('Nodes', 'nodes', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.status(),
  PREDEFINED_COLUMNS.roles(),
  PREDEFINED_COLUMNS.age(),
  PREDEFINED_COLUMNS.version(),
  PREDEFINED_COLUMNS.internalIP(),
]);

// 开发环境下的配置验证
if (process.env.NODE_ENV === 'development') {
  // 验证所有配置的结构
  const allConfigs = [
    PODS_CONFIG, DEPLOYMENTS_CONFIG, STATEFULSETS_CONFIG, DAEMONSETS_CONFIG,
    CRONJOBS_CONFIG, JOBS_CONFIG, INGRESS_CONFIG, SERVICES_CONFIG,
    EVENTS_CONFIG, PVCS_CONFIG, PVS_CONFIG, STORAGECLASSES_CONFIG,
    CONFIGMAPS_CONFIG, SECRETS_CONFIG, NAMESPACES_CONFIG, NODES_CONFIG
  ];

  allConfigs.forEach((config, index) => {
    const configNames = [
      'PODS_CONFIG', 'DEPLOYMENTS_CONFIG', 'STATEFULSETS_CONFIG', 'DAEMONSETS_CONFIG',
      'CRONJOBS_CONFIG', 'JOBS_CONFIG', 'INGRESS_CONFIG', 'SERVICES_CONFIG',
      'EVENTS_CONFIG', 'PVCS_CONFIG', 'PVS_CONFIG', 'STORAGECLASSES_CONFIG',
      'CONFIGMAPS_CONFIG', 'SECRETS_CONFIG', 'NAMESPACES_CONFIG', 'NODES_CONFIG'
    ];
    
    if (!config.title || !config.apiEndpoint || !config.resourceType || !config.columns) {
      console.error(`Configuration validation failed for ${configNames[index]}:`, config);
    }
  });
}
