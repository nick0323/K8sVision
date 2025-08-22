import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import { formatDateTime } from './utils';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import NamespaceSelect from './components/NamespaceSelect';
import PageHeader from './components/PageHeader';
import { SEARCH_PLACEHOLDER } from './constants';
import { PAGE_SIZE } from './constants';
import { useFilterRows } from './utils';
import Pagination from './Pagination';

export default function EventsPage({ collapsed, onToggleCollapsed }) {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [namespace, setNamespace] = useState('');
  const [page, setPage] = useState(1);
  const pageSize = PAGE_SIZE;
  const [pageMeta, setPageMeta] = useState({});

  const fetchData = useCallback(() => {
    setLoading(true);
    const params = new URLSearchParams({
      limit: pageSize.toString(),
      offset: ((page-1)*pageSize).toString(),
    });
    
    if (namespace) {
      params.append('namespace', namespace);
    }
    
    fetch(`/api/events?${params}`)
      .then(res => res.json())
      .then(res => {
        // 确保数据始终是数组
        const dataList = res.data || res || [];
        setData(Array.isArray(dataList) ? dataList : []);
        setPageMeta(res.page || {});
      })
      .catch(error => {
        console.error('Failed to fetch Events:', error);
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

  const filteredRows = useFilterRows(data, ['namespace', 'name', 'reason', 'message', 'type', 'firstSeen', 'lastSeen', 'count'], search);

  return (
    <div>
      <PageHeader
        title="Events"
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
          onChange={e => setSearch(e.target.value)}
        />
        <RefreshButton onClick={fetchData} />
      </PageHeader>
      <CommonTable
        className="events-table"
        columns={[
          { title: 'Type', dataIndex: 'type', render: (val, row, i, isTooltip) => {
            if (isTooltip) return val;
            const eventTypeClass = val === 'Warning' ? 'event-type-warning' : 'event-type-normal';
            return <span className={`event-type ${eventTypeClass}`}>{val}</span>;
          }
        },
        { title: 'Reason', dataIndex: 'reason', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
        { title: 'Message', dataIndex: 'message', render: (val, row, i, isTooltip) => isTooltip ? val : <span className="event-message">{val}</span> },
        { title: 'Name', dataIndex: 'name', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
        { title: 'Namespace', dataIndex: 'namespace', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
        { title: 'FirstSeen', dataIndex: 'firstSeen', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
        { title: 'LastSeen', dataIndex: 'lastSeen', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
        { title: 'Duration', dataIndex: 'duration', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
        { title: 'Count', dataIndex: 'count', render: (val, row, i, isTooltip) => isTooltip ? val : <span>{val}</span> },
        ]}
        data={filteredRows.slice().sort((a, b) => new Date(b.lastSeen) - new Date(a.lastSeen))}
        pageSize={pageSize}
        currentPage={page}
        onPageChange={setPage}
        total={pageMeta?.total || filteredRows.length}
        emptyText="No events"
        hasFixedPagination={false}
      />
      <Pagination
        currentPage={page}
        total={pageMeta?.total || filteredRows.length}
        pageSize={pageSize}
        onPageChange={setPage}
        fixedBottom={true}
      />
    </div>
  );
} 