/**
 * 数据处理工具函数模块
 */

import { useMemo } from 'react';

/**
 * 过滤行数据的Hook
 * @param {Array} rows - 原始数据行
 * @param {Array} fields - 要搜索的字段名数组
 * @param {string} search - 搜索关键词
 * @returns {Array} 过滤后的数据行
 */
export function useFilterRows(rows, fields, search) {
  return useMemo(() => {
    if (!search || !search.trim()) return rows;
    const kw = search.trim().toLowerCase();
    return rows.filter(row =>
      fields.some(f => (row[f] || '').toString().toLowerCase().includes(kw))
    );
  }, [rows, fields, search]);
}

/**
 * 深度克隆对象
 * @param {*} obj - 要克隆的对象
 * @returns {*} 克隆后的对象
 */
export function deepClone(obj) {
  if (obj === null || typeof obj !== 'object') return obj;
  if (obj instanceof Date) return new Date(obj.getTime());
  if (obj instanceof Array) return obj.map(item => deepClone(item));
  if (typeof obj === 'object') {
    const clonedObj = {};
    for (const key in obj) {
      if (obj.hasOwnProperty(key)) {
        clonedObj[key] = deepClone(obj[key]);
      }
    }
    return clonedObj;
  }
  return obj;
}

/**
 * 安全的JSON解析
 * @param {string} str - 要解析的JSON字符串
 * @param {*} defaultValue - 解析失败时的默认值
 * @returns {*} 解析后的对象或默认值
 */
export function safeJsonParse(str, defaultValue = null) {
  try {
    return JSON.parse(str);
  } catch (error) {
    console.warn('JSON解析失败:', error);
    return defaultValue;
  }
}

/**
 * 安全的JSON字符串化
 * @param {*} obj - 要字符串化的对象
 * @param {string} defaultValue - 字符串化失败时的默认值
 * @returns {string} 字符串化后的JSON字符串或默认值
 */
export function safeJsonStringify(obj, defaultValue = '{}') {
  try {
    return JSON.stringify(obj);
  } catch (error) {
    console.warn('JSON字符串化失败:', error);
    return defaultValue;
  }
}

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

/**
 * 节流函数
 * @param {Function} func - 要节流的函数
 * @param {number} limit - 限制时间（毫秒）
 * @returns {Function} 节流后的函数
 */
export function throttle(func, limit) {
  let inThrottle;
  return function executedFunction(...args) {
    if (!inThrottle) {
      func.apply(this, args);
      inThrottle = true;
      setTimeout(() => inThrottle = false, limit);
    }
  };
}
