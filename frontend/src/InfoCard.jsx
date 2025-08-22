import React from 'react';
import './App.css';

export default function InfoCard({ icon, title, value, status, children }) {
  return (
    <div className="overview-card">
      {icon && (
        <div className="overview-card-icon-col">
          <span className="overview-icon">{icon}</span>
        </div>
      )}
      <div className="overview-card-content-col" style={{minHeight: 110, display: 'flex', flexDirection: 'column', justifyContent: 'center'}}>
        <div className="overview-title">{title}</div>
        {value !== undefined && (
          <div className="overview-value">{value}</div>
        )}
        {status && (
          <div className="overview-status">{status}</div>
        )}
        {children}
      </div>
    </div>
  );
} 