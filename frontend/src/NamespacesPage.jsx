import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import PageHeader from './components/PageHeader';
import ResourceDetailModal from './components/ResourceDetailModal';
import { SEARCH_PLACEHOLDER, EMPTY_TEXT, PAGE_SIZE } from './constants';
import { useFilterRows } from './utils';
import Pagination from './Pagination';

export default function NamespacesPage({ collapsed, onToggleCollapsed }) {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const pageSize = PAGE_SIZE;
  const [pageMeta, setPageMeta] = useState({});

  // 详情模态框状态
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [currentResource, setCurrentResource] = useState(null);

  const fetchData = useCallback(() => {
    setLoading(true);
    fetch(`/api/namespaces?limit=${pageSize}&offset=${(page-1)*pageSize}`)
      .then(res => res.json())
      .then(res => {
        setData(res.data || res);
        setPageMeta(res.page || {});
      })
      .finally(() => setLoading(false));
  }, [page, pageSize]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const filteredRows = useFilterRows(data, ['name', 'status'], search);

  // 处理资源点击
  const handleResourceClick = (resource) => {
    setCurrentResource({
      type: 'namespaces',
      namespace: null, // 命名空间是集群级别的资源
      name: resource.name
    });
    setDetailModalVisible(true);
  };

  return (
    <div>
      <PageHeader
        title="Namespaces"
        onToggleCollapsed={onToggleCollapsed}
        collapsed={collapsed}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <SearchInput
            placeholder={SEARCH_PLACEHOLDER}
            value={search}
            onChange={e => setSearch(e.target.value)}
          />
          <RefreshButton onClick={fetchData} />
        </div>
      </PageHeader>
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
          { title: 'State', dataIndex: 'status', render: (val, row, i, isTooltip) => isTooltip ? val : <span className={`status-tag ${val === 'Active' ? 'status-ready' : 'status-running'}`}>{val}</span> },
        ]}
        data={filteredRows}
        pageSize={pageSize}
        currentPage={page}
        onPageChange={setPage}
        total={pageMeta?.total || filteredRows.length}
        emptyText={EMPTY_TEXT}
      />
      <Pagination
        currentPage={page}
        total={pageMeta?.total || filteredRows.length}
        pageSize={pageSize}
        onPageChange={setPage}
        fixedBottom={true}
      />

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