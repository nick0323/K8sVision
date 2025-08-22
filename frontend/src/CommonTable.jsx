import React from 'react';
import { EMPTY_TEXT } from './constants';

export default function CommonTable({ 
  columns = [], 
  data = [], 
  pageSize = 20, 
  currentPage = 1, 
  total = 0, 
  onPageChange, 
  emptyText = EMPTY_TEXT, 
  className = '',
  hasFixedPagination = false 
}) {
  const totalPages = Math.ceil(total / pageSize);
  
  // 根据是否有固定分页来决定表格卡片的类名
  const tableCardClass = hasFixedPagination ? 'table-card has-fixed-pagination' : 'table-card';
  
  // 根据数据量决定表格内容区域的类名 - 数据少时使用compact模式
  const isCompact = data.length <= 3; // 3行或更少时使用紧凑模式
  const tableContentClass = isCompact ? 'table-content-area compact' : 'table-content-area';

  return (
    <div className={tableCardClass}>
      <div className={tableContentClass}>
        <table className={`table ${className}`}>
          <thead>
            <tr>
              {columns.map(col => (
                <th key={col.key || col.dataIndex}>{col.label || col.title}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {data.length === 0 ? (
              <tr>
                <td colSpan={columns.length} style={{ textAlign: 'center', color: '#c0c4cc', fontSize: 'var(--font-size-sm)' }}>
                  {emptyText}
                </td>
              </tr>
            ) : (
              data.map((row, i) => (
                <tr key={i}>
                  {columns.map(col => {
                    // 直接渲染内容，不再需要Tooltip
                    const cellContent = col.render ? col.render(row[col.dataIndex], row, i, false) : row[col.dataIndex];
                    
                    return (
                      <td key={col.key || col.dataIndex}>
                        {cellContent}
                      </td>
                    );
                  })}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
      {/* 分页区域已移除，由外部页面负责渲染 */}
    </div>
  );
}