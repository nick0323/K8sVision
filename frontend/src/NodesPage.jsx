import React, { useEffect, useState, useCallback } from 'react';
import CommonTable from './CommonTable';
import { FaSync } from 'react-icons/fa';
import RefreshButton from './components/RefreshButton';
import SearchInput from './components/SearchInput';
import { SEARCH_PLACEHOLDER, PAGE_SIZE } from './constants';
import { useFilterRows } from './utils';
import Pagination from './Pagination';

export default function NodesPage() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const pageSize = PAGE_SIZE;

  const fetchData = useCallback((page, pageSize) => {
    setLoading(true);
    fetch(`/api/nodes?limit=${pageSize}&offset=${(page-1)*pageSize}`)
      .then(res => res.json())
      .then(res => {
        setData(Array.isArray(res.data) ? res.data : []);
        setTotal(res.page?.total || 0);
      })
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    fetchData(page, pageSize);
  }, [fetchData, page, pageSize]);

  // 工具函数：单位转换
  function formatPercent(val) {
    if (val === null || val === undefined || isNaN(val)) return '-';
    return Number(val).toFixed(1) + '%';
  }
  function formatRole(roleArr) {
    if (!Array.isArray(roleArr) || roleArr.length === 0) return '-';
    return roleArr.join(', ');
  }

  const filteredRows = useFilterRows(data, ['name', 'ip', 'status', 'role'], search);

  return (
    <div>
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'flex-end', alignItems: 'center', gap: 8 }}>
        <SearchInput
          placeholder={SEARCH_PLACEHOLDER}
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        <RefreshButton onClick={() => fetchData(page, pageSize)} />
      </div>
      <CommonTable
        columns={[
          { title: 'Name', dataIndex: 'name', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'IP', dataIndex: 'ip', render: (val, row, i, isTooltip) => isTooltip ? val : <span title={val}>{val}</span> },
          { title: 'Role', dataIndex: 'role', render: (val, row, i, isTooltip) => {
              // 统一排序，优先 controlplane > etcd > worker
              const order = ['controlplane', 'etcd', 'worker'];
              let arr = Array.isArray(row.role) ? [...row.role] : (row.role ? [row.role] : []);
              arr = arr.sort((a, b) => {
                const ia = order.indexOf(a);
                const ib = order.indexOf(b);
                if (ia === -1 && ib === -1) return a.localeCompare(b);
                if (ia === -1) return 1;
                if (ib === -1) return -1;
                return ia - ib;
              });
              const text = arr.join(', ') || '-';
              return isTooltip ? text : <span title={text}>{text}</span>;
            }
          },
          { title: 'Status', dataIndex: 'status', render: (val, row, i, isTooltip) => isTooltip ? val : <span className={`status-tag ${(val === 'Ready' || val === 'Running' || val === 'Healthy' || val === 'Normal' || val === 'Active') ? 'event-type-normal' : 'event-type-warning'}`} title={val}>{val}</span> },
          { title: 'CPU Usage', dataIndex: 'cpuUsage', render: (val, row, i, isTooltip) => {
              const text = row.cpuUsage !== undefined && row.cpuUsage !== null ? Number(row.cpuUsage).toFixed(1) + '%' : '-';
              return isTooltip ? text : <span title={text}>{text}</span>;
            }
          },
          { title: 'Memory Usage', dataIndex: 'memoryUsage', render: (val, row, i, isTooltip) => {
              const text = row.memoryUsage !== undefined && row.memoryUsage !== null ? Number(row.memoryUsage).toFixed(1) + '%' : '-';
              return isTooltip ? text : <span title={text}>{text}</span>;
            }
          },
          { title: 'Pods', dataIndex: 'podsUsed', render: (val, row, i, isTooltip) => {
              const text = `${row.podsUsed ?? '-'} / ${row.podsCapacity ?? '-'}`;
              return isTooltip ? text : <span title={text}>{text}</span>;
            }
          },
        ]}
        data={filteredRows}
        pageSize={pageSize}
        currentPage={page}
        total={total || filteredRows.length}
        onPageChange={setPage}
        emptyText="No data"
      />
      <Pagination
        currentPage={page}
        total={total || filteredRows.length}
        pageSize={pageSize}
        onPageChange={setPage}
      />
    </div>
  );
} 