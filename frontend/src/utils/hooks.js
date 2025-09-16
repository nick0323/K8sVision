/**
 * 自定义Hooks工具模块
 */

import { useState, useEffect, useRef } from 'react';
import { getAuthHeaders, logout } from './authUtils';

/**
 * 通用数据获取Hook
 * @param {string} url - API URL
 * @param {Object} options - fetch选项
 * @returns {Object} 包含data, loading, error的对象
 */
export function useFetch(url, options = {}) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!url) return;
    let ignore = false;
    setLoading(true);
    setError(null);
    
    try {
      // 合并认证头
      const authHeaders = getAuthHeaders();
      const fetchOptions = {
        ...options,
        headers: {
          ...authHeaders,
          ...(options.headers || {})
        }
      };
      
      fetch(url, fetchOptions)
        .then(res => {
          if (!res.ok) {
            if (res.status === 401) {
              logout();
            }
            throw new Error(res.statusText || 'Network error');
          }
          return res.json();
        })
        .then(res => {
          if (!ignore) setData(res.data || res);
        })
        .catch(e => {
          if (!ignore) setError(e.message || '请求失败');
        })
        .finally(() => {
          if (!ignore) setLoading(false);
        });
    } catch (e) {
      if (!ignore) {
        setError('请求配置错误');
        setLoading(false);
      }
    }
    
    return () => { ignore = true; };
  }, [url, JSON.stringify(options)]);

  return { data, loading, error };
}

// 注意：已移除未使用的 hooks: useLocalStorage, useWindowSize, useClickOutside, useKeyPress
// 如需使用这些功能，请重新添加相应的 hook 实现
