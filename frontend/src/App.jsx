import React, { useState, lazy, Suspense } from 'react';
import './App.css';
import { FaServer, FaLayerGroup, FaDesktop, FaBell, FaCubes, FaTasks, FaChartPie, FaDatabase, FaShieldAlt, FaProjectDiagram, FaSync, FaNetworkWired, FaRegClock } from 'react-icons/fa';
import { LuSquareDashed } from 'react-icons/lu';
import { MENU_LIST, API_MAP } from './constants';
import LoginPage from './LoginPage';
import { FiLogOut } from 'react-icons/fi';

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

function isLoggedIn() {
  return !!localStorage.getItem('token');
}

function setAuthHeader(onLogout) {
  const token = localStorage.getItem('token');
  if (token) {
    window.fetch = ((origFetch) => async (url, options = {}) => {
      options.headers = options.headers || {};
      options.headers['Authorization'] = 'Bearer ' + token;
      const res = await origFetch(url, options);
      if (res.status === 401 && onLogout) {
        onLogout();
        return res;
      }
      return res;
    })(window.fetch);
  }
}

export default function App() {
  const [login, setLogin] = React.useState(isLoggedIn());
  React.useEffect(() => { setAuthHeader(() => {
    localStorage.removeItem('token');
    setLogin(false);
  }); }, [login]);
  if (!login) {
    return <LoginPage onLogin={() => setLogin(true)} />;
  }
  // 退出登录按钮
  const handleLogout = () => {
    localStorage.removeItem('token');
    setLogin(false);
    window.location.reload(); // 退出登录后强制刷新页面
  };
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
        <button className="logout-btn" onClick={handleLogout}>
          <span className="icon"><FiLogOut /></span>
          <span>Sign out</span>
        </button>
      </div>
      <div className="main-content">
        <Suspense fallback={<div style={{textAlign:'center',color:'#888',padding:'32px 0'}}>加载中...</div>}>
          <CurrentPage />
        </Suspense>
      </div>
    </div>
  );
}