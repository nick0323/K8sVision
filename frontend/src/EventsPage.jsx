import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import PageHeader from './components/PageHeader';
import { SEARCH_PLACEHOLDER, EMPTY_TEXT } from './constants';
import { usePagination } from './hooks/usePagination';
import { useSimpleSearch } from './hooks/useSimpleSearch';
import Pagination from './Pagination';

export default function EventsPage({ collapsed, onToggleCollapsed }) {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
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
    handleSearchChange,
    handleSearchSubmit,
    clearSearch,
    hasSearchResults,
    isSearchActive
  } = useSimpleSearch('/api/events', null);

  const fetchData = useCallback(() => {
    setLoading(true);
    const params = new URLSearchParams({
      limit: pageSize.toString(),
      offset: ((page-1)*pageSize).toString(),
    });
    
    fetch(`/api/events?${params}`)
      .then(res => res.json())
      .then(res => {
        // 确保数据始终是数组
        const dataList = res.data || res || [];
        setData(Array.isArray(dataList) ? dataList : []);
        setPageMeta(res.page || {});
      })
      .catch(error => {
        setData([]); // 出错时设置为空数组
        setPageMeta({});
      })
      .finally(() => setLoading(false));
  }, [page, pageSize]);

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

  // 确定要显示的数据
  const displayData = isSearchActive ? searchResults : data;
  const totalCount = isSearchActive ? searchTotal : (pageMeta?.total || 0);

  return (
    <div>
      <PageHeader
        title="Events"
        collapsed={collapsed}
        onToggleCollapsed={onToggleCollapsed}
      >
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
        className="events-table"
        columns={[
          { 
            title: 'Type', 
            dataIndex: 'type', 
            width: '100px',
            render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              const isNormal = val === 'Normal';
              const isWarning = val === 'Warning';
              let statusClass = 'status-running';
              if (isNormal) statusClass = 'status-ready';
              else if (isWarning) statusClass = 'status-failed';
              return <span className={`status-tag ${statusClass}`}>{val}</span>;
            }
          },
          { title: 'Reason', dataIndex: 'reason', width: '150px', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'Message', dataIndex: 'message', width: '300px', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { 
            title: 'Name', 
            dataIndex: 'name',
            width: '200px',
            render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              return <span>{val}</span>;
            }
          },
          { title: 'Namespace', dataIndex: 'namespace', width: '120px', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'FirstSeen', dataIndex: 'firstSeen', width: '150px', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'LastSeen', dataIndex: 'lastSeen', width: '150px', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'Duration', dataIndex: 'duration', width: '100px', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'Count', dataIndex: 'count', width: '80px', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
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

    </div>
  );
} 