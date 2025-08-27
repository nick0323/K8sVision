import { useState, useCallback } from 'react';
import { createSearchQuery, normalizeApiResponse } from '../utils/apiUtils';

/**
 * 简化的搜索Hook
 * 使用通用API工具避免复杂的依赖问题
 * @param {string} apiUrl - API端点URL
 * @param {string} namespace - 命名空间
 * @returns {Object} 搜索状态和操作方法
 */
export function useSimpleSearch(apiUrl, namespace) {
  const [search, setSearch] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [isSearching, setIsSearching] = useState(false);
  const [searchTotal, setSearchTotal] = useState(0);
  const [hasSearched, setHasSearched] = useState(false); // 标记是否已执行过搜索

  // 执行搜索
  const executeSearch = useCallback(async (searchTerm) => {
    if (!searchTerm || !searchTerm.trim()) {
      setSearchResults([]);
      setSearchTotal(0);
      setHasSearched(false);
      return;
    }

    setIsSearching(true);
    setHasSearched(true); // 标记已执行搜索
    
    try {
      const response = await createSearchQuery(apiUrl, searchTerm, { namespace });
      const { data, page } = normalizeApiResponse(response);
      
      setSearchResults(data);
      setSearchTotal(page?.total || data.length);
    } catch (error) {
      console.error('Search failed:', error);
      setSearchResults([]);
      setSearchTotal(0);
    } finally {
      setIsSearching(false);
    }
  }, [apiUrl, namespace]);

  // 处理搜索输入变化
  const handleSearchChange = useCallback((e) => {
    const value = e.target.value;
    setSearch(value);
  }, []);

  // 处理搜索提交
  const handleSearchSubmit = useCallback((e) => {
    e.preventDefault();
    if (search.trim()) {
      executeSearch(search);
    }
  }, [search, executeSearch]);

  // 清除搜索
  const clearSearch = useCallback(() => {
    setSearch('');
    setSearchResults([]);
    setSearchTotal(0);
    setHasSearched(false); // 重置搜索标记
  }, []);

  return {
    // 状态
    search,
    searchResults,
    isSearching,
    searchTotal,
    hasSearched, // 是否已执行过搜索
    
    // 操作方法
    setSearch,
    handleSearchChange,
    handleSearchSubmit,
    clearSearch,
    executeSearch,
    
    // 计算属性
    hasSearchResults: searchResults.length > 0,
    // 修复：只有当实际执行了搜索时才进入搜索状态
    isSearchActive: search.trim().length > 0 && hasSearched
  };
}
