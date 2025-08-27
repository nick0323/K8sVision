/**
 * 日期工具函数模块
 */

/**
 * 格式化日期时间为标准格式
 * @param {string} str - 日期字符串
 * @returns {string} 格式化后的日期时间字符串
 */
export function formatDateTime(str) {
  if (!str) return '';
  const d = new Date(str);
  if (isNaN(d.getTime())) return str;
  const pad = n => n < 10 ? '0' + n : n;
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}

/**
 * 格式化相对时间
 * @param {string} dateStr - 日期字符串
 * @returns {string} 相对时间字符串
 */
export function formatRelativeTime(dateStr) {
  if (!dateStr) return '';
  const now = new Date();
  const d = new Date(dateStr);
  const diff = Math.floor((now - d) / 1000);
  if (diff < 60) return `${diff} seconds ago`;
  if (diff < 3600) return `${Math.floor(diff/60)} minutes ago`;
  if (diff < 86400) return `${Math.floor(diff/3600)} hours ago`;
  return `${Math.floor(diff/86400)} days ago`;
}

/**
 * 计算两个时间之间的持续时间
 * @param {string} startTime - 开始时间
 * @param {string} endTime - 结束时间
 * @returns {string} 持续时间字符串
 */
export function calculateDuration(startTime, endTime) {
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
}
