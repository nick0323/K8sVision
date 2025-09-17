import React, { useState, useEffect } from 'react';
import './PodDetailPage.css';

// 状态指示器组件
const StatusIndicator = ({ status, phase }) => {
  const getStatusColor = () => {
    if (phase === 'Running') return '#52c41a';
    if (phase === 'Pending') return '#faad14';
    if (phase === 'Failed') return '#ff4d4f';
    if (phase === 'Succeeded') return '#1890ff';
    if (phase === 'Unknown') return '#d9d9d9';
    return '#d9d9d9';
  };

  return (
    <div className="status-indicator">
      <div 
        className="status-dot" 
        style={{ backgroundColor: getStatusColor() }}
      />
      <span className="status-text">{status}</span>
    </div>
  );
};

// 状态概览卡片组件
const StatusOverviewCard = ({ pod }) => {
  const getReadyContainers = () => {
    if (!pod.status?.containerStatuses) return '0 / 0';
    const ready = pod.status.containerStatuses.filter(c => c.ready).length;
    const total = pod.status.containerStatuses.length;
    return `${ready} / ${total}`;
  };

  const getRestartCount = () => {
    if (!pod.status?.containerStatuses) return 0;
    return pod.status.containerStatuses.reduce((sum, c) => sum + c.restartCount, 0);
  };

  return (
    <div className="status-overview-card">
      <h3>Status Overview</h3>
      <div className="status-grid">
        <div className="status-item">
          <div className="status-label">Phase</div>
          <div className="status-value">
            <StatusIndicator 
              status={pod.status?.phase || 'Unknown'} 
              phase={pod.status?.phase}
            />
          </div>
        </div>
        <div className="status-item">
          <div className="status-label">Ready Containers</div>
          <div className="status-value">{getReadyContainers()}</div>
        </div>
        <div className="status-item">
          <div className="status-label">Restart Count</div>
          <div className="status-value">{getRestartCount()}</div>
        </div>
        <div className="status-item">
          <div className="status-label">Node</div>
          <div className="status-value">
            <a href="#" className="node-link">{pod.spec?.nodeName || '-'}</a>
          </div>
        </div>
      </div>
    </div>
  );
};

// Pod信息卡片组件
const PodInfoCard = ({ pod }) => {
  const formatTime = (timestamp) => {
    if (!timestamp) return '-';
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / (1000 * 60));
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    
    if (diffDays > 0) return `${date.toLocaleDateString()} (${diffDays}d)`;
    if (diffHours > 0) return `${date.toLocaleDateString()} ${date.toLocaleTimeString()} (${diffHours}h)`;
    if (diffMins > 0) return `${date.toLocaleDateString()} ${date.toLocaleTimeString()} (${diffMins}m)`;
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
  };

  const getOwner = () => {
    if (pod.metadata?.ownerReferences && pod.metadata.ownerReferences.length > 0) {
      const owner = pod.metadata.ownerReferences[0];
      return `${owner.kind}/${owner.name}`;
    }
    return '-';
  };

  return (
    <div className="pod-info-card">
      <h3>Pod Information</h3>
      <div className="info-grid">
        <div className="info-item">
          <div className="info-label">Created</div>
          <div className="info-value">{formatTime(pod.metadata?.creationTimestamp)}</div>
        </div>
        <div className="info-item">
          <div className="info-label">Pod IP</div>
          <div className="info-value">{pod.status?.podIP || '-'}</div>
        </div>
        <div className="info-item">
          <div className="info-label">Owner</div>
          <div className="info-value">
            <a href="#" className="owner-link">{getOwner()}</a>
          </div>
        </div>
        <div className="info-item">
          <div className="info-label">Started</div>
          <div className="info-value">
            {pod.status?.containerStatuses?.[0]?.state?.running?.startedAt ? 
              formatTime(pod.status.containerStatuses[0].state.running.startedAt) : '-'}
          </div>
        </div>
        <div className="info-item">
          <div className="info-label">Host IP</div>
          <div className="info-value">{pod.status?.hostIP || '-'}</div>
        </div>
      </div>
    </div>
  );
};

// 标签和注解卡片组件
const LabelsAnnotationsCard = ({ pod }) => {
  return (
    <div className="labels-annotations-card">
      <div className="labels-section">
        <h4>Labels</h4>
        <div className="labels-display">
          {pod.metadata?.labels ? (
            Object.entries(pod.metadata.labels).map(([key, value]) => (
              <span key={key} className="label-tag">
                {key}: {value}
              </span>
            ))
          ) : (
            <div className="no-data">No labels</div>
          )}
        </div>
      </div>
      
      <div className="annotations-section">
        <h4>Annotations</h4>
        <div className="annotations-display">
          {pod.metadata?.annotations ? (
            Object.entries(pod.metadata.annotations).map(([key, value]) => (
              <div key={key} className="annotation-item">
                <span className="annotation-key">{key}</span>
                <span className="annotation-value">{value}</span>
              </div>
            ))
          ) : (
            <div className="no-data">No annotations</div>
          )}
        </div>
      </div>
    </div>
  );
};

