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

export function useFetch(url, options) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!url) return;
    let ignore = false;
    setLoading(true);
    setError(null);
    fetch(url, options)
      .then(res => {
        if (!res.ok) throw new Error(res.statusText || 'Network error');
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