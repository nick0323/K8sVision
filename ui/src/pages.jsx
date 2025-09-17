/**
 * 统一的页面组件定义
 * 使用工厂模式消除所有页面组件的重复代码
 */
import React from 'react';
import ResourcePage from './components/ResourcePage';
import OverviewPage from './OverviewPage';
import {
  PODS_CONFIG,
  DEPLOYMENTS_CONFIG,
  STATEFULSETS_CONFIG,
  DAEMONSETS_CONFIG,
  CRONJOBS_CONFIG,
  JOBS_CONFIG,
  INGRESS_CONFIG,
  SERVICES_CONFIG,
  EVENTS_CONFIG,
  PVCS_CONFIG,
  PVS_CONFIG,
  STORAGECLASSES_CONFIG,
  CONFIGMAPS_CONFIG,
  SECRETS_CONFIG,
  NAMESPACES_CONFIG,
  NODES_CONFIG
} from './constants/pageConfigs';

// 资源页面工厂函数
const createResourcePage = (config) => {
  const ResourcePageComponent = ({ collapsed, onToggleCollapsed }) => (
    <ResourcePage
      {...config}
      collapsed={collapsed}
      onToggleCollapsed={onToggleCollapsed}
    />
  );
  
  // 设置显示名称便于调试
  ResourcePageComponent.displayName = `${config.title}Page`;
  
  return ResourcePageComponent;
};

// 页面组件映射 - 使用工厂函数自动生成
export const PAGE_COMPONENTS = {
  overview: OverviewPage,
  
  // 工作负载
  pods: createResourcePage(PODS_CONFIG),
  deployments: createResourcePage(DEPLOYMENTS_CONFIG),
  statefulsets: createResourcePage(STATEFULSETS_CONFIG),
  daemonsets: createResourcePage(DAEMONSETS_CONFIG),
  jobs: createResourcePage(JOBS_CONFIG),
  cronjobs: createResourcePage(CRONJOBS_CONFIG),
  
  // 网络
  ingress: createResourcePage(INGRESS_CONFIG),
  services: createResourcePage(SERVICES_CONFIG),
  
  // 存储
  pvcs: createResourcePage(PVCS_CONFIG),
  pvs: createResourcePage(PVS_CONFIG),
  storageclasses: createResourcePage(STORAGECLASSES_CONFIG),
  
  // 配置
  configmaps: createResourcePage(CONFIGMAPS_CONFIG),
  secrets: createResourcePage(SECRETS_CONFIG),
  
  // 其他
  namespaces: createResourcePage(NAMESPACES_CONFIG),
  nodes: createResourcePage(NODES_CONFIG),
  events: createResourcePage(EVENTS_CONFIG),
};

// 导出个别组件（向后兼容）
export const PodsPage = PAGE_COMPONENTS.pods;
export const DeploymentsPage = PAGE_COMPONENTS.deployments;
export const StatefulSetsPage = PAGE_COMPONENTS.statefulsets;
export const DaemonSetsPage = PAGE_COMPONENTS.daemonsets;
export const CronJobsPage = PAGE_COMPONENTS.cronjobs;
export const JobsPage = PAGE_COMPONENTS.jobs;
export const IngressPage = PAGE_COMPONENTS.ingress;
export const ServicesPage = PAGE_COMPONENTS.services;
export const EventsPage = PAGE_COMPONENTS.events;
export const PVCsPage = PAGE_COMPONENTS.pvcs;
export const PVsPage = PAGE_COMPONENTS.pvs;
export const StorageClassesPage = PAGE_COMPONENTS.storageclasses;
export const ConfigMapsPage = PAGE_COMPONENTS.configmaps;
export const SecretsPage = PAGE_COMPONENTS.secrets;
export const NamespacesPage = PAGE_COMPONENTS.namespaces;
export const NodesPage = PAGE_COMPONENTS.nodes;