// 容器卡片组件
const ContainersCard = ({ pod }) => {
  return (
    <div className="containers-card">
      <h3>Containers</h3>
      <div className="containers-list">
        {pod.spec?.containers?.map((container, index) => (
          <div key={index} className="container-item">
            <div className="container-header">
              <div className="container-name">{container.name}</div>
              <div className="container-status">
                <span className="status-badge ready">Always</span>
              </div>
            </div>
            <div className="container-details">
              <div className="detail-row">
                <span className="detail-label">Image:</span>
                <span className="detail-value">{container.image}</span>
              </div>
              {container.ports && container.ports.length > 0 && (
                <div className="detail-row">
                  <span className="detail-label">Ports:</span>
                  <span className="detail-value">
                    {container.ports.map(port => 
                      `${port.containerPort}/${port.protocol || 'TCP'}`
                    ).join(', ')}
                  </span>
                </div>
              )}
            </div>
          </div>
        )) || <div className="no-data">No container information</div>}
      </div>
    </div>
  );
};

// 标签页组件
const InfoTabs = ({ pod, activeTab, onTabChange }) => {
  const tabs = [
    { key: 'overview', label: 'Overview', icon: '📊' },
    { key: 'yaml', label: 'YAML', icon: '📄' },
    { key: 'logs', label: 'Logs', icon: '📋' },
    { key: 'terminal', label: 'Terminal', icon: '💻' },
    { key: 'volumes', label: 'Volumes', icon: '💾', badge: pod.spec?.volumes?.length || 0 },
    { key: 'related', label: 'Related', icon: '🔗' },
    { key: 'events', label: 'Events', icon: '📝' },
    { key: 'monitor', label: 'Monitor', icon: '📈' }
  ];

  return (
    <div className="info-tabs">
      <div className="tab-header">
        {tabs.map(tab => (
          <button
            key={tab.key}
            className={`tab-button ${activeTab === tab.key ? 'active' : ''}`}
            onClick={() => onTabChange(tab.key)}
          >
            <span className="tab-icon">{tab.icon}</span>
            {tab.label}
            {tab.badge && <span className="tab-badge">{tab.badge}</span>}
          </button>
        ))}
      </div>
      
      <div className="tab-content">
        {activeTab === 'overview' && <OverviewTab pod={pod} />}
        {activeTab === 'yaml' && <YAMLTab pod={pod} />}
        {activeTab === 'logs' && <LogsTab pod={pod} />}
        {activeTab === 'terminal' && <TerminalTab pod={pod} />}
        {activeTab === 'volumes' && <VolumesTab pod={pod} />}
        {activeTab === 'related' && <RelatedTab pod={pod} />}
        {activeTab === 'events' && <EventsTab pod={pod} />}
        {activeTab === 'monitor' && <MonitorTab pod={pod} />}
      </div>
    </div>
  );
};

// 概览标签页
const OverviewTab = ({ pod }) => {
  return (
    <div className="overview-tab">
      <StatusOverviewCard pod={pod} />
      <PodInfoCard pod={pod} />
      <LabelsAnnotationsCard pod={pod} />
      <ContainersCard pod={pod} />
    </div>
  );
};

// YAML标签页
const YAMLTab = ({ pod }) => {
  const yamlContent = `apiVersion: v1
kind: Pod
metadata:
  name: ${pod.metadata?.name || 'example-pod'}
  namespace: ${pod.metadata?.namespace || 'default'}
  labels:
    app: ${pod.metadata?.labels?.app || 'example'}
spec:
  containers:
  - name: ${pod.spec?.containers?.[0]?.name || 'example-container'}
    image: ${pod.spec?.containers?.[0]?.image || 'nginx:latest'}
    ports:
    - containerPort: 80`;

  return (
    <div className="yaml-tab">
      <div className="yaml-header">
        <h3>YAML Configuration</h3>
        <div className="yaml-actions">
          <button className="btn btn-primary">Edit</button>
          <button className="btn btn-secondary">Copy</button>
          <button className="btn btn-secondary">Download</button>
        </div>
      </div>
      <div className="yaml-content">
        <pre className="yaml-text">{yamlContent}</pre>
      </div>
    </div>
  );
};

// 日志标签页
const LogsTab = ({ pod }) => {
  const [logs, setLogs] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setTimeout(() => {
      setLogs(`2024-01-01T10:00:00Z [INFO] Starting nginx...
2024-01-01T10:00:01Z [INFO] nginx started successfully
2024-01-01T10:00:02Z [INFO] Listening on port 80
2024-01-01T10:00:03Z [INFO] Ready to handle requests`);
      setLoading(false);
    }, 1000);
  }, []);

  if (loading) {
    return <div className="loading">Loading logs...</div>;
  }

  return (
    <div className="logs-tab">
      <div className="logs-header">
        <div className="logs-controls">
          <button className="btn btn-primary">Real-time Logs</button>
          <button className="btn btn-secondary">Download Logs</button>
          <button className="btn btn-secondary">Clear Logs</button>
        </div>
        <div className="logs-filter">
          <input type="text" placeholder="Search logs..." className="log-search" />
        </div>
      </div>
      <div className="logs-content">
        <pre className="logs-text">{logs}</pre>
      </div>
    </div>
  );
};

