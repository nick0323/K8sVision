import React from 'react';

export default function SearchInput({ value, onChange, placeholder, style }) {
  return (
    <input
      type="text"
      placeholder={placeholder}
      value={value}
      onChange={onChange}
      style={{ width: 200, height: 28, fontSize: 'var(--font-size-sm)', borderRadius: 6, border: '1px solid #d9d9d9', padding: '0 10px', ...style }}
    />
  );
} 