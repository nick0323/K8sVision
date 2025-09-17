import React, { useState, lazy, Suspense, useCallback, useMemo, useEffect } from 'react';
import './App.css';
import LoadingSpinner from './components/LoadingSpinner';
import ErrorBoundary from './components/ErrorBoundary';
import { useOptimizedState } from './hooks/useOptimizedState';
import { FaLayerGroup, FaDesktop, FaBell, FaCubes, FaTasks, FaChartPie, FaDatabase, FaShieldAlt, FaProjectDiagram, FaSync, FaNetworkWired, FaRegClock, FaChevronDown, FaChevronRight, FaChevronLeft, FaHdd, FaBoxOpen, FaBoxes, FaCog, FaKey } from 'react-icons/fa';
import { LuSquareDashed } from 'react-icons/lu';
import { MENU_LIST, API_MAP } from './constants';
import LoginPage from './LoginPage';
import { FiLogOut } from 'react-icons/fi';

// 懒加载页面组件 - 使用统一的页面组件定义
const OverviewPage = lazy(() => import('./OverviewPage'));

// 导入页面组件
import { PAGE_COMPONENTS } from './pages.jsx';

// 解构页面组件
const {
  pods: PodsPage,
  deployments: DeploymentsPage,
  statefulsets: StatefulSetsPage,
  daemonsets: DaemonSetsPage,
  cronjobs: CronJobsPage,
  jobs: JobsPage,
  ingress: IngressPage,
  services: ServicesPage,
  events: EventsPage,
  pvcs: PVCsPage,
  pvs: PVsPage,
  storageclasses: StorageClassesPage,
  configmaps: ConfigMapsPage,
  secrets: SecretsPage,
  namespaces: NamespacesPage,
  nodes: NodesPage
} = PAGE_COMPONENTS;

// 图标映射 - 使用useMemo优化
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

// 页面组件映射 - 使用统一的页面组件
const CURRENT_PAGE_COMPONENTS = {
  ...PAGE_COMPONENTS
};

// 预加载策略：当用户hover到菜单项时预加载对应组件
const preloadComponent = (componentKey) => {
  const Component = CURRENT_PAGE_COMPONENTS[componentKey];
  if (Component && Component.preload) {
    Component.preload();
  }
};

// 认证相关工具函数
const authUtils = {
  isLoggedIn: () => !!localStorage.getItem('token'),
  
  setAuthHeader: (onLogout) => {
    const token = localStorage.getItem('token');
    
    if (token) {
      // 检查token格式
      const segments = token.split('.');
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
        
        const res = await origFetch(url, options);
        
        if (res.status === 401) {
          console.warn('收到401错误，清除token');
          localStorage.removeItem('token');
          if (onLogout) onLogout();
        }
        
        return res;
      })(window.fetch);
      
      window._fetchInterceptorSet = true;
    } else {
      // 如果没有token，恢复原始fetch
      if (window._originalFetch) {
        window.fetch = window._originalFetch;
        window._fetchInterceptorSet = false;
      }
    }
  }
};

