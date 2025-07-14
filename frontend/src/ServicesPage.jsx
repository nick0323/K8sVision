import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import { useFilterRows } from './utils';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import { SEARCH_PLACEHOLDER, PAGE_SIZE } from './constants';
import Pagination from './Pagination';

export default function ServicesPage() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const pageSize = PAGE_SIZE;
  const [pageMeta, setPageMeta] = useState({});

  const fetchData = useCallback(() => {
    setLoading(true);
    fetch(`/api/services?limit=${pageSize}&offset=${(page-1)*pageSize}`)
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
          { title: 'Type', dataIndex: 'type', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'ClusterIP', dataIndex: 'clusterIP', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Ports', dataIndex: 'ports', render: (val, row, i, isTooltip) => {
              const text = Array.isArray(val) ? val.join(', ') : val;
              return isTooltip ? text : <span title={text}>{text}</span>;
            }
          },
        ]}
        data={useFilterRows(data, ['namespace', 'name', 'type', 'clusterIP', 'ports'], search)}
        pageSize={pageSize}
        currentPage={page}
        onPageChange={setPage}
        total={pageMeta?.total || data.length}
        emptyText="No data"
      />
      <Pagination
        currentPage={page}
        total={pageMeta?.total || data.length}
        pageSize={pageSize}
        onPageChange={setPage}
      />
    </div>
  );
} 