import React from 'react';
import ResourcePage from './components/ResourcePage';
import { DEPLOYMENTS_CONFIG } from './constants/pageConfigs';

export default function DeploymentsPage({ collapsed, onToggleCollapsed }) {
  return (
    <ResourcePage
      {...DEPLOYMENTS_CONFIG}
      collapsed={collapsed}
      onToggleCollapsed={onToggleCollapsed}
    />
  );
} 