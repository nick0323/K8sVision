import React from 'react';

/**
 * 创建状态渲染器
 * @param {Object} statusMap - 状态映射对象
 * @returns {Function} 状态渲染函数
 */
export function createStatusRenderer(statusMap = {}) {
  return (value, row, index, isTooltip) => {
    if (isTooltip) return value;
    
    if (!value) return '-';
    
    // 使用配置的状态映射
    if (statusMap[value]) {
      return (
        <span className={`status-tag ${statusMap[value]}`}>
          {value}
        </span>
      );
    }
    
    // 默认状态判断逻辑
    const isHealthy = value === 'Running' || value === 'Succeeded' || value === 'Ready' || 
                     value === 'Healthy' || value === 'Normal' || value === 'Active' || 
                     value === 'Bound' || value === 'Available';
    const isFailed = value === 'Failed' || value === 'Error' || value === 'CrashLoopBackOff' || 
                     value === 'Unhealthy' || value === 'Warning';
    const isPending = value === 'Pending' || value === 'ContainerCreating' || 
                     value === 'PodInitializing' || value === 'Creating';
    
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
  };
}

/**
 * 通用状态标签组件
 */
export function StatusTag({ value, className = '', size = 'default' }) {
  if (!value) return null;
  
  const isHealthy = value === 'Running' || value === 'Succeeded' || value === 'Ready' || 
                   value === 'Healthy' || value === 'Normal' || value === 'Active' || 
                   value === 'Bound' || value === 'Available';
  const isFailed = value === 'Failed' || value === 'Error' || value === 'CrashLoopBackOff' || 
                   value === 'Unhealthy' || value === 'Warning';
  const isPending = value === 'Pending' || value === 'ContainerCreating' || 
                   value === 'PodInitializing' || value === 'Creating';
  
  let statusClass = 'status-running';
  if (isHealthy) {
    statusClass = 'status-ready';
  } else if (isFailed) {
    statusClass = 'status-failed';
  } else if (isPending) {
    statusClass = 'status-pending';
  }
  
  return (
    <span className={`status-tag ${statusClass} ${className} ${size}`}>
      {value}
    </span>
  );
}

