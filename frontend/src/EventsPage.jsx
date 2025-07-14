import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import { formatDateTime } from './utils';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import { SEARCH_PLACEHOLDER } from './constants';
import { PAGE_SIZE } from './constants';
import { useFilterRows } from './utils';
import Pagination from './Pagination';

export default function EventsPage() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const pageSize = PAGE_SIZE;
  const [pageMeta, setPageMeta] = useState({});

  const fetchData = useCallback(() => {
    setLoading(true);
    fetch(`/api/events?limit=${pageSize}&offset=${(page-1)*pageSize}`)
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

  const filteredRows = useFilterRows(data, ['namespace', 'name', 'reason', 'message', 'type', 'firstSeen', 'lastSeen', 'count'], search);

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
          { title: 'Type', dataIndex: 'type', render: (val, row, i, isTooltip) => isTooltip ? val : <span className={`status-tag ${val === 'Warning' ? 'event-type-warning' : 'event-type-normal'}`}>{val}</span> },
          { title: 'Reason', dataIndex: 'reason', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Message', dataIndex: 'message', render: (val, row, i, isTooltip) => isTooltip ? val : <span className="event-message" title={val}>{val}</span> },
          { title: 'Name', dataIndex: 'name', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Namespace', dataIndex: 'namespace', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'First Seen', dataIndex: 'firstSeen', width: 140, render: (val, row, i, isTooltip) => isTooltip ? formatDateTime(val) : <span className="ellipsis-cell">{formatDateTime(val)}</span> },
          { title: 'Last Seen', dataIndex: 'lastSeen', width: 140, render: (val, row, i, isTooltip) => isTooltip ? formatDateTime(val) : <span className="ellipsis-cell">{formatDateTime(val)}</span> },
          { title: 'Count', dataIndex: 'count', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
        ]}
        data={filteredRows.slice().sort((a, b) => new Date(b.lastSeen) - new Date(a.lastSeen))}
        pageSize={pageSize}
        currentPage={page}
        onPageChange={setPage}
        total={pageMeta?.total || filteredRows.length}
        emptyText="No events"
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