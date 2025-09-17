import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from '../CommonTable';
import RefreshButton from './RefreshButton';
import SearchInput from './SearchInput';
import NamespaceSelect from './NamespaceSelect';
import PageHeader from './PageHeader';
import ResourceDetailModal from './ResourceDetailModal';
import { createNameRenderer } from '../utils/commonRenderers';
import { createPaginatedQuery, createSearchQuery, normalizeApiResponse } from '../utils/apiUtils';
import { SEARCH_PLACEHOLDER, EMPTY_TEXT } from '../constants';
import { usePagination } from '../hooks/usePagination';
import { useSimpleSearch } from '../hooks/useSimpleSearch';
import Pagination from '../Pagination';
import './ResourcePage.css';

/**
 * 通用资源页面组件
 * 消除所有页面组件的重复代码
 */
export default function ResourcePage({ 
  title, 
  apiEndpoint, 
  resourceType, 
  columns, 
  collapsed, 
  onToggleCollapsed,
  statusMap = {}, // 状态映射对象
  extraActions = null, // 额外的操作按钮
  searchPlaceholder = SEARCH_PLACEHOLDER,
  namespaceFilter = true // 是否需要namespace筛选器
}) {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [namespace, setNamespace] = useState('');
  const [pageMeta, setPageMeta] = useState({});

  // 使用分页Hook
  const {
    page,
    pageSize,
    handlePageChange,
    handlePageSizeChange,
    resetPagination,
    pageSizeOptions
  } = usePagination();

  // 使用简化的搜索Hook
  const {
    search,
    searchResults,
    isSearching,
    searchTotal,
    hasSearched,
    handleSearchChange,
    handleSearchSubmit,
    clearSearch,
    hasSearchResults,
    isSearchActive
  } = useSimpleSearch(apiEndpoint, namespaceFilter ? namespace : '');

  // 详情模态框状态
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [currentResource, setCurrentResource] = useState(null);

  // 使用通用API工具获取数据
  const fetchData = useCallback(async () => {
    setLoading(true);
    
    try {
      const response = await createPaginatedQuery(apiEndpoint, {
        page,
        pageSize,
        namespace: namespaceFilter ? namespace : ''
      });
      
      const { data: dataList, page: pageInfo } = normalizeApiResponse(response);
      setData(dataList);
      setPageMeta(pageInfo);
    } catch (error) {
      console.error('Failed to fetch data:', error);
      setData([]);
      setPageMeta({});
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, namespace, apiEndpoint, namespaceFilter]);

  // 增强的搜索输入变化处理，空值时重新获取数据
  const handleSearchInputChange = useCallback((e) => {
    const value = e.target.value;
    handleSearchChange(e);
    
    // 如果搜索框为空，清除搜索状态并重新获取原始数据
    if (!value.trim()) {
      clearSearch();
      fetchData();
    }
  }, [handleSearchChange, clearSearch, fetchData]);

  // 增强的清除搜索函数，会重新获取数据
  const handleClearSearch = useCallback(() => {
    clearSearch();
    fetchData();
  }, [clearSearch, fetchData]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // 当namespace变化时重置到第一页
  useEffect(() => {
    resetPagination();
  }, [namespace, resetPagination]);

  // 决定显示哪些数据：如果有搜索结果则显示搜索结果，否则显示原始数据
  const displayData = isSearchActive ? searchResults : data;
  const displayTotal = isSearchActive ? searchTotal : (pageMeta?.total || data.length);

  // 处理资源点击
  const handleResourceClick = (resource) => {
    setCurrentResource({
      type: resourceType,
      namespace: resource.namespace,
      name: resource.name
    });
    setDetailModalVisible(true);
  };

  // 处理列配置，添加默认的资源名称点击
  const processedColumns = columns.map(col => {
    if (col.dataIndex === 'name') {
      return {
        ...col,
        render: createNameRenderer(handleResourceClick)
      };
    }
    
    return col;
  });

  return (
    <div>
      <PageHeader
        title={title}
        collapsed={collapsed}
        onToggleCollapsed={onToggleCollapsed}
      >
        {/* 只有需要namespace筛选器的资源才显示namespace选择器 */}
        {namespaceFilter && (
          <NamespaceSelect
            value={namespace}
            onChange={setNamespace}
            placeholder="All Namespaces"
          />
        )}
        <SearchInput
          placeholder={searchPlaceholder}
          value={search}
          onChange={handleSearchInputChange}
          onSubmit={handleSearchSubmit}
          onClear={handleClearSearch}
          isSearching={isSearching}
          hasSearchResults={hasSearchResults}
          showSearchButton={false}
          showClearButton={false}
        />
        {extraActions}
        <RefreshButton onClick={fetchData} />
      </PageHeader>

      {/* 搜索状态提示 */}
      {isSearchActive && (
        <div className={`search-status-banner ${searchTotal > 0 ? 'has-results' : 'no-results'}`}>
          <div className="search-status-content">
            <div className="search-status-info">
              <span className="search-icon">🔍</span>
              <span>
                Search results for <strong>"{search}"</strong>: {
                  searchTotal > 0 
                    ? `Found ${searchTotal} matching resource${searchTotal > 1 ? 's' : ''}`
                    : 'No matching resources found'
                }
              </span>
              {searchTotal === 0 && (
                <span className="search-tip">
                  💡 Try different keywords or partial names
                </span>
              )}
            </div>
            <button
              onClick={handleClearSearch}
              className="clear-search-btn"
            >
              Clear Search
            </button>
          </div>
        </div>
      )}

      <CommonTable
        columns={processedColumns}
        data={displayData}
        pageSize={pageSize}
        currentPage={page}
        onPageChange={handlePageChange}
        total={displayTotal}
        emptyText={isSearchActive ? `No results found for "${search}"` : EMPTY_TEXT}
        hasFixedPagination={false}
      />
      
      {/* 只有在非搜索状态下才显示分页 */}
      {!isSearchActive && (
        <Pagination
          currentPage={page}
          total={displayTotal}
          pageSize={pageSize}
          onPageChange={handlePageChange}
          onPageSizeChange={handlePageSizeChange}
          pageSizeOptions={pageSizeOptions}
          fixedBottom={true}
        />
      )}

      {/* 资源详情模态框 */}
      <ResourceDetailModal
        visible={detailModalVisible}
        resourceType={currentResource?.type}
        namespace={currentResource?.namespace}
        name={currentResource?.name}
        onClose={() => setDetailModalVisible(false)}
      />
    </div>
  );
}
