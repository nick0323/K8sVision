import React from 'react';

export default function SearchInput({ value, onChange, placeholder, style }) {
  return (
    <input
      type="text"
      value={value}
      onChange={onChange}
      placeholder={placeholder}
      style={{ width: 200, height: 28, fontSize: 15, borderRadius: 6, border: '1px solid #d9d9d9', padding: '0 10px', ...style }}
    />
  );
} 