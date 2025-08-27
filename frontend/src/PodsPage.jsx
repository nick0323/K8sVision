import React from 'react';
import ResourcePage from './components/ResourcePage';
import { PODS_CONFIG } from './constants/pageConfigs';

export default function PodsPage({ collapsed, onToggleCollapsed }) {
  return (
    <ResourcePage
      {...PODS_CONFIG}
      collapsed={collapsed}
      onToggleCollapsed={onToggleCollapsed}
    />
  );
} 