import { useState, useEffect, useMemo } from 'react';

export function formatDateTime(str) {
  if (!str) return '';
  const d = new Date(str);
  if (isNaN(d.getTime())) return str;
  const pad = n => n < 10 ? '0' + n : n;
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}

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

export function useFilterRows(rows, fields, search) {
  return useMemo(() => {
    if (!search || !search.trim()) return rows;
    const kw = search.trim().toLowerCase();
    return rows.filter(row =>
      fields.some(f => (row[f] || '').toString().toLowerCase().includes(kw))
    );
  }, [rows, fields, search]);
}

// 认证相关工具函数
export function getAuthHeaders() {
  const token = localStorage.getItem('token');
  if (!token) {
    return {};
  }
  
  // 验证token格式
  const segments = token.split('.');
  if (segments.length !== 3) {
    localStorage.removeItem('token');
    return {};
  }
  
  return {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  };
}

export function validateToken() {
  const token = localStorage.getItem('token');
  if (!token) return false;
  
  const segments = token.split('.');
  if (segments.length !== 3) {
    localStorage.removeItem('token');
    return false;
  }
  
  return true;
}

export function logout() {
  localStorage.removeItem('token');
  // 不清除 remembered_username 和 remembered_password，保持"记住我"功能
  // 触发页面刷新，让App组件重新检查登录状态
  window.location.reload();
} 