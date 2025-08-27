import React, { useState, useEffect } from 'react';

export default function NamespaceSelect({ value, onChange, placeholder = "All Namespaces" }) {
  const [namespaces, setNamespaces] = useState([]);
  const [loading, setLoading] = useState(false);

  // 获取namespace列表
  const fetchNamespaces = async () => {
    setLoading(true);
    try {
      const response = await fetch('/api/namespaces');
      const result = await response.json();
      // 确保数据始终是数组
      const namespaceList = result.data || result || [];
      setNamespaces(Array.isArray(namespaceList) ? namespaceList : []);
    } catch (error) {
      setNamespaces([]); // 出错时设置为空数组
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNamespaces();
  }, []);

  return (
    <div className="namespace-select">
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        disabled={loading}
        title={loading ? '正在加载命名空间...' : '选择命名空间'}
      >
        <option value="">{loading ? '加载中...' : placeholder}</option>
        {!loading && namespaces.map((ns) => (
          <option key={ns.name} value={ns.name}>
            {ns.name}
          </option>
        ))}
      </select>
    </div>
  );
} 