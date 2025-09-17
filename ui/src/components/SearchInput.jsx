import React, { useState } from 'react';

export default function SearchInput({ 
  value, 
  onChange, 
  onSubmit, 
  onClear,
  placeholder, 
  style,
  isSearching = false,
  hasSearchResults = false,
  showSearchButton = false,
  showClearButton = false
}) {
  const [isFocused, setIsFocused] = useState(false);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (onSubmit && value.trim()) {
      onSubmit(e);
    }
  };

  const handleClear = () => {
    if (onClear) {
      onClear();
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      handleSubmit(e);
    } else if (e.key === 'Escape') {
      handleClear();
    }
  };

  return (
    <div className="search-input-container" style={{ display: 'flex', alignItems: 'center', gap: '8px', ...style }}>
      <div className="search-input-wrapper" style={{ position: 'relative', display: 'flex', alignItems: 'center' }}>
        <input
          type="text"
          placeholder={placeholder}
          value={value}
          onChange={onChange}
          onKeyDown={handleKeyDown}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
          style={{ 
            width: 200, 
            height: 28, // 与namespace选择器高度保持一致
            fontSize: 'var(--font-size-sm)', 
            borderRadius: 6, 
            border: `1px solid ${isFocused ? '#1890ff' : '#d9d9d9'}`, 
            padding: '0 10px', // 移除左侧图标的内边距
            outline: 'none',
            transition: 'all 0.2s ease'
          }}
        />
        
        {/* 搜索状态指示器 */}
        {isSearching && (
          <div style={{
            position: 'absolute',
            right: '10px',
            color: '#1890ff',
            fontSize: '12px'
          }}>
            ⏳
          </div>
        )}

        {/* 搜索结果指示器 */}
        {hasSearchResults && !isSearching && (
          <div style={{
            position: 'absolute',
            right: '10px',
            color: '#52c41a',
            fontSize: '12px'
          }}>
            ✓
          </div>
        )}
      </div>
    </div>
  );
} 