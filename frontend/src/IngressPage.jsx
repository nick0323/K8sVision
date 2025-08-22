import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import NamespaceSelect from './components/NamespaceSelect';
import PageHeader from './components/PageHeader';
import ResourceDetailModal from './components/ResourceDetailModal';
import { SEARCH_PLACEHOLDER, EMPTY_TEXT, PAGE_SIZE } from './constants';
import { useFilterRows } from './utils';
import Pagination from './Pagination';

export default function IngressPage({ collapsed, onToggleCollapsed }) {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [namespace, setNamespace] = useState('');
  const [page, setPage] = useState(1);
  const pageSize = PAGE_SIZE;
  const [pageMeta, setPageMeta] = useState({});

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
    
    fetch(`/api/ingress?${params}`)
      .then(res => res.json())
      .then(res => {
        // 确保数据始终是数组
        const dataList = res.data || res || [];
        setData(Array.isArray(dataList) ? dataList : []);
        setPageMeta(res.page || {});
      })
      .catch(error => {
        console.error('Failed to fetch Ingress:', error);
        setData([]); // 出错时设置为空数组
        setPageMeta({});
      })
      .finally(() => setLoading(false));
  }, [page, pageSize, namespace]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // 当namespace变化时重置到第一页
  useEffect(() => {
    setPage(1);
  }, [namespace]);

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
      <PageHeader
        title="Ingress"
        onToggleCollapsed={onToggleCollapsed}
        collapsed={collapsed}
      >
        <NamespaceSelect
          value={namespace}
          onChange={setNamespace}
          placeholder="All Namespaces"
        />
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