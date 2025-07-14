import React from 'react';
import { EMPTY_TEXT } from './constants';

function Tooltip({ children, content }) {
  const [show, setShow] = React.useState(false);
  const ref = React.useRef();
  // 只有内容溢出时才显示 Tooltip
  const isOverflow = () => {
    if (!ref.current) return false;
    return ref.current.scrollWidth > ref.current.clientWidth;
  };
  return (
    <span
      className="cell-ellipsis"
      style={{ position: 'relative', display: 'inline-block', verticalAlign: 'middle' }}
      onMouseEnter={() => isOverflow() && setShow(true)}
      onMouseLeave={() => setShow(false)}
    >
      <span ref={ref} style={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', display: 'inline-block' }}>
        {children}
      </span>
      {show && content && (
        <span style={{
          position: 'absolute',
          left: 0,
          top: '100%',
          zIndex: 99,
          background: '#222',
          color: '#fff',
          padding: '4px 10px',
          borderRadius: 4,
          fontSize: 13,
          whiteSpace: 'pre-wrap',
          marginTop: 4,
          maxWidth: 400,
          boxShadow: '0 2px 8px rgba(0,0,0,0.15)'
        }}>{content}</span>
      )}
    </span>
  );
}

export default function CommonTable({ columns = [], data = [], pageSize = 20, currentPage = 1, total = 0, onPageChange, emptyText = EMPTY_TEXT }) {
  const totalPages = Math.ceil(total / pageSize);

  return (
    <div className="table-card">
      <div className="table-content-area">
        <table className="table">
          <thead>
            <tr>
              {columns.map(col => (
                <th key={col.key || col.dataIndex}>{col.title}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {data.length === 0 ? (
              <tr>
                <td colSpan={columns.length} style={{ textAlign: 'center', color: '#c0c4cc', fontSize: 16 }}>
                  {emptyText}
                </td>
              </tr>
            ) : (
              data.map((row, i) => (
                <tr key={i}>
                  {columns.map(col => {
                    // 支持 render(text, row, i, isTooltip) 约定
                    const rawText = col.render ? col.render(row[col.dataIndex], row, i, true) : row[col.dataIndex];
                    const cellContent = col.render ? col.render(row[col.dataIndex], row, i, false) : row[col.dataIndex];
                    return (
                      <td key={col.key || col.dataIndex}>
                        <Tooltip content={rawText}>{cellContent}</Tooltip>
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