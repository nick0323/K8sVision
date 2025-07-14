import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import { FaSync } from 'react-icons/fa';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
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

  const filteredRows = useFilterRows(data, ['namespace', 'name', 'class', 'hosts', 'path', 'targetService'], search);

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
          { title: 'Class', dataIndex: 'class', render: (val, row, i, isTooltip) => {
              const text = val ? val : '-';
              return isTooltip ? text : <span title={text}>{text}</span>;
            }
          },
          { title: 'Hosts', dataIndex: 'hosts', render: (val, row, i, isTooltip) => {
              const text = (Array.isArray(val) && val.length > 0) ? val.join(', ') : (val ? val : '-');
              return isTooltip ? text : <span title={text}>{text}</span>;
            }
          },
          { title: 'Path', dataIndex: 'path', render: (val, row, i, isTooltip) => {
              const text = Array.isArray(val) ? val.join(', ') : (val || '-');
              return isTooltip ? text : <span title={text}>{text}</span>;
            }
          },
          { title: 'Target Service', dataIndex: 'targetService', render: (val, row, i, isTooltip) => {
              const text = Array.isArray(val) ? val.join(', ') : (val || '-');
              return isTooltip ? text : <span title={text}>{text}</span>;
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
    </div>
  );
} 