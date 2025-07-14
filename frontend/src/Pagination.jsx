import React from 'react';

export default function Pagination({ currentPage, total, pageSize, onPageChange }) {
  if (total <= pageSize) return null;
  const totalPages = Math.ceil(total / pageSize);
  const pages = [];
  if (totalPages <= 7) {
    for (let i = 1; i <= totalPages; i++) pages.push(i);
  } else {
    pages.push(1);
    if (currentPage > 4) pages.push('...');
    for (let i = Math.max(2, currentPage - 2); i <= Math.min(totalPages - 1, currentPage + 2); i++) {
      pages.push(i);
    }
    if (currentPage < totalPages - 3) pages.push('...');
    pages.push(totalPages);
  }
  return (
    <div className="table-pagination-area">
      <button
        onClick={() => currentPage > 1 && onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
      >上一页</button>
      {pages.map((p, idx) =>
        p === '...'
          ? <span key={"ellipsis-"+idx} style={{margin: '0 4px', color: '#bbb', userSelect: 'none'}}>...</span>
          : <button
              key={p}
              onClick={() => onPageChange(p)}
              style={{
                margin: '0 4px',
                padding: '2px 10px',
                borderRadius: 4,
                border: currentPage === p ? '1.5px solid #1890ff' : '1px solid #d9d9d9',
                background: currentPage === p ? '#e6f7ff' : '#fff',
                color: currentPage === p ? '#1890ff' : '#333',
                cursor: 'pointer',
                fontWeight: currentPage === p ? 600 : 400
              }}
            >{p}</button>
      )}
      <button
        onClick={() => currentPage < totalPages && onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
      >下一页</button>
    </div>
  );
} 