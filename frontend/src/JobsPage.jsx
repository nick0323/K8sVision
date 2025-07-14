import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import { FaSync } from 'react-icons/fa';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import { SEARCH_PLACEHOLDER, PAGE_SIZE } from './constants';
import { useFilterRows } from './utils';
import Pagination from './Pagination';

export default function JobsPage() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const pageSize = PAGE_SIZE;
  const [pageMeta, setPageMeta] = useState({});

  const fetchData = useCallback(() => {
    setLoading(true);
    fetch(`/api/jobs?limit=${pageSize}&offset=${(page-1)*pageSize}`)
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

  const filteredRows = useFilterRows(data, ['namespace', 'name', 'completions', 'succeeded', 'failed', 'startTime', 'completionTime'], search);

  return (
    <div>
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'flex-end', alignItems: 'center', gap: 8 }}>
        <SearchInput
          value={search}
          onChange={e => setSearch(e.target.value)}
          placeholder={SEARCH_PLACEHOLDER}
        />
        <RefreshButton onClick={fetchData} />
      </div>
      <CommonTable
        columns={[
          { title: 'Name', dataIndex: 'name', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Namespace', dataIndex: 'namespace', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Completions', dataIndex: 'completions', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Succeeded', dataIndex: 'succeeded', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Failed', dataIndex: 'failed', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Start Time', dataIndex: 'startTime', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Completion Time', dataIndex: 'completionTime', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Status', dataIndex: 'status', render: (val, row, i, isTooltip) =>
              isTooltip
                ? val
                : <span className={`status-tag ${val === 'Succeeded' ? 'event-type-normal' : 'event-type-warning'}`} title={val}>{val}</span>
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
    </div>
  );
} 