export default function App() {
  const [login, setLogin] = React.useState(authUtils.isLoggedIn());
  const [tab, setTab] = useOptimizedState('overview', {
    cache: true,
    cacheKey: 'current_tab',
  });
  
  // 分组折叠状态 - 使用优化的状态管理
  const [openGroups, setOpenGroups] = useOptimizedState((() => {
    const state = {};
    MENU_LIST.forEach(g => { state[g.group] = true; });
    return state;
  })(), {
    cache: true,
    cacheKey: 'open_groups',
  });
  
  // 侧边栏是否折叠 - 使用优化的状态管理
  const [collapsed, setCollapsed] = useOptimizedState((() => {
    try {
      return JSON.parse(localStorage.getItem('sider_collapsed') || 'false');
    } catch (e) { return false; }
  })(), {
    cache: true,
    cacheKey: 'sider_collapsed',
  });
  
  // 使用useCallback优化函数，避免不必要的重渲染
  const toggleCollapsed = useCallback(() => {
    setCollapsed(prev => {
      const next = !prev;
      localStorage.setItem('sider_collapsed', JSON.stringify(next));
      return next;
    });
  }, []);
  
  // 折叠态气泡提示
  const [tip, setTip] = useState({ visible: false, text: '', x: 0, y: 0 });
  
  const showTip = useCallback((e, text) => {
    if (!collapsed) return;
    const rect = e.currentTarget.getBoundingClientRect();
    setTip({
      visible: true,
      text,
      x: rect.right + 12,
      y: rect.top + rect.height / 2,
    });
  }, [collapsed]);
  
  const hideTip = useCallback(() => setTip(t => ({ ...t, visible: false })), []);
  
  // 处理登录状态变化
  const handleLogin = useCallback(() => {
    setLogin(true);
    // 稍后设置认证头，确保状态已更新
    setTimeout(() => {
      authUtils.setAuthHeader(() => {
        localStorage.removeItem('token');
        setLogin(false);
      });
    }, 100);
  }, []);
  
  // 退出登录按钮
  const handleLogout = useCallback(() => {
    localStorage.removeItem('token');
    setLogin(false);
  }, []);
  
  const toggleGroup = useCallback((group) => {
    setOpenGroups(prev => ({ ...prev, [group]: !prev[group] }));
  }, []);
  
  // 处理标签页切换，添加预加载
  const handleTabChange = useCallback((newTab) => {
    setTab(newTab);
    // 预加载相邻的组件
    const menuItems = MENU_LIST.flatMap(g => g.items);
    const currentIndex = menuItems.findIndex(item => item.key === newTab);
    if (currentIndex !== -1) {
      // 预加载下一个组件
      if (currentIndex + 1 < menuItems.length) {
        preloadComponent(menuItems[currentIndex + 1].key);
      }
      // 预加载上一个组件
      if (currentIndex - 1 >= 0) {
        preloadComponent(menuItems[currentIndex - 1].key);
      }
    }
  }, []);
  
  // 使用useMemo优化菜单渲染
  const menuItemsList = useMemo(() => {
    return MENU_LIST.flatMap(g => g.items);
  }, []);
  
  // 使用useMemo优化当前页面组件
  const CurrentPage = useMemo(() => {
    return CURRENT_PAGE_COMPONENTS[tab] || (() => <div style={{padding:32,textAlign:'center',color:'#888'}}>页面开发中</div>);
  }, [tab]);
  
  useEffect(() => {
    // 保存原始fetch函数
    if (!window._originalFetch) {
      window._originalFetch = window.fetch;
    }
    
    // 初始化认证头
    if (authUtils.isLoggedIn()) {
      authUtils.setAuthHeader(() => {
        localStorage.removeItem('token');
        setLogin(false);
      });
    }
  }, []); 
  
  if (!login) {
    return <LoginPage onLogin={handleLogin} />;
  }

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
                onClick={() => handleTabChange(item.key)}
                onMouseEnter={(e) => {
                  showTip(e, item.label);
                  preloadComponent(item.key);
                }}
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
                    onClick={() => handleTabChange(item.key)}
                    onMouseEnter={(e) => {
                      showTip(e, item.label);
                      preloadComponent(item.key);
                    }}
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
        <ErrorBoundary
          fallback={({ error, onRetry }) => (
            <div className="error-fallback">
              <h3>页面加载失败</h3>
              <p>{error?.message || '未知错误'}</p>
              <button onClick={onRetry} className="retry-btn">
                重试
              </button>
            </div>
          )}
        >
          <Suspense fallback={
            <LoadingSpinner 
              type="spinner" 
              text="Loading..." 
              size="lg"
              className="app-loading"
            />
          }>
            <CurrentPage collapsed={collapsed} onToggleCollapsed={toggleCollapsed} />
          </Suspense>
        </ErrorBoundary>
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