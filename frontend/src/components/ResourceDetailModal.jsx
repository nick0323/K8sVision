import React, { useState, useEffect } from 'react';
import './ResourceDetailModal.css';

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
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="Status" value={data.status} type="status" />
              <DetailItem label="Pod IP" value={data.podIP} />
              <DetailItem label="Node" value={data.nodeName} />
              <DetailItem label="Start Time" value={data.startTime} />
            </DetailCard>

            <DetailCard title="Containers">
              {data.containers && data.containers.map((container, index) => (
                <DetailItem 
                  key={index} 
                  label={`Container ${index + 1}`} 
                  value={container} 
                  type="container"
                />
              ))}
              {(!data.containers || data.containers.length === 0) && (
                <DetailItem label="" value="No containers found" />
              )}
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'deployments':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="Status" value={data.status} type="status" />
              <DetailItem label="Replicas" value={data.replicas} />
              <DetailItem label="Available" value={data.available} />
              <DetailItem label="Desired" value={data.desired} />
              <DetailItem label="Strategy" value={data.strategy} />
            </DetailCard>

            <DetailCard title="Image Info">
              <DetailItem label="Image" value={data.image} />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'statefulsets':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="Status" value={data.status} type="status" />
              <DetailItem label="Replicas" value={data.replicas} />
              <DetailItem label="Available" value={data.available} />
              <DetailItem label="Desired" value={data.desired} />
              <DetailItem label="Service Name" value={data.serviceName} />
            </DetailCard>

            <DetailCard title="Image Info">
              <DetailItem label="Image" value={data.image} />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'daemonsets':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="Status" value={data.status} type="status" />
              <DetailItem label="Available" value={data.available} />
              <DetailItem label="Desired" value={data.desired} />
            </DetailCard>

            <DetailCard title="Image Info">
              <DetailItem label="Image" value={data.image} />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'jobs':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="Status" value={data.status} type="status" />
              <DetailItem label="Completions" value={data.completions} />
              {/* 新增字段 */}
              <DetailItem label="Succeeded" value={data.succeeded} />
              <DetailItem label="Failed" value={data.failed} />
              <DetailItem label="Start Time" value={data.startTime} />
              <DetailItem label="Completion Time" value={data.completionTime} />
              <DetailItem label="Duration" value={calculateDuration(data.startTime, data.completionTime)} />
            </DetailCard>

            <DetailCard title="Image Info">
              <DetailItem label="Image" value={data.image} />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'cronjobs':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="Status" value={data.status} type="status" />
              <DetailItem label="Schedule" value={data.schedule} />
              <DetailItem label="Last Schedule" value={data.lastScheduleTime} />
              {/* 新增字段 */}
              <DetailItem label="Suspended" value={data.suspend ? 'Yes' : 'No'} />
              <DetailItem label="Active Jobs" value={data.active} />
            </DetailCard>

            <DetailCard title="Image Info">
              <DetailItem label="Image" value={data.image} />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'services':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="Type" value={data.type} />
              <DetailItem label="Cluster IP" value={data.clusterIP} />
              <DetailItem label="Status" value={data.status} type="status" />
            </DetailCard>

            <DetailCard title="Ports">
              <DetailItem label="" value={data.ports} type="ports" />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'ingress':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="Status" value={data.status} type="status" />
              <DetailItem label="Address" value={data.address} />
              <DetailItem label="Class" value={data.class} />
            </DetailCard>

            <DetailCard title="Hosts & Routes">
              <DetailItem label="Hosts" value={data.hosts} type="hosts" />
              <DetailItem label="Ports" value={data.ports} type="ports" />
              <DetailItem label="Paths" value={data.path} type="paths" />
              <DetailItem label="Target Services" value={data.targetService} type="targetServices" />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'nodes':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="IP" value={data.ip} />
              <DetailItem label="Status" value={data.status} type="status" />
              {/* 修复百分比显示 */}
              <DetailItem label="CPU Usage" value={data.cpuUsage} type="percentage" />
              <DetailItem label="Memory Usage" value={data.memoryUsage} type="percentage" />
              <DetailItem label="Pods Used" value={data.podsUsed} />
              <DetailItem label="Pods Capacity" value={data.podsCapacity} />
            </DetailCard>

            <DetailCard title="Role">
              <DetailItem label="" value={data.role} type="role" />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'namespaces':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Status" value={data.status} type="status" />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'configmaps':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="Data Count" value={data.dataCount} />
            </DetailCard>

            <DetailCard title="Data">
              <DetailItem label="" value={data.data} type="configData" />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'secrets':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="Type" value={data.type} />
              <DetailItem label="Data Count" value={data.dataCount} />
            </DetailCard>

            <DetailCard title="Data">
              <DetailItem label="" value={data.data} type="secretData" />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'pvcs':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="Status" value={data.status} type="status" />
              <DetailItem label="Capacity" value={data.capacity} />
              <DetailItem label="Storage Class" value={data.storageClass} />
              <DetailItem label="Volume Name" value={data.volumeName} />
            </DetailCard>

            <DetailCard title="Access Mode">
              <DetailItem label="" value={data.accessMode} type="accessMode" />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'pvs':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Status" value={data.status} type="status" />
              <DetailItem label="Capacity" value={data.capacity} />
              <DetailItem label="Storage Class" value={data.storageClass} />
              <DetailItem label="Claim Ref" value={data.claimRef} />
              <DetailItem label="Reclaim Policy" value={data.reclaimPolicy} />
            </DetailCard>

            <DetailCard title="Access Mode">
              <DetailItem label="" value={data.accessMode} type="accessMode" />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      case 'storageclasses':
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Provisioner" value={data.provisioner} />
              <DetailItem label="Reclaim Policy" value={data.reclaimPolicy} />
              <DetailItem label="Volume Binding Mode" value={data.volumeBindingMode} />
              <DetailItem label="Is Default" value={data.isDefault} />
            </DetailCard>

            <DetailCard title="Parameters">
              <DetailItem label="" value={data.parameters} type="parameters" />
            </DetailCard>

            <DetailCard title="Labels">
              <DetailItem label="" value={data.labels} type="labels" />
            </DetailCard>

            <DetailCard title="Annotations">
              <DetailItem label="" value={data.annotations} type="annotations" />
            </DetailCard>
          </div>
        );

      default:
        return (
          <div className="detail-content">
            <DetailCard title="Basic Info">
              <DetailItem label="Namespace" value={data.namespace} />
              <DetailItem label="Status" value={data.status} type="status" />
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

// Detail card component
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

// Detail item component
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
      
      case 'annotations':
        if (typeof value === 'object' && value !== null) {
          return (
            <div className="annotations-container">
              {Object.entries(value).map(([key, val]) => (
                <div key={key} className="annotation-item">
                  <span className="annotation-key">{key}</span>
                  <span className="annotation-value">{val}</span>
                </div>
              ))}
            </div>
          );
        }
        return <span>{JSON.stringify(value)}</span>;

      case 'ports':
        if (Array.isArray(value)) {
          return (
            <div className="labels-container">
              {value.map((port, index) => (
                <div key={index} className="label-item">
                  <span className="label-key">Port {index + 1}</span>
                  <span className="label-value">{port}</span>
                </div>
              ))}
            </div>
          );
        }
        return <span>{value}</span>;

      case 'hosts':
        if (Array.isArray(value)) {
          return (
            <div className="labels-container">
              {value.map((host, index) => (
                <div key={index} className="label-item">
                  <span className="label-key">Host {index + 1}</span>
                  <span className="label-value">{host}</span>
                </div>
              ))}
            </div>
          );
        }
        return <span>{value}</span>;

      case 'role':
        if (Array.isArray(value)) {
          return (
            <div className="labels-container">
              {value.map((role, index) => (
                <div key={index} className="label-item">
                  <span className="label-key">Role {index + 1}</span>
                  <span className="label-value">{role}</span>
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

      // 新增路径和服务目标显示
      case 'paths':
      case 'targetServices':
        if (Array.isArray(value)) {
          return (
            <div className="labels-container">
              {value.map((item, index) => (
                <div key={index} className="label-item">
                  <span className="label-key">{type === 'paths' ? 'Path' : 'Service'} {index + 1}</span>
                  <span className="label-value">{item}</span>
                </div>
              ))}
            </div>
          );
        }
        return <span>{value}</span>;

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