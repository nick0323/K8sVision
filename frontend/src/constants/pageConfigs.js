// 页面配置常量
// 使用通用表格工具消除重复代码
import { PREDEFINED_COLUMNS } from '../utils/tableUtils';

// 通用配置生成器
const createResourceConfig = (title, resourceType, columns, statusMap = {}, tableConfig = {}, namespaceFilter = true) => ({
  title,
  apiEndpoint: `/api/${resourceType}`,
  resourceType,
  columns,
  statusMap,
  tableConfig,
  namespaceFilter
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
  PREDEFINED_COLUMNS.suspend(),
  PREDEFINED_COLUMNS.custom('Last Schedule', 'lastScheduleTime'),
  PREDEFINED_COLUMNS.custom('Active', 'active'),
  PREDEFINED_COLUMNS.status(),
]);

// Jobs页面配置
export const JOBS_CONFIG = createResourceConfig('Jobs', 'jobs', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.custom('Completions', 'completions'),
  PREDEFINED_COLUMNS.custom('Succeeded', 'succeeded'),
  PREDEFINED_COLUMNS.custom('Failed', 'failed'),
  PREDEFINED_COLUMNS.startTime(),
  PREDEFINED_COLUMNS.completionTime(),
  PREDEFINED_COLUMNS.status(),
]);

// Ingress页面配置
export const INGRESS_CONFIG = createResourceConfig('Ingresses', 'ingress', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.class(),
  PREDEFINED_COLUMNS.hosts(),
  PREDEFINED_COLUMNS.path(),
  PREDEFINED_COLUMNS.targetService(),
  PREDEFINED_COLUMNS.state(),
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
  PREDEFINED_COLUMNS.eventType(),
  PREDEFINED_COLUMNS.reason(),
  PREDEFINED_COLUMNS.message(),
  PREDEFINED_COLUMNS.eventName(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.lastTimestamp(),
  PREDEFINED_COLUMNS.duration(),
  PREDEFINED_COLUMNS.count(),
], {}, {
  // 特殊表格配置：允许内容换行展示
  className: 'events-table-wrap-content',
  scroll: false,
  wrap: true
});

// PVCs页面配置
export const PVCS_CONFIG = createResourceConfig('PersistentVolumeClaims', 'pvcs', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.pvcState(),
  PREDEFINED_COLUMNS.volumeName(),
  PREDEFINED_COLUMNS.capacity(),
  PREDEFINED_COLUMNS.accessMode(),
  PREDEFINED_COLUMNS.storageClass(),
]);

// PVs页面配置 - 集群级别资源
export const PVS_CONFIG = createResourceConfig('PersistentVolumes', 'pvs', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.pvState(),
  PREDEFINED_COLUMNS.capacity(),
  PREDEFINED_COLUMNS.accessMode(),
  PREDEFINED_COLUMNS.storageClass(),
  PREDEFINED_COLUMNS.claimRef(),
  PREDEFINED_COLUMNS.reclaimPolicy(),
], {}, {}, false); // 集群级别资源，不需要namespace筛选器

// StorageClasses页面配置 - 集群级别资源
export const STORAGECLASSES_CONFIG = createResourceConfig('StorageClasses', 'storageclasses', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.custom('Provisioner', 'provisioner'),
  PREDEFINED_COLUMNS.custom('ReclaimPolicy', 'reclaimPolicy'),
  PREDEFINED_COLUMNS.custom('VolumeBindingMode', 'volumeBindingMode'),
  PREDEFINED_COLUMNS.isDefault(),
], {}, {}, false); // 集群级别资源，不需要namespace筛选器

// ConfigMaps页面配置
export const CONFIGMAPS_CONFIG = createResourceConfig('ConfigMaps', 'configmaps', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.data(),
  PREDEFINED_COLUMNS.keys(),
], {}, {
  // 特殊表格配置：允许内容换行展示
  className: 'configmaps-table-wrap-content',
  scroll: false,
  wrap: true
});

// Secrets页面配置
export const SECRETS_CONFIG = createResourceConfig('Secrets', 'secrets', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.namespace(),
  PREDEFINED_COLUMNS.type(),
  PREDEFINED_COLUMNS.data(),
  PREDEFINED_COLUMNS.keys(),
]);

// Namespaces页面配置 - 集群级别资源
export const NAMESPACES_CONFIG = createResourceConfig('Namespaces', 'namespaces', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.status(),
], {}, {}, false); // 集群级别资源，不需要namespace筛选器

// Nodes页面配置 - 集群级别资源
export const NODES_CONFIG = createResourceConfig('Nodes', 'nodes', [
  PREDEFINED_COLUMNS.name(),
  PREDEFINED_COLUMNS.internalIP(),
  PREDEFINED_COLUMNS.roles(),
  PREDEFINED_COLUMNS.cpuUsage(),
  PREDEFINED_COLUMNS.memoryUsage(),
  PREDEFINED_COLUMNS.pods(),
  PREDEFINED_COLUMNS.status(),
], {}, {}, false); // 集群级别资源，不需要namespace筛选器

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
