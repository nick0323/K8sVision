import React, { useState, lazy, Suspense } from 'react';
import './App.css';
import { FaLayerGroup, FaDesktop, FaBell, FaCubes, FaTasks, FaChartPie, FaDatabase, FaShieldAlt, FaProjectDiagram, FaSync, FaNetworkWired, FaRegClock, FaChevronDown, FaChevronRight, FaChevronLeft, FaHdd, FaBoxOpen, FaBoxes, FaCog, FaKey } from 'react-icons/fa';
import { LuSquareDashed } from 'react-icons/lu';
import { MENU_LIST, API_MAP } from './constants';
import LoginPage from './LoginPage';
import { FiLogOut } from 'react-icons/fi';

// 懒加载页面组件
const OverviewPage = lazy(() => import('./OverviewPage'));
const NamespacesPage = lazy(() => import('./NamespacesPage'));
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
// 存储资源页面组件
const PVCsPage = lazy(() => import('./PVCsPage'));
const PVsPage = lazy(() => import('./PVsPage'));
const StorageClassesPage = lazy(() => import('./StorageClassesPage'));
// 配置资源页面组件
const ConfigMapsPage = lazy(() => import('./ConfigMapsPage'));
const SecretsPage = lazy(() => import('./SecretsPage'));

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
  FaHdd: <FaHdd />,
  FaBoxOpen: <FaBoxOpen />,
  FaBoxes: <FaBoxes />,
  FaCog: <FaCog />,
  FaKey: <FaKey />,
  LuSquareDashed: <LuSquareDashed />,
};

// 页面组件映射
const PAGE_COMPONENTS = {
  overview: OverviewPage,
  namespaces: NamespacesPage,
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
  // 存储资源页面组件
  pvcs: PVCsPage,
  pvs: PVsPage,
  storageclasses: StorageClassesPage,
  // 配置资源页面组件
  configmaps: ConfigMapsPage,
  secrets: SecretsPage,
};

function isLoggedIn() {
  return !!localStorage.getItem('token');
}

function setAuthHeader(onLogout) {
  const token = localStorage.getItem('token');
  console.log('设置认证头，token:', token ? 'exists' : 'not found');
  
  if (token) {
    // 检查token格式
    const segments = token.split('.');
    console.log('Token segments:', segments.length);
    if (segments.length !== 3) {
      console.warn('Token格式错误，清除token');
      localStorage.removeItem('token');
      if (onLogout) onLogout();
      return;
    }
    
    // 重置fetch拦截器
    if (window._originalFetch) {
      window.fetch = window._originalFetch;
    }
    
    window.fetch = ((origFetch) => async (url, options = {}) => {
      options.headers = options.headers || {};
      options.headers['Authorization'] = `Bearer ${token}`;
      console.log('发送请求:', url, '带Authorization头');
      
      const res = await origFetch(url, options);
      console.log('响应状态:', res.status);
      
      if (res.status === 401) {
        console.warn('收到401错误，清除token');
        localStorage.removeItem('token');
        if (onLogout) onLogout();
      }
      
      return res;
    })(window.fetch);
    
    window._fetchInterceptorSet = true;
    console.log('Fetch拦截器设置成功');
  } else {
    // 如果没有token，恢复原始fetch
    if (window._originalFetch) {
      window.fetch = window._originalFetch;
      window._fetchInterceptorSet = false;
      console.log('恢复原始fetch');
    }
  }
}

