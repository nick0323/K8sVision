/**
 * 数据处理工具函数模块
 */

// 移除了 useMemo 导入，因为相关函数已被清理

// 注意：已移除未使用的函数: useFilterRows, deepClone, safeJsonParse, safeJsonStringify
// 如需使用这些功能，请重新添加相应的函数实现

/**
 * 防抖函数
 * @param {Function} func - 要防抖的函数
 * @param {number} wait - 等待时间（毫秒）
 * @returns {Function} 防抖后的函数
 */
export function debounce(func, wait) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}

// 注意：已移除未使用的 throttle 函数
// 如需使用节流功能，请重新添加 throttle 函数实现
