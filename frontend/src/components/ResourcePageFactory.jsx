import React from 'react';
import ResourcePage from './ResourcePage';

/**
 * 资源页面组件工厂
 * 自动生成所有资源页面，消除重复代码
 */
export function createResourcePage(config) {
  return function ResourcePageComponent({ collapsed, onToggleCollapsed }) {
    return (
      <ResourcePage
        {...config}
        collapsed={collapsed}
        onToggleCollapsed={onToggleCollapsed}
      />
    );
  };
}

/**
 * 批量创建资源页面组件
 * @param {Object} configs - 页面配置对象
 * @returns {Object} 页面组件对象
 */
export function createResourcePages(configs) {
  const pages = {};
  
  Object.entries(configs).forEach(([key, config]) => {
    const pageKey = key.replace('_CONFIG', '').toLowerCase();
    pages[pageKey] = createResourcePage(config);
  });
  
  return pages;
}

/**
 * 页面组件映射
 * 自动导入和创建所有页面组件
 */
export const PAGE_COMPONENTS = {
  // 这里可以动态导入，但为了保持简单，我们手动映射
  // 实际使用时，这些组件会通过createResourcePages自动生成
};

// 导出工厂函数供外部使用
export { createResourcePage as default };

