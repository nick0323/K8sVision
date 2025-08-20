import React from 'react';
import { EMPTY_TEXT } from './constants';

export default function CommonTable({ columns = [], data = [], pageSize = 20, currentPage = 1, total = 0, onPageChange, emptyText = EMPTY_TEXT, className = '' }) {
  const totalPages = Math.ceil(total / pageSize);

  return (
    <div className="table-card">
      <div className="table-content-area">
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