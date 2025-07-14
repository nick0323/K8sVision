import React from 'react';
import { FaSync } from 'react-icons/fa';

export default function RefreshButton({ onClick, title = '刷新', style }) {
  return (
    <button
      className="btn-refresh"
      onClick={onClick}
      title={title}
      style={style}
    >
      <FaSync />
    </button>
  );
} 