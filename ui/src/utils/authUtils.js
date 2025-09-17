/**
 * 认证工具函数模块
 */

/**
 * 获取认证头信息
 * @returns {Object} 包含Authorization和Content-Type的headers对象
 */
export function getAuthHeaders() {
  const token = localStorage.getItem('token');
  if (!token) {
    return {};
  }
  
  // 验证token格式
  const segments = token.split('.');
  if (segments.length !== 3) {
    localStorage.removeItem('token');
    return {};
  }
  
  return {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  };
}

/**
 * 验证token是否有效
 * @returns {boolean} token是否有效
 */
export function validateToken() {
  const token = localStorage.getItem('token');
  if (!token) return false;
  
  const segments = token.split('.');
  if (segments.length !== 3) {
    localStorage.removeItem('token');
    return false;
  }
  
  return true;
}

/**
 * 退出登录
 */
export function logout() {
  localStorage.removeItem('token');
  // 不清除 remembered_username 和 remembered_password，保持"记住我"功能
  // 触发页面刷新，让App组件重新检查登录状态
  window.location.reload();
}

/**
 * 设置认证拦截器
 * @param {Function} onLogout - 退出登录回调函数
 */
export function setAuthInterceptor(onLogout) {
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
