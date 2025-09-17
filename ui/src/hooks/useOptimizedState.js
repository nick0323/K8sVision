import { useState, useCallback, useRef, useEffect } from 'react';

/**
 * 优化的状态管理Hook
 * 提供防抖、节流、缓存等功能
 */
export const useOptimizedState = (initialValue, options = {}) => {
  const {
    debounceMs = 0,
    throttleMs = 0,
    cache = false,
    cacheKey = null,
    onUpdate = null
  } = options;

  // 缓存相关
  const cacheRef = useRef(new Map());
  
  // 获取缓存值
  const getCachedValue = useCallback((key) => {
    if (!cache || !key) return null;
    const cached = cacheRef.current.get(key);
    if (cached && Date.now() - cached.timestamp < 300000) { // 5分钟过期
      return cached.value;
    }
    return null;
  }, [cache]);

  // 设置缓存值
  const setCachedValue = useCallback((key, value) => {
    if (!cache || !key) return;
    cacheRef.current.set(key, {
      value,
      timestamp: Date.now()
    });
  }, [cache]);

  // 初始化状态
  const getInitialValue = useCallback(() => {
    if (typeof initialValue === 'function') {
      return initialValue();
    }
    return initialValue;
  }, [initialValue]);

  const [state, setState] = useState(() => {
    const cached = getCachedValue(cacheKey);
    if (cached !== null) {
      return cached;
    }
    return getInitialValue();
  });

  // 防抖和节流相关
  const debounceRef = useRef(null);
  const throttleRef = useRef(null);
  const lastUpdateRef = useRef(0);

  // 更新状态的核心函数
  const updateState = useCallback((newValue) => {
    const updateFn = (value) => {
      setState(value);
      if (cache && cacheKey) {
        setCachedValue(cacheKey, value);
      }
      if (onUpdate) {
        onUpdate(value);
      }
    };

    // 如果设置了防抖
    if (debounceMs > 0) {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
      debounceRef.current = setTimeout(() => {
        updateFn(newValue);
      }, debounceMs);
      return;
    }

    // 如果设置了节流
    if (throttleMs > 0) {
      const now = Date.now();
      if (now - lastUpdateRef.current < throttleMs) {
        return;
      }
      lastUpdateRef.current = now;
    }

    updateFn(newValue);
  }, [debounceMs, throttleMs, cache, cacheKey, onUpdate, getCachedValue, setCachedValue]);

  // 清理定时器
  useEffect(() => {
    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, []);

  return [state, updateState];
};