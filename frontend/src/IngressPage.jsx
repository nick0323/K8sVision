import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import { FaSync } from 'react-icons/fa';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import ResourceDetailModal from './components/ResourceDetailModal';
import { SEARCH_PLACEHOLDER, PAGE_SIZE } from './constants';
import { useFilterRows } from './utils';
import Pagination from './Pagination';

export default function IngressPage() {
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
    fetch(`/api/ingress?limit=${pageSize}&offset=${(page-1)*pageSize}`)
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
      type: 'ingress',
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
              const displayVal = val || '-';
              return (
                <span 
                  className="resource-name-link" 
                  onClick={() => handleResourceClick(row)}
                >
                  {displayVal}
                </span>
              );
            }
          },
          { title: 'Namespace', dataIndex: 'namespace', render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              const displayVal = val || '-';
              return <span>{displayVal}</span>;
            }
          },
          { title: 'Class', dataIndex: 'class', render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              const displayVal = val || '-';
              return <span>{displayVal}</span>;
            }
          },
          { title: 'Hosts', dataIndex: 'hosts', render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              const displayVal = Array.isArray(val) ? val.join(', ') : val;
              return <span>{displayVal}</span>;
            }
          },
          { title: 'Path', dataIndex: 'path', render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              const displayVal = Array.isArray(val) ? val.join(', ') : val;
              return <span>{displayVal || '-'}</span>;
            }
          },
          { title: 'TargetService', dataIndex: 'targetService', render: (val, row, i, isTooltip) => {
              if (isTooltip) return val;
              const text = Array.isArray(val) ? val.join(', ') : val;
              return <span>{text}</span>;
            }
          },
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