// 终端标签页
const TerminalTab = ({ pod }) => {
  return (
    <div className="terminal-tab">
      <div className="terminal-header">
        <h3>Terminal</h3>
        <div className="terminal-actions">
          <button className="btn btn-primary">Connect</button>
          <button className="btn btn-secondary">Disconnect</button>
        </div>
      </div>
      <div className="terminal-content">
        <div className="terminal-placeholder">
          <div className="terminal-icon">💻</div>
          <p>Click "Connect" button to start terminal session</p>
        </div>
      </div>
    </div>
  );
};

// 卷标签页
const VolumesTab = ({ pod }) => {
  return (
    <div className="volumes-tab">
      <div className="volumes-header">
        <h3>Volumes</h3>
      </div>
      <div className="volumes-list">
        {pod.spec?.volumes?.map((volume, index) => (
          <div key={index} className="volume-item">
            <div className="volume-name">{volume.name}</div>
            <div className="volume-type">
              {volume.configMap ? 'ConfigMap' : 
               volume.secret ? 'Secret' : 
               volume.emptyDir ? 'EmptyDir' : 'Other'}
            </div>
          </div>
        )) || <div className="no-data">No volumes mounted</div>}
      </div>
    </div>
  );
};

// 相关资源标签页
const RelatedTab = ({ pod }) => {
  return (
    <div className="related-tab">
      <div className="related-header">
        <h3>Related Resources</h3>
      </div>
      <div className="related-list">
        <div className="related-item">
          <span className="related-label">Service:</span>
          <a href="#" className="related-link">ai-check-robot-service</a>
        </div>
        <div className="related-item">
          <span className="related-label">ConfigMap:</span>
          <a href="#" className="related-link">ai-check-robot-config</a>
        </div>
      </div>
    </div>
  );
};

// 事件标签页
const EventsTab = ({ pod }) => {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setTimeout(() => {
      setEvents([
        {
          type: 'Normal',
          reason: 'Scheduled',
          message: 'Successfully assigned default/nginx-pod to node-1',
          lastTimestamp: new Date().toISOString(),
          count: 1
        },
        {
          type: 'Normal',
          reason: 'Pulling',
          message: 'Pulling image "nginx:latest"',
          lastTimestamp: new Date(Date.now() - 60000).toISOString(),
          count: 1
        }
      ]);
      setLoading(false);
    }, 1000);
  }, []);

  if (loading) {
    return <div className="loading">Loading events...</div>;
  }

  return (
    <div className="events-tab">
      <div className="events-list">
        {events.map((event, index) => (
          <div key={index} className={`event-item ${event.type.toLowerCase()}`}>
            <div className="event-header">
              <span className={`event-type ${event.type.toLowerCase()}`}>
                {event.type}
              </span>
              <span className="event-reason">{event.reason}</span>
              <span className="event-time">
                {new Date(event.lastTimestamp).toLocaleString()}
              </span>
            </div>
            <div className="event-message">{event.message}</div>
            <div className="event-count">Count: {event.count}</div>
          </div>
        ))}
        {events.length === 0 && (
          <div className="no-data">No events</div>
        )}
      </div>
    </div>
  );
};

// 监控标签页
const MonitorTab = ({ pod }) => {
  return (
    <div className="monitor-tab">
      <div className="monitor-header">
        <h3>Monitor</h3>
      </div>
      <div className="monitor-content">
        <div className="monitor-placeholder">
          <div className="monitor-icon">📈</div>
          <p>Loading monitoring data...</p>
        </div>
      </div>
    </div>
  );
};

// 操作栏组件
const ActionBar = ({ pod }) => {
  const handleAction = (action) => {
    console.log(`执行操作: ${action}`, pod);
  };

  return (
    <div className="action-bar">
      <div className="action-group">
        <button 
          className="btn btn-primary"
          onClick={() => handleAction('refresh')}
        >
          🔄 Refresh
        </button>
      </div>
      
      <div className="action-group">
        <button 
          className="btn btn-danger"
          onClick={() => handleAction('delete')}
        >
          🗑️ Delete
        </button>
      </div>
    </div>
  );
};

// 主组件
const PodDetailPage = ({ podData }) => {
  const [activeTab, setActiveTab] = useState('overview');
  const [pod, setPod] = useState(podData);

  useEffect(() => {
    setPod(podData);
  }, [podData]);

  if (!pod) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="pod-detail-page">
      <div className="pod-detail-container">
        <div className="page-header">
          <div className="header-content">
            <h1 className="resource-name">{pod.metadata?.name || 'Unknown Pod'}</h1>
                    <div className="resource-meta">
          <span className="namespace">Namespace: {pod.metadata?.namespace || 'default'}</span>
        </div>
          </div>
          <ActionBar pod={pod} />
        </div>
        
        <InfoTabs 
          pod={pod} 
          activeTab={activeTab} 
          onTabChange={setActiveTab} 
        />
      </div>
    </div>
  );
};

export default PodDetailPage;