export default function App() {
  const [login, setLogin] = React.useState(isLoggedIn());
  const [tab, setTab] = useState('overview');
  // 分组折叠状态
  const [openGroups, setOpenGroups] = useState(() => {
    // 默认全部展开
    const state = {};
    MENU_LIST.forEach(g => { state[g.group] = true; });
    return state;
  });
  // 侧边栏是否折叠
  const [collapsed, setCollapsed] = useState(() => {
    try {
      return JSON.parse(localStorage.getItem('sider_collapsed') || 'false');
    } catch (e) { return false; }
  });
  const toggleCollapsed = React.useCallback(() => {
    setCollapsed(prev => {
      const next = !prev;
      localStorage.setItem('sider_collapsed', JSON.stringify(next));
      return next;
    });
  }, []);
  // 折叠态气泡提示
  const [tip, setTip] = useState({ visible: false, text: '', x: 0, y: 0 });
  const showTip = React.useCallback((e, text) => {
    if (!collapsed) return;
    const rect = e.currentTarget.getBoundingClientRect();
    setTip({
      visible: true,
      text,
      x: rect.right + 12,
      y: rect.top + rect.height / 2,
    });
  }, [collapsed]);
  const hideTip = React.useCallback(() => setTip(t => ({ ...t, visible: false })), []);
  
  React.useEffect(() => {
    // 保存原始fetch函数
    if (!window._originalFetch) {
      window._originalFetch = window.fetch;
    }
    
    // 初始化认证头
    if (isLoggedIn()) {
      setAuthHeader(() => {
        localStorage.removeItem('token');
        setLogin(false);
      });
    }
  }, []); 
  
  // 处理登录状态变化
  const handleLogin = React.useCallback(() => {
    setLogin(true);
    // 稍后设置认证头，确保状态已更新
    setTimeout(() => {
      setAuthHeader(() => {
        localStorage.removeItem('token');
        setLogin(false);
      });
    }, 100);
  }, []);
  
  // 退出登录按钮
  const handleLogout = React.useCallback(() => {
    localStorage.removeItem('token');
    setLogin(false);
  }, []);
  
  const toggleGroup = React.useCallback((group) => {
    setOpenGroups(prev => ({ ...prev, [group]: !prev[group] }));
  }, []);
  
  if (!login) {
    return <LoginPage onLogin={handleLogin} />;
  }
  
  const CurrentPage = PAGE_COMPONENTS[tab] || (() => <div style={{padding:32,textAlign:'center',color:'#888'}}>页面开发中</div>);

  return (
    <div className="layout-root" data-sider-collapsed={collapsed}>
      <div className={`sider-menu ${collapsed ? 'collapsed' : ''}`}>
        <div className="logo-area">
          <span className="logo-text-full" style={{fontSize: 'var(--font-size-2xl)', fontWeight: 700, color: '#2563eb', letterSpacing: '1px'}}>KubeVision</span>
          <span className="logo-text-compact" style={{fontSize: 'var(--font-size-xl)', fontWeight: 800, color: '#2563eb', letterSpacing: '1px'}}>KV</span>
        </div>
        <div className="sider-scroll">
          <ul>
            {/* Overview 单独渲染，放在分组菜单前，与分组一起滚动 */}
            {MENU_LIST[0].items.map(item => (
              <li
                key={item.key}
                className={tab === item.key ? 'active' : ''}
                onClick={() => setTab(item.key)}
                onMouseEnter={(e)=>showTip(e,item.label)}
                onMouseLeave={hideTip}
                data-tip={item.label}
              >
                <span className="icon">{ICON_MAP[item.icon]}</span>
                <span>{item.label}</span>
              </li>
            ))}
            {/* 其余分组可折叠渲染 */}
            {MENU_LIST.slice(1).map(group => (
              <React.Fragment key={group.group}>
                {!collapsed && (
                  <li
                    className="menu-group-title"
                    style={{cursor:'default',fontWeight:700,marginTop:16,display:'flex',alignItems:'center',userSelect:'none',fontSize:'var(--font-size-base)',color:'#888',justifyContent:'space-between'}}
                  >
                    <span>{group.group}</span>
                    <span
                      style={{marginLeft:8,cursor:'pointer',display:'flex',alignItems:'center'}}
                      onClick={e => { e.stopPropagation(); toggleGroup(group.group); }}
                    >
                      {openGroups[group.group] ? <FaChevronDown size={12}/> : <FaChevronRight size={12}/>} 
                    </span>
                  </li>
                )}
                {(collapsed || openGroups[group.group]) ? group.items.map(item => (
                  <li
                    key={item.key}
                    className={tab === item.key ? 'active' : ''}
                    onClick={() => setTab(item.key)}
                    onMouseEnter={(e)=>showTip(e,item.label)}
                    onMouseLeave={hideTip}
                    data-tip={item.label}
                  >
                    <span className="icon">{ICON_MAP[item.icon]}</span>
                    <span>{item.label}</span>
                  </li>
                )) : null}
              </React.Fragment>
            ))}
          </ul>
        </div>
        <div className="sider-bottom">
          <button className="logout-btn" onClick={handleLogout}>
            <span className="icon"><FiLogOut /></span>
            <span>Sign out</span>
          </button>
        </div>
      </div>
      <div className="main-content">
        <Suspense fallback={<div style={{textAlign:'center',color:'#888',padding:'32px 0'}}>Loading...</div>}>
          <CurrentPage collapsed={collapsed} onToggleCollapsed={toggleCollapsed} />
        </Suspense>
      </div>
      {collapsed && tip.visible && (
        <div style={{
          position: 'fixed',
          left: tip.x,
          top: tip.y,
          transform: 'translateY(-50%)',
          background: '#fff',
          color: '#222',
          padding: '8px 10px',
          borderRadius: 10,
          fontSize: 12,
          lineHeight: 1,
          boxShadow: '0 6px 20px rgba(0,0,0,0.10)',
          border: '1px solid rgba(226, 232, 240, 1)',
          pointerEvents: 'none',
          zIndex: 2000
        }}>{tip.text}</div>
      )}
    </div>
  );
}