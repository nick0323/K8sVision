import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import PageHeader from './components/PageHeader';
import ResourceDetailModal from './components/ResourceDetailModal';
import { SEARCH_PLACEHOLDER, EMPTY_TEXT, PAGE_SIZE } from './constants';
import { useFilterRows } from './utils';
import Pagination from './Pagination';

export default function NodesPage({ collapsed, onToggleCollapsed }) {
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
    fetch(`/api/nodes?limit=${pageSize}&offset=${(page-1)*pageSize}`)
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
      type: 'nodes',
      namespace: null, // 节点是集群级别的资源
      name: resource.name
    });
    setDetailModalVisible(true);
  };

  return (
    <div>
      <PageHeader
        title="Nodes"
        collapsed={collapsed}
        onToggleCollapsed={onToggleCollapsed}
      >
        <SearchInput
          placeholder={SEARCH_PLACEHOLDER}
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        <RefreshButton onClick={fetchData} />
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
          { title: 'IP', dataIndex: 'ip', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'Role', dataIndex: 'role', render: (val, row, i, isTooltip) => {
              if (isTooltip) return Array.isArray(val) ? val.join(', ') : (val || 'worker');
              const text = Array.isArray(val) ? val.join(', ') : (val || 'worker');
              return <span>{text}</span>;
            }
          },
          { title: 'CPUUsage', dataIndex: 'cpuUsage', render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              const roundedVal = typeof val === 'number' ? val.toFixed(2) : val;
              return <span>{roundedVal}%</span>;
            }
          },
          { title: 'MemoryUsage', dataIndex: 'memoryUsage', render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              const roundedVal = typeof val === 'number' ? val.toFixed(2) : val;
              return <span>{roundedVal}%</span>;
            }
          },
          { title: 'Pods', dataIndex: 'podsUsed', render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              const podsText = `${val || 0}`;
              return <span>{podsText}</span>;
            }
          },
          { title: 'State', dataIndex: 'status', render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              
              const isHealthy = val === 'Running' || val === 'Succeeded' || val === 'Ready' || val === 'Healthy' || val === 'Normal' || val === 'Active' || val === 'Bound';
              const isFailed = val === 'Failed' || val === 'Error' || val === 'CrashLoopBackOff';
              const isPending = val === 'Pending' || val === 'ContainerCreating' || val === 'PodInitializing';
              
              let statusClass = 'status-running';
              if (isHealthy) {
                statusClass = 'status-ready';
              } else if (isFailed) {
                statusClass = 'status-failed';
              } else if (isPending) {
                statusClass = 'status-pending';
              }
              
              return <span className={`status-tag ${statusClass}`}>{val}</span>;
            }
          },
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