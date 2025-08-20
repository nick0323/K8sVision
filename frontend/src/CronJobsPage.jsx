import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import { FaSync } from 'react-icons/fa';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import ResourceDetailModal from './components/ResourceDetailModal';
import { SEARCH_PLACEHOLDER, PAGE_SIZE } from './constants';
import { useFilterRows } from './utils';
import Pagination from './Pagination';

export default function CronJobsPage() {
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
    fetch(`/api/cronjobs?limit=${pageSize}&offset=${(page-1)*pageSize}`)
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

  const filteredRows = useFilterRows(data, ['namespace', 'name', 'status'], search);

  // 处理资源点击
  const handleResourceClick = (resource) => {
    setCurrentResource({
      type: 'cronjobs',
      namespace: resource.namespace,
      name: resource.name
    });
    setDetailModalVisible(true);
  };

  return (
    <div>
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'flex-end', alignItems: 'center', gap: 8 }}>
        <SearchInput
          placeholder={SEARCH_PLACEHOLDER}
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        <RefreshButton onClick={fetchData} />
      </div>
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
          { title: 'Schedule', dataIndex: 'schedule', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'Suspend', dataIndex: 'suspend', render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              return <span className={`status-tag ${val ? 'status-pending' : 'status-ready'}`}>
                {val ? 'Yes' : 'No'}
              </span>;
            }
          },
          { title: 'Active', dataIndex: 'active', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'LastSchedule', dataIndex: 'lastScheduleTime', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
          { title: 'Status', dataIndex: 'status', render: (val, row, i, isTooltip) => {
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
        emptyText="No data"
      />
      <Pagination
        currentPage={page}
        total={pageMeta?.total || filteredRows.length}
        pageSize={pageSize}
        onPageChange={setPage}
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