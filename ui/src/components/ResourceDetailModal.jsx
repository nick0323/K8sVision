import React, { useState, useEffect } from 'react';
import './ResourceDetailModal.css';
import {
  PodDetails,
  DeploymentDetails,
  StatefulSetDetails,
  DaemonSetDetails,
  JobDetails,
  CronJobDetails,
  ServiceDetails,
  IngressDetails,
  NodeDetails,
  NamespaceDetails,
  ConfigMapDetails,
  SecretDetails,
  PVCDetails,
  PVDetails,
  StorageClassDetails
} from './ResourceDetailComponents';

// 新增工具函数：时间计算
const calculateDuration = (startTime, endTime) => {
  if (!startTime) return '-';
  if (!endTime) return 'Running';
  
  const start = new Date(startTime);
  const end = new Date(endTime);
  const diffMs = end - start;
  
  if (diffMs < 1000) return `${diffMs}ms`;
  if (diffMs < 60000) return `${Math.round(diffMs / 1000)}s`;
  if (diffMs < 3600000) return `${Math.round(diffMs / 60000)}m${Math.round((diffMs % 60000) / 1000)}s`;
  
  const hours = Math.floor(diffMs / 3600000);
  const minutes = Math.round((diffMs % 3600000) / 60000);
  return `${hours}h${minutes}m`;
};

// 详情卡片组件
function DetailCard({ title, children }) {
  return (
    <div className="detail-card">
      <h3 className="detail-card-title">{title}</h3>
      <div className="detail-card-content">
        {children}
      </div>
    </div>
  );
}

// 详情项组件
function DetailItem({ label, value, type = 'text' }) {
  const renderValue = () => {
    if (!value) return '-';

    switch (type) {
      case 'status':
        const isHealthy = value === 'Running' || value === 'Succeeded' || value === 'Ready' || value === 'Healthy' || value === 'Normal' || value === 'Active' || value === 'Bound';
        const isFailed = value === 'Failed' || value === 'Error' || value === 'CrashLoopBackOff';
        const isPending = value === 'Pending' || value === 'ContainerCreating' || value === 'PodInitializing';
        
        let statusClass = 'status-running';
        if (isHealthy) {
          statusClass = 'status-ready';
        } else if (isFailed) {
          statusClass = 'status-failed';
        } else if (isPending) {
          statusClass = 'status-pending';
        }
        
        return (
          <span className={`status-tag ${statusClass}`}>
            {value}
          </span>
        );
      
      case 'container':
        // 检查是否包含镜像信息 (格式: "name (image)")
        const containerMatch = value.match(/^(.+?)\s*\((.+)\)$/);
        if (containerMatch) {
          const [, containerName, image] = containerMatch;
          return (
            <div className="container-item">
              <span className="label-key">{containerName}</span>
              <span className="label-value">{image}</span>
            </div>
          );
        }
        return value;
      
      case 'labels':
      case 'annotations':
        if (typeof value === 'object' && value !== null) {
          const containerClass = type === 'labels' ? 'labels-container' : 'annotations-container';
          const itemClass = type === 'labels' ? 'label-item' : 'annotation-item';
          const keyClass = type === 'labels' ? 'label-key' : 'annotation-key';
          const valueClass = type === 'labels' ? 'label-value' : 'annotation-value';
          
          return (
            <div className={containerClass}>
              {Object.entries(value).map(([key, val]) => (
                <div key={key} className={itemClass}>
                  <span className={keyClass}>{key}</span>
                  <span className={valueClass}>{val}</span>
                </div>
              ))}
            </div>
          );
        }
        return <span>{JSON.stringify(value)}</span>;

      case 'ports':
      case 'hosts':
      case 'role':
      case 'paths':
      case 'targetServices':
        if (Array.isArray(value)) {
          const label = type === 'ports' ? 'Port' : 
                       type === 'hosts' ? 'Host' : 
                       type === 'role' ? 'Role' : 
                       type === 'paths' ? 'Path' : 'Service';
          
          return (
            <div className="labels-container">
              {value.map((item, index) => (
                <div key={index} className="label-item">
                  <span className="label-key">{label} {index + 1}</span>
                  <span className="label-value">{item}</span>
                </div>
              ))}
            </div>
          );
        }
        return <span>{value}</span>;

      // 增强AccessMode兼容性处理
      case 'accessMode':
        let modes = [];
        if (Array.isArray(value)) {
          modes = value;
        } else if (typeof value === 'string') {
          // 处理逗号分隔字符串
          modes = value.split(',').map(mode => mode.trim()).filter(mode => mode);
        }
        
        if (modes.length > 0) {
          return (
            <div className="labels-container">
              {modes.map((mode, index) => (
                <div key={index} className="label-item">
                  <span className="label-key">Mode {index + 1}</span>
                  <span className="label-value">{mode}</span>
                </div>
              ))}
            </div>
          );
        }
        return <span>-</span>;

      // 新增百分比类型处理
      case 'percentage':
        const roundedValue = typeof value === 'number' ? value.toFixed(2) : value;
        return <span>{roundedValue}%</span>;

      case 'configData':
      case 'secretData':
        if (typeof value === 'object' && value !== null) {
          return (
            <div className="labels-container">
              {Object.entries(value).map(([key, val]) => (
                <div key={key} className="label-item">
                  <span className="label-key">{key}</span>
                  <span className="label-value">{type === 'secretData' ? '***' : val}</span>
                </div>
              ))}
            </div>
          );
        }
        return <span>{JSON.stringify(value)}</span>;

      case 'parameters':
        if (typeof value === 'object' && value !== null) {
          return (
            <div className="labels-container">
              {Object.entries(value).map(([key, val]) => (
                <div key={key} className="label-item">
                  <span className="label-key">{key}</span>
                  <span className="label-value">{val}</span>
                </div>
              ))}
            </div>
          );
        }
        return <span>{JSON.stringify(value)}</span>;
      
      default:
        return <span>{value}</span>;
    }
  };

  return (
    <div className="detail-item">
      {label && <div className="detail-label">{label}:</div>}
      <div className="detail-value">{renderValue()}</div>
    </div>
  );
}































