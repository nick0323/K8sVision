import React, { useState, lazy, Suspense } from 'react';
import './App.css';
import { FaServer, FaLayerGroup, FaDesktop, FaBell, FaCubes, FaTasks, FaChartPie, FaDatabase, FaShieldAlt, FaProjectDiagram, FaSync, FaNetworkWired, FaRegClock } from 'react-icons/fa';
import { LuSquareDashed } from 'react-icons/lu';
import { MENU_LIST, API_MAP } from './constants';

// 懒加载页面组件
const OverviewPage = lazy(() => import('./OverviewPage'));
const NodesPage = lazy(() => import('./NodesPage'));
const PodsPage = lazy(() => import('./PodsPage'));
const DeploymentsPage = lazy(() => import('./DeploymentsPage'));
const StatefulSetsPage = lazy(() => import('./StatefulSetsPage'));
const DaemonSetsPage = lazy(() => import('./DaemonSetsPage'));
const CronJobsPage = lazy(() => import('./CronJobsPage'));
const JobsPage = lazy(() => import('./JobsPage'));
const IngressPage = lazy(() => import('./IngressPage'));
const ServicesPage = lazy(() => import('./ServicesPage'));
const EventsPage = lazy(() => import('./EventsPage'));

// 图标映射
const ICON_MAP = {
  FaChartPie: <FaChartPie />,
  FaDesktop: <FaDesktop />,
  FaCubes: <FaCubes />,
  FaLayerGroup: <FaLayerGroup />,
  FaDatabase: <FaDatabase />,
  FaShieldAlt: <FaShieldAlt />,
  FaRegClock: <FaRegClock />,
  FaTasks: <FaTasks />,
  FaNetworkWired: <FaNetworkWired />,
  FaProjectDiagram: <FaProjectDiagram />,
  FaBell: <FaBell />,
};

// 页面组件映射
const PAGE_COMPONENTS = {
  overview: OverviewPage,
  nodes: NodesPage,
  pods: PodsPage,
  deployments: DeploymentsPage,
  statefulsets: StatefulSetsPage,
  daemonsets: DaemonSetsPage,
  cronjobs: CronJobsPage,
  jobs: JobsPage,
  ingress: IngressPage,
  services: ServicesPage,
  events: EventsPage,
};

export default function App() {
  const [tab, setTab] = useState('overview');

  const CurrentPage = PAGE_COMPONENTS[tab];

  return (
    <div className="layout-root">
      <div className="sider-menu">
        <ul>
          {MENU_LIST.map(item => (
            <li
              key={item.key}
              className={tab === item.key ? 'active' : ''}
              onClick={() => setTab(item.key)}
            >
              <span className="icon">{ICON_MAP[item.icon]}</span>
              <span>{item.label}</span>
            </li>
          ))}
        </ul>
      </div>
      <div className="main-content">
        <Suspense fallback={<div style={{textAlign:'center',color:'#888',padding:'32px 0'}}>加载中...</div>}>
          <CurrentPage />
        </Suspense>
      </div>
    </div>
  );
}