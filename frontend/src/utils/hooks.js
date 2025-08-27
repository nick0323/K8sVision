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

/**
 * 本地存储Hook
 * @param {string} key - 存储键名
 * @param {*} initialValue - 初始值
 * @returns {Array} [storedValue, setValue]
 */
export function useLocalStorage(key, initialValue) {
  const [storedValue, setStoredValue] = useState(() => {
    try {
      const item = window.localStorage.getItem(key);
      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      console.warn(`Error reading localStorage key "${key}":`, error);
      return initialValue;
    }
  });

  const setValue = (value) => {
    try {
      const valueToStore = value instanceof Function ? value(storedValue) : value;
      setStoredValue(valueToStore);
      window.localStorage.setItem(key, JSON.stringify(valueToStore));
    } catch (error) {
      console.warn(`Error setting localStorage key "${key}":`, error);
    }
  };

  return [storedValue, setValue];
}

/**
 * 窗口大小Hook
 * @returns {Object} 包含width和height的对象
 */
export function useWindowSize() {
  const [windowSize, setWindowSize] = useState({
    width: window.innerWidth,
    height: window.innerHeight,
  });

  useEffect(() => {
    function handleResize() {
      setWindowSize({
        width: window.innerWidth,
        height: window.innerHeight,
      });
    }

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  return windowSize;
}

/**
 * 点击外部Hook
 * @param {Function} callback - 点击外部时的回调函数
 * @returns {Object} ref对象
 */
export function useClickOutside(callback) {
  const ref = useRef();

  useEffect(() => {
    const listener = (event) => {
      if (!ref.current || ref.current.contains(event.target)) {
        return;
      }
      callback(event);
    };

    document.addEventListener('mousedown', listener);
    document.addEventListener('touchstart', listener);

    return () => {
      document.removeEventListener('mousedown', listener);
      document.removeEventListener('touchstart', listener);
    };
  }, [callback]);

  return ref;
}

/**
 * 键盘事件Hook
 * @param {string} targetKey - 目标按键
 * @param {Function} callback - 按键事件回调函数
 */
export function useKeyPress(targetKey, callback) {
  useEffect(() => {
    const keyPressHandler = (event) => {
      if (event.key === targetKey) {
        callback(event);
      }
    };

    document.addEventListener('keydown', keyPressHandler);
    return () => {
      document.removeEventListener('keydown', keyPressHandler);
    };
  }, [targetKey, callback]);
}
