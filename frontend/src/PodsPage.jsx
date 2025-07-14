import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import { FaSync } from 'react-icons/fa';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import { SEARCH_PLACEHOLDER, EMPTY_TEXT, PAGE_SIZE } from './constants';
import { useFilterRows } from './utils';
import Pagination from './Pagination';

export default function PodsPage() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const pageSize = PAGE_SIZE;
  const [pageMeta, setPageMeta] = useState({});

  const fetchData = useCallback(() => {
    setLoading(true);
    fetch(`/api/pods?limit=${pageSize}&offset=${(page-1)*pageSize}`)
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

  const filteredRows = useFilterRows(data, ['namespace', 'name', 'status', 'podIP', 'nodeName'], search);

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
          { title: 'Name', dataIndex: 'name', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Namespace', dataIndex: 'namespace', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Status', dataIndex: 'status', render: (val, row, i, isTooltip) =>
              isTooltip
                ? val
                : <span className={`status-tag ${(val === 'Running' || val === 'Succeeded') ? 'event-type-normal' : 'event-type-warning'}`} title={val}>{val}</span>
          },
          { title: 'CPU Usage', dataIndex: 'cpuUsage', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Memory Usage', dataIndex: 'memoryUsage', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Pod IP', dataIndex: 'podIP', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Node', dataIndex: 'nodeName', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
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
      />
    </div>
  );
} 