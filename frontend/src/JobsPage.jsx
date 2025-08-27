import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import NamespaceSelect from './components/NamespaceSelect';
import PageHeader from './components/PageHeader';
import ResourceDetailModal from './components/ResourceDetailModal';
import { SEARCH_PLACEHOLDER, EMPTY_TEXT, PAGE_SIZE } from './constants';
import Pagination from './Pagination';
import { usePagination } from './hooks/usePagination';
import { useSimpleSearch } from './hooks/useSimpleSearch';

export default function JobsPage({ collapsed, onToggleCollapsed }) {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [namespace, setNamespace] = useState('');
  const [pageMeta, setPageMeta] = useState({});

  // 使用分页hook
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
  } = useSimpleSearch('/api/jobs', namespace);

  // 详情模态框状态
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [currentResource, setCurrentResource] = useState(null);

  const fetchData = useCallback(() => {
    setLoading(true);
    const params = new URLSearchParams({
      limit: pageSize.toString(),
      offset: ((page-1)*pageSize).toString(),
    });
    
    if (namespace) {
      params.append('namespace', namespace);
    }
    
    fetch(`/api/jobs?${params}`)
      .then(res => res.json())
      .then(res => {
        // 确保数据始终是数组
        const dataList = res.data || res || [];
        setData(Array.isArray(dataList) ? dataList : []);
        setPageMeta(res.page || {});
      })
      .catch(error => {
        console.error('Failed to fetch Jobs:', error);
        setData([]); // 出错时设置为空数组
        setPageMeta({});
      })
      .finally(() => setLoading(false));
  }, [page, pageSize, namespace]);

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

  // 处理资源点击
  const handleResourceClick = (resource) => {
    setCurrentResource({
      type: 'jobs',
      namespace: resource.namespace,
      name: resource.name
    });
    setDetailModalVisible(true);
  };

  // 确定要显示的数据
  const displayData = isSearchActive ? searchResults : data;
  const totalCount = isSearchActive ? searchTotal : (pageMeta?.total || 0);

  return (
    <div>
      <PageHeader
        title="Jobs"
        collapsed={collapsed}
        onToggleCollapsed={onToggleCollapsed}
      >
        <NamespaceSelect
          value={namespace}
          onChange={setNamespace}
          placeholder="All Namespaces"
        />
        <SearchInput
          placeholder={SEARCH_PLACEHOLDER}
          value={search}
          onChange={handleSearchInputChange}
          onSubmit={handleSearchSubmit}
          onClear={handleClearSearch}
          isSearching={isSearching}
          hasSearchResults={hasSearchResults}
          showSearchButton={false}
          showClearButton={false}
        />
        <RefreshButton onClick={fetchData} />
      </PageHeader>

      {/* 搜索状态提示 */}
      {isSearchActive && (
        <div style={{
          padding: '12px 16px',
          margin: '8px 0',
          background: searchTotal > 0 ? '#f0f8ff' : '#fff7e6',
          border: `1px solid ${searchTotal > 0 ? '#d6e4ff' : '#ffd591'}`,
          borderRadius: '6px',
          fontSize: 'var(--font-size-sm)',
          color: searchTotal > 0 ? '#1890ff' : '#d46b08'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <span>🔍</span>
              <span>
                Search results for <strong>"{search}"</strong>: {
                  searchTotal > 0 
                    ? `Found ${searchTotal} matching resource${searchTotal > 1 ? 's' : ''}`
                    : 'No matching resources found'
                }
              </span>
              {searchTotal === 0 && (
                <span style={{ fontSize: 'var(--font-size-xs)', opacity: 0.8 }}>
                  💡 Try different keywords or partial names
                </span>
              )}
            </div>
            <button
              onClick={handleClearSearch}
              style={{
                padding: '4px 12px',
                fontSize: 'var(--font-size-xs)',
                border: `1px solid ${searchTotal > 0 ? '#1890ff' : '#ffd591'}`,
                borderRadius: '4px',
                background: 'white',
                color: searchTotal > 0 ? '#1890ff' : '#d46b08',
                cursor: 'pointer',
                transition: 'all 0.2s ease'
              }}
              onMouseEnter={(e) => {
                e.target.style.background = searchTotal > 0 ? '#1890ff' : '#d46b08';
                e.target.style.color = 'white';
              }}
              onMouseLeave={(e) => {
                e.target.style.background = 'white';
                e.target.style.color = searchTotal > 0 ? '#1890ff' : '#d46b08';
              }}
            >
              Clear Search
            </button>
          </div>
        </div>
      )}

      <CommonTable
        columns={[
          { 
            title: 'Name', 
            dataIndex: 'name', 
            render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              return (
                <span 
                  className="resource-name-link" 
                  onClick={() => handleResourceClick(row)}
                >
                  {val}
                </span>
              );
            }
          },
          { title: 'Namespace', dataIndex: 'namespace', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'Completions', dataIndex: 'completions', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'Succeeded', dataIndex: 'succeeded', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'Failed', dataIndex: 'failed', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'StartTime', dataIndex: 'startTime', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'CompletionTime', dataIndex: 'completionTime', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'State', dataIndex: 'status', render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              
              // Job状态判断逻辑
              const isCompleted = val === 'Completed' || val === 'Succeeded';
              const isRunning = val === 'Running' || val === 'Active';
              const isFailed = val === 'Failed' || val === 'Error';
              const isPending = val === 'Pending' || val === 'Creating';
              
              let statusClass = 'status-running';
              if (isCompleted) {
                statusClass = 'status-ready';
              } else if (isRunning) {
                statusClass = 'status-running';
              } else if (isFailed) {
                statusClass = 'status-failed';
              } else if (isPending) {
                statusClass = 'status-pending';
              }
              
              return <span className={`status-tag ${statusClass}`}>{val}</span>;
            } },
        ]}
        data={displayData}
        pageSize={pageSize}
        currentPage={page}
        onPageChange={handlePageChange}
        total={totalCount}
        emptyText={isSearchActive ? `No results found for "${search}"` : EMPTY_TEXT}
        hasFixedPagination={false}
      />

      {/* 只在非搜索状态下显示分页 */}
      {!isSearchActive && (
        <Pagination
          currentPage={page}
          total={totalCount}
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