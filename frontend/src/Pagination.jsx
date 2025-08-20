import React from 'react';

export default function Pagination({ currentPage, total, pageSize, onPageChange }) {
  if (total <= pageSize) return null;
  const totalPages = Math.ceil(total / pageSize);
  
  return (
    <div className="table-pagination-area">
      {/* 左侧：总行数 */}
      <div className="pagination-total">
        {total} row(s) total
      </div>
      
      {/* 右侧：页码信息和导航按钮 */}
      <div className="pagination-controls">
        <span className="pagination-info">
          Page {currentPage} of {totalPages}
        </span>
        <button
          className="pagination-btn"
          onClick={() => currentPage > 1 && onPageChange(currentPage - 1)}
          disabled={currentPage === 1}
          title="Previous page"
        >
          ←
        </button>
        <button
          className="pagination-btn"
          onClick={() => currentPage < totalPages && onPageChange(currentPage + 1)}
          disabled={currentPage === totalPages}
          title="Next page"
        >
          →
        </button>
      </div>
    </div>
  );
} 