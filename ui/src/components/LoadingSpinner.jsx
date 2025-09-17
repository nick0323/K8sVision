import React from 'react';
import './LoadingSpinner.css';

/**
 * 加载动画组件
 * @param {Object} props - 组件属性
 * @param {string} props.type - 加载动画类型 ('spinner', 'skeleton', 'pulse', 'progress')
 * @param {string} props.text - 加载文本
 * @param {string} props.size - 大小 ('sm', 'md', 'lg')
 * @param {string} props.className - 额外的CSS类名
 * @returns {JSX.Element} 加载动画组件
 */
export default function LoadingSpinner({ 
  type = 'spinner', 
  text = 'Loading...', 
  size = 'md',
  className = '' 
}) {
  const sizeClasses = {
    sm: 'loading-sm',
    md: 'loading-md',
    lg: 'loading-lg'
  };

  const renderLoadingContent = () => {
    switch (type) {
      case 'skeleton':
        return (
          <div className={`skeleton-container ${sizeClasses[size]} ${className}`}>
            <div className="skeleton skeleton-title"></div>
            <div className="skeleton skeleton-text"></div>
            <div className="skeleton skeleton-text"></div>
            <div className="skeleton skeleton-text"></div>
          </div>
        );
      
      case 'pulse':
        return (
          <div className={`loading-container ${className}`}>
            <div className={`pulse ${sizeClasses[size]}`}>
              <div className="loading-spinner"></div>
            </div>
            {text && <div className="loading-text">{text}</div>}
          </div>
        );
      
      case 'progress':
        return (
          <div className={`loading-container ${className}`}>
            <div className="progress-bar">
              <div className="progress-fill"></div>
            </div>
            {text && <div className="loading-text">{text}</div>}
          </div>
        );
      
      case 'spinner':
      default:
        return (
          <div className={`loading-container ${className}`}>
            <div className={`loading-spinner ${sizeClasses[size]}`}></div>
            {text && <div className="loading-text">{text}</div>}
          </div>
        );
    }
  };

  return renderLoadingContent();
}

/**
 * 骨架屏加载组件
 * @param {Object} props - 组件属性
 * @param {number} props.rows - 骨架行数
 * @param {string} props.className - 额外的CSS类名
 * @returns {JSX.Element} 骨架屏组件
 */
export function SkeletonLoader({ rows = 3, className = '' }) {
  return (
    <div className={`skeleton-loader ${className}`}>
      {Array.from({ length: rows }).map((_, index) => (
        <div key={index} className="skeleton skeleton-text"></div>
      ))}
    </div>
  );
}

/**
 * 表格骨架屏组件
 * @param {Object} props - 组件属性
 * @param {number} props.rows - 骨架行数
 * @param {number} props.columns - 骨架列数
 * @param {string} props.className - 额外的CSS类名
 * @returns {JSX.Element} 表格骨架屏组件
 */
export function TableSkeleton({ rows = 5, columns = 4, className = '' }) {
  return (
    <div className={`table-skeleton ${className}`}>
      {/* 表头骨架 */}
      <div className="skeleton-row skeleton-header">
        {Array.from({ length: columns }).map((_, index) => (
          <div key={index} className="skeleton skeleton-text"></div>
        ))}
      </div>
      
      {/* 表行骨架 */}
      {Array.from({ length: rows }).map((_, rowIndex) => (
        <div key={rowIndex} className="skeleton-row">
          {Array.from({ length: columns }).map((_, colIndex) => (
            <div key={colIndex} className="skeleton skeleton-text"></div>
          ))}
        </div>
      ))}
    </div>
  );
}

/**
 * 卡片骨架屏组件
 * @param {Object} props - 组件属性
 * @param {string} props.className - 额外的CSS类名
 * @returns {JSX.Element} 卡片骨架屏组件
 */
export function CardSkeleton({ className = '' }) {
  return (
    <div className={`card-skeleton ${className}`}>
      <div className="skeleton skeleton-title"></div>
      <div className="skeleton skeleton-text"></div>
      <div className="skeleton skeleton-text"></div>
      <div className="skeleton skeleton-button"></div>
    </div>
  );
}
