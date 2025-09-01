import React from 'react';
import PodDetailPage from './PodDetailPage';

// 模拟的Pod数据 - 更符合真实场景
const mockPodData = {
  metadata: {
    name: 'ai-check-robot-7b9ff5f8dc-gvx8n',
    namespace: 'finq',
    creationTimestamp: '2025-08-28T08:46:55Z',
    labels: {
      name: 'ai-check-robot',
      'pod-template-hash': '7b9ff5f8dc'
    },
    annotations: {
      'cattle.io/timestamp': '2025-08-28T08:49:05Z'
    },
    ownerReferences: [
      {
        kind: 'ReplicaSet',
        name: 'ai-check-robot-7b9ff5f8dc'
      }
    ]
  },
  spec: {
    nodeName: 'node1',
    priority: 0,
    containers: [
      {
        name: 'ai-check-robot',
        image: 'hub.qmhost1.com/parser/ai-check-robot:develop',
        imagePullPolicy: 'IfNotPresent',
        ports: [
          {
            containerPort: 8080,
            protocol: 'TCP'
          }
        ],
        resources: {
          requests: {
            memory: '512Mi',
            cpu: '250m'
          },
          limits: {
            memory: '1Gi',
            cpu: '500m'
          }
        }
      }
    ],
    volumes: [
      {
        name: 'config-volume',
        configMap: {
          name: 'ai-check-robot-config'
        }
      },
      {
        name: 'data-volume',
        emptyDir: {}
      },
      {
        name: 'logs-volume',
        emptyDir: {}
      }
    ]
  },
  status: {
    phase: 'Running',
    podIP: '10.42.0.96',
    hostIP: '172.17.0.2',
    qosClass: 'Burstable',
    containerStatuses: [
      {
        name: 'ai-check-robot',
        ready: true,
        restartCount: 0,
        state: {
          running: {
            startedAt: '2025-08-28T08:46:55Z'
          }
        }
      }
    ],
    conditions: [
      {
        type: 'Ready',
        status: 'True',
        lastProbeTime: null,
        lastTransitionTime: '2025-08-28T08:46:55Z'
      },
      {
        type: 'ContainersReady',
        status: 'True',
        lastProbeTime: null,
        lastTransitionTime: '2025-08-28T08:46:55Z'
      },
      {
        type: 'PodScheduled',
        status: 'True',
        lastProbeTime: null,
        lastTransitionTime: '2025-08-28T08:46:55Z'
      }
    ]
  }
};

const PodDetailDemo = () => {
  return (
    <div style={{ height: '100vh', overflow: 'auto' }}>
      <div style={{ 
        background: 'var(--bg-secondary)', 
        padding: 'var(--spacing-lg)', 
        borderBottom: '1px solid var(--border-primary)',
        marginBottom: 'var(--spacing-lg)'
      }}>
        <h1 style={{ 
          margin: '0 0 var(--spacing-sm) 0', 
          fontSize: 'var(--font-size-2xl)', 
          fontWeight: '700', 
          color: 'var(--text-primary)',
          fontFamily: 'var(--font-main)'
        }}>
          Kubernetes Pod 资源详情页面
        </h1>
        <p style={{ 
          margin: '0', 
          color: 'var(--text-tertiary)', 
          fontSize: 'var(--font-size-sm)',
          fontFamily: 'var(--font-main)'
        }}>
          这是一个专业的Kubernetes Pod资源详情页面，展示了现代化的信息架构设计和用户体验优化
        </p>
      </div>
      <PodDetailPage podData={mockPodData} />
    </div>
  );
};

export default PodDetailDemo;
