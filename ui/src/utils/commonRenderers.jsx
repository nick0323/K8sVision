import React from 'react';
import { StatusTag } from '../components/StatusRenderer';

/**
 * 通用列渲染器集合
 * 消除重复的渲染逻辑
 */

// 资源名称渲染器（可点击）
export function createNameRenderer(onClick, className = 'resource-name-link') {
  return (value, row, index, isTooltip) => {
    if (isTooltip) return value;
    
    if (!value) return '-';
    
    return (
      <span 
        className={className}
        onClick={() => onClick(row)}
        style={{ cursor: 'pointer' }}
      >
        {value}
      </span>
    );
  };
}

// 命名空间渲染器
export function createNamespaceRenderer() {
  return (value, row, index, isTooltip) => {
    if (isTooltip) return value;
    
    if (!value) return 'default';
    
    return (
      <span className="namespace-tag">
        {value}
      </span>
    );
  };
}

// 状态渲染器
export function createStatusRenderer(statusMap = {}) {
  return (value, row, index, isTooltip) => {
    if (isTooltip) return value;
    
    if (!value) return '-';
    
    return <StatusTag value={value} />;
  };
}

// 标签渲染器
export function createLabelsRenderer() {
  return (value, row, index, isTooltip) => {
    if (isTooltip) return JSON.stringify(value);
    
    if (!value || typeof value !== 'object') return '-';
    
    return (
      <div className="labels-container">
        {Object.entries(value).slice(0, 3).map(([key, val]) => (
          <div key={key} className="label-item">
            <span className="label-key">{key}</span>
            <span className="label-value">{val}</span>
          </div>
        ))}
        {Object.keys(value).length > 3 && (
          <span className="more-labels">+{Object.keys(value).length - 3} more</span>
        )}
      </div>
    );
  };
}

// 时间渲染器
export function createTimeRenderer() {
  return (value, row, index, isTooltip) => {
    if (isTooltip) return value;
    
    if (!value) return '-';
    
    try {
      const date = new Date(value);
      if (isNaN(date.getTime())) return value;
      
      const now = new Date();
      const diffMs = now - date;
      const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
      
      if (diffDays === 0) {
        return 'Today';
      } else if (diffDays === 1) {
        return 'Yesterday';
      } else if (diffDays < 7) {
        return `${diffDays} days ago`;
      } else {
        return date.toLocaleDateString();
      }
    } catch (e) {
      return value;
    }
  };
}

// 详细时间渲染器（显示完整日期时间格式，上海时间）
export function createDetailedTimeRenderer() {
  return (value, row, index, isTooltip) => {
    if (isTooltip) return value;
    
    if (!value) return '-';
    
    try {
      const date = new Date(value);
      if (isNaN(date.getTime())) return value;
      
      // 转换为上海时间（东八区 UTC+8）
      const shanghaiTime = new Date(date.getTime() + (8 * 60 * 60 * 1000));
      
      // 显示完整的日期时间格式：YYYY-MM-DD HH:MM:SS
      return shanghaiTime.toISOString().slice(0, 19).replace('T', ' ');
    } catch (e) {
      return value;
    }
  };
}

// 数字渲染器
export function createNumberRenderer(suffix = '', formatter = null) {
  return (value, row, index, isTooltip) => {
    if (isTooltip) return value;
    
    if (value === null || value === undefined) return '-';
    
    const num = Number(value);
    if (isNaN(num)) return value;
    
    if (formatter) {
      return formatter(num);
    }
    
    return `${num}${suffix}`;
  };
}

// 布尔值渲染器
export function createBooleanRenderer() {
  return (value, row, index, isTooltip) => {
    if (isTooltip) return value ? 'true' : 'false';
    
    if (typeof value === 'boolean') {
      return (
        <span className={`boolean-tag ${value ? 'true' : 'false'}`}>
          {value ? 'true' : 'false'}
        </span>
      );
    }
    
    return value || '-';
  };
}

// 数组渲染器
export function createArrayRenderer(separator = ', ', maxItems = 3) {
  return (value, row, index, isTooltip) => {
    if (isTooltip) return Array.isArray(value) ? value.join(separator) : value;
    
    if (!Array.isArray(value)) return value || '-';
    
    if (value.length === 0) return '-';
    
    // 去重并过滤空值
    const uniqueValues = [...new Set(value.filter(item => item && item.trim()))];
    
    const displayItems = uniqueValues.slice(0, maxItems);
    const hasMore = uniqueValues.length > maxItems;
    
    return (
      <span>
        {displayItems.join(separator)}
        {hasMore && <span className="more-items"> (+{value.length - maxItems})</span>}
      </span>
    );
  };
}

// 去重数组渲染器
export function createUniqueArrayRenderer(separator = ', ', maxItems = 3) {
  return (value, row, index, isTooltip) => {
    if (!Array.isArray(value)) return value || '-';
    
    if (value.length === 0) return '-';
    
    // 去重并过滤空值
    const uniqueValues = [...new Set(value.filter(item => item && item.trim()))];
    
    if (isTooltip) return uniqueValues.join(separator);
    
    const displayItems = uniqueValues.slice(0, maxItems);
    const hasMore = uniqueValues.length > maxItems;
    
    return (
      <span>
        {displayItems.join(separator)}
        {hasMore && <span className="more-items"> (+{uniqueValues.length - maxItems})</span>}
      </span>
    );
  };
}

// 资源使用量渲染器
export function createUsageRenderer() {
  return (value, row, index, isTooltip) => {
    if (isTooltip) return value;
    
    if (!value) return '-';
    
    // 处理CPU和内存使用量
    if (typeof value === 'string') {
      if (value.includes('m')) {
        // CPU millicores
        const cores = parseFloat(value) / 1000;
        return `${cores.toFixed(2)} cores`;
      } else if (value.includes('Mi') || value.includes('Ki')) {
        // 内存
        return value;
      }
    }
    
    // 处理数值类型的百分比
    if (typeof value === 'number' || !isNaN(Number(value))) {
      const num = Number(value);
      // 如果数值在0-100范围内，认为是百分比
      if (num >= 0 && num <= 100) {
        return `${num.toFixed(1)}%`;
      }
      // 如果数值大于100，可能是小数形式，转换为百分比
      if (num > 1) {
        return `${(num / 100).toFixed(1)}%`;
      }
      // 如果数值在0-1之间，直接转换为百分比
      if (num >= 0 && num <= 1) {
        return `${(num * 100).toFixed(1)}%`;
      }
    }
    
    return value;
  };
}
