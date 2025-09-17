import { useState, useCallback, useEffect, useRef } from 'react';
import { debounce } from '../utils/dataUtils';

/**
 * 全局搜索Hook
 * 支持搜索全部数据，而不仅仅是当前页面的数据
 * 新增防抖、缓存和错误重试机制
 * @param {string} apiUrl - API端点URL
 * @param {string} namespace - 命名空间
 * @param {Function} fetchData - 获取数据的函数
 * @param {Object} options - 配置选项
 * @returns {Object} 搜索状态和操作方法
 */
export function useGlobalSearch(apiUrl, namespace, fetchData, options = {}) {
  const {
    debounceMs = 300,           // 防抖延迟时间
    cacheTimeout = 5 * 60 * 1000, // 缓存超时时间（5分钟）
    maxRetries = 3,             // 最大重试次数
    searchLimit = 1000          // 搜索时获取的最大数据量
  } = options;

  const [search, setSearch] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [isSearching, setIsSearching] = useState(false);
  const [searchTotal, setSearchTotal] = useState(0);
  const [searchError, setSearchError] = useState(null);
  const [retryCount, setRetryCount] = useState(0);
  
  // 使用ref来避免初始化时的依赖问题
  const fetchDataRef = useRef(fetchData);
  fetchDataRef.current = fetchData;

  // 搜索缓存
  const searchCache = useRef(new Map());
  
  // 搜索历史
  const searchHistory = useRef([]);

  // 生成缓存键
  const getCacheKey = useCallback((searchTerm, namespace) => {
    return `${searchTerm}:${namespace || 'all'}`;
  }, []);

  // 检查缓存是否有效
  const isCacheValid = useCallback((cacheEntry) => {
    return cacheEntry && (Date.now() - cacheEntry.timestamp) < cacheTimeout;
  }, [cacheTimeout]);

  // 从缓存获取搜索结果
  const getFromCache = useCallback((searchTerm, namespace) => {
    const cacheKey = getCacheKey(searchTerm, namespace);
    const cacheEntry = searchCache.current.get(cacheKey);
    
    if (isCacheValid(cacheEntry)) {
      return cacheEntry.data;
    }
    
    // 清除过期缓存
    searchCache.current.delete(cacheKey);
    return null;
  }, [getCacheKey, isCacheValid]);

  // 保存到缓存
  const saveToCache = useCallback((searchTerm, namespace, data) => {
    const cacheKey = getCacheKey(searchTerm, namespace);
    searchCache.current.set(cacheKey, {
      data,
      timestamp: Date.now()
    });
    
    // 限制缓存大小，避免内存泄漏
    if (searchCache.current.size > 100) {
      const firstKey = searchCache.current.keys().next().value;
      searchCache.current.delete(firstKey);
    }
  }, [getCacheKey]);

  // 执行搜索
  const executeSearch = useCallback(async (searchTerm, retryAttempt = 0) => {
    if (!searchTerm || !searchTerm.trim()) {
      setSearchResults([]);
      setSearchTotal(0);
      setSearchError(null);
      return;
    }

    // 检查缓存
    const cachedResults = getFromCache(searchTerm, namespace);
    if (cachedResults) {
      setSearchResults(cachedResults.results);
      setSearchTotal(cachedResults.total);
      setSearchError(null);
      setIsSearching(false);
      return;
    }

    setIsSearching(true);
    setSearchError(null);
    
    try {
      // 构建搜索参数
      const params = new URLSearchParams({
        search: searchTerm.trim(),
        limit: searchLimit.toString(),
      });
      
      if (namespace) {
        params.append('namespace', namespace);
      }

      const response = await fetch(`${apiUrl}?${params}`);
      if (!response.ok) {
        throw new Error(`搜索请求失败: ${response.status} ${response.statusText}`);
      }

      const result = await response.json();
      const data = result.data || [];
      const total = result.page?.total || data.length;
      
      const searchData = {
        results: data,
        total
      };
      
      // 保存到缓存
      saveToCache(searchTerm, namespace, searchData);
      
      // 添加到搜索历史
      if (!searchHistory.current.includes(searchTerm)) {
        searchHistory.current.unshift(searchTerm);
        // 限制历史记录数量
        if (searchHistory.current.length > 20) {
          searchHistory.current.pop();
        }
      }
      
      setSearchResults(data);
      setSearchTotal(total);
      setRetryCount(0);
    } catch (error) {
      console.error('搜索失败:', error);
      setSearchError(error.message);
      
      // 重试机制
      if (retryAttempt < maxRetries) {
        setRetryCount(retryAttempt + 1);
        setTimeout(() => {
          executeSearch(searchTerm, retryAttempt + 1);
        }, 1000 * Math.pow(2, retryAttempt)); // 指数退避
      } else {
        setSearchResults([]);
        setSearchTotal(0);
      }
    } finally {
      setIsSearching(false);
    }
  }, [apiUrl, namespace, getFromCache, saveToCache, searchLimit, maxRetries]);

  // 防抖搜索
  const debouncedSearch = useCallback(
    debounce((searchTerm) => {
      executeSearch(searchTerm);
    }, debounceMs),
    [executeSearch, debounceMs]
  );

  // 处理搜索输入变化
  const handleSearchChange = useCallback((e) => {
    const value = e.target.value;
    setSearch(value);
    
    // 如果搜索框为空，清除搜索结果并重新获取原始数据
    if (!value.trim()) {
      setSearchResults([]);
      setSearchTotal(0);
      setSearchError(null);
      // 使用ref来安全调用fetchData
      if (fetchDataRef.current) {
        fetchDataRef.current();
      }
    } else {
      // 使用防抖搜索
      debouncedSearch(value);
    }
  }, [debouncedSearch]);

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
    setSearchError(null);
    setRetryCount(0);
    // 使用ref来安全调用fetchData
    if (fetchDataRef.current) {
      fetchDataRef.current();
    }
  }, []);

  // 重试搜索
  const retrySearch = useCallback(() => {
    if (search.trim()) {
      executeSearch(search);
    }
  }, [search, executeSearch]);

  // 获取搜索建议
  const getSearchSuggestions = useCallback(() => {
    return searchHistory.current.filter(term => 
      term.toLowerCase().includes(search.toLowerCase())
    ).slice(0, 5);
  }, [search]);

  // 当命名空间变化时，清除搜索
  useEffect(() => {
    if (search) {
      clearSearch();
    }
  }, [namespace, clearSearch]);

  // 清理防抖函数
  useEffect(() => {
    return () => {
      debouncedSearch.cancel && debouncedSearch.cancel();
    };
  }, [debouncedSearch]);

  return {
    // 状态
    search,
    searchResults,
    isSearching,
    searchTotal,
    searchError,
    retryCount,
    
    // 操作方法
    setSearch,
    handleSearchChange,
    handleSearchSubmit,
    clearSearch,
    executeSearch,
    retrySearch,
    
    // 计算属性
    hasSearchResults: searchResults.length > 0,
    isSearchActive: search.trim().length > 0,
    hasError: !!searchError,
    
    // 搜索建议
    searchSuggestions: getSearchSuggestions(),
    
    // 缓存管理
    clearCache: () => searchCache.current.clear(),
    getCacheSize: () => searchCache.current.size
  };
}