// 主组件
export default function ResourceDetailModal({ 
  visible, 
  resourceType, 
  namespace, 
  name, 
  onClose 
}) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (visible && resourceType && name) {
      fetchResourceDetail(resourceType, namespace, name);
    } else if (!visible) {
      setData(null);
      setError(null);
    }
  }, [visible, resourceType, namespace, name]);

  const fetchResourceDetail = async (resourceType, namespace, name) => {
    try {
      setLoading(true);
      setError(null);
      
      // Build API path
      const apiPath = namespace 
        ? `/api/${resourceType}/${namespace}/${name}`
        : `/api/${resourceType}/${name}`;
      
      const response = await fetch(apiPath);
      const result = await response.json();
      
      if (response.ok) {
        setData(result.data);
      } else {
        setError(result.message || 'Failed to fetch details');
      }
    } catch (err) {
      setError('Network error');
    } finally {
      setLoading(false);
    }
  };

  const handleRefresh = () => {
    if (resourceType && name) {
      fetchResourceDetail(resourceType, namespace, name);
    }
  };

  // 根据资源类型渲染不同的详情内容
  const renderResourceDetails = () => {
    if (!data) return null;

    switch (resourceType) {
      case 'pods':
        return <PodDetails data={data} />;
      case 'deployments':
        return <DeploymentDetails data={data} />;
      case 'statefulsets':
        return <StatefulSetDetails data={data} />;
      case 'daemonsets':
        return <DaemonSetDetails data={data} />;
      case 'jobs':
        return <JobDetails data={data} />;
      case 'cronjobs':
        return <CronJobDetails data={data} />;
      case 'services':
        return <ServiceDetails data={data} />;
      case 'ingress':
        return <IngressDetails data={data} />;
      case 'nodes':
        return <NodeDetails data={data} />;
      case 'namespaces':
        return <NamespaceDetails data={data} />;
      case 'configmaps':
        return <ConfigMapDetails data={data} />;
      case 'secrets':
        return <SecretDetails data={data} />;
      case 'pvcs':
        return <PVCDetails data={data} />;
      case 'pvs':
        return <PVDetails data={data} />;
      case 'storageclasses':
        return <StorageClassDetails data={data} />;
      default:
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="State" value={data.status} type="status" />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );
    }
  };

  if (!visible) return null;

  return (
    <div className="resource-detail-modal" onClick={onClose}>
      <div className="modal-content" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <div className="modal-title-section">
            <h2 className="modal-title">{name}</h2>
            <span className="modal-subtitle">{resourceType.toUpperCase()}</span>
          </div>
          <div className="modal-actions">
            <button className="modal-refresh" onClick={handleRefresh} title="Refresh">
              ↻
            </button>
            <button className="modal-close" onClick={onClose} title="Close">
              ×
            </button>
          </div>
        </div>
        <div className="modal-body">
          {loading && (
            <div className="loading-container">
              <div className="loading-spinner"></div>
              <div className="loading-text">Loading...</div>
            </div>
          )}
          
          {error && (
            <div className="error-container">
              <div className="error-icon">⚠️</div>
              <div className="error-text">Error: {error}</div>
              <button className="error-retry" onClick={handleRefresh}>
                Retry
              </button>
            </div>
          )}
          
          {data && renderResourceDetails()}
        </div>
      </div>
    </div>
  );
} 