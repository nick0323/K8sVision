/**
 * 通用API工具函数
 * 消除重复的API调用逻辑
 */

// 基础API配置
const API_CONFIG = {
  baseURL: '',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  }
};

// 获取认证token
const getAuthToken = () => {
  return localStorage.getItem('token');
};

// 构建请求头
const buildHeaders = (customHeaders = {}) => {
  const token = getAuthToken();
  const headers = { ...API_CONFIG.headers, ...customHeaders };
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  return headers;
};

// 构建查询参数
const buildQueryParams = (params = {}) => {
  const queryParams = new URLSearchParams();
  
  Object.entries(params).forEach(([key, value]) => {
    if (value !== null && value !== undefined && value !== '') {
      queryParams.append(key, value.toString());
    }
  });
  
  return queryParams.toString();
};

// 通用GET请求
export const apiGet = async (endpoint, params = {}, options = {}) => {
  try {
    const queryString = buildQueryParams(params);
    const url = queryString ? `${endpoint}?${queryString}` : endpoint;
    
    const response = await fetch(url, {
      method: 'GET',
      headers: buildHeaders(options.headers),
      ...options
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    return { success: true, data, response };
  } catch (error) {
    console.error('API GET request failed:', error);
    return { success: false, error: error.message, data: null };
  }
};

// 通用POST请求
export const apiPost = async (endpoint, body = {}, options = {}) => {
  try {
    const response = await fetch(endpoint, {
      method: 'POST',
      headers: buildHeaders(options.headers),
      body: JSON.stringify(body),
      ...options
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    return { success: true, data, response };
  } catch (error) {
    console.error('API POST request failed:', error);
    return { success: false, error: error.message, data: null };
  }
};

// 通用PUT请求
export const apiPut = async (endpoint, body = {}, options = {}) => {
  try {
    const response = await fetch(endpoint, {
      method: 'PUT',
      headers: buildHeaders(options.headers),
      body: JSON.stringify(body),
      ...options
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    return { success: true, data, response };
  } catch (error) {
    console.error('API PUT request failed:', error);
    return { success: false, error: error.message, data: null };
  }
};

// 通用DELETE请求
export const apiDelete = async (endpoint, options = {}) => {
  try {
    const response = await fetch(endpoint, {
      method: 'DELETE',
      headers: buildHeaders(options.headers),
      ...options
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    return { success: true, data, response };
  } catch (error) {
    console.error('API DELETE request failed:', error);
    return { success: false, error: error.message, data: null };
  }
};

// 分页查询工具
export const createPaginatedQuery = (endpoint, options = {}) => {
  const {
    page = 1,
    pageSize = 20,
    namespace = '',
    search = '',
    sortBy = '',
    sortOrder = 'asc'
  } = options;
  
  const params = {
    limit: pageSize,
    offset: (page - 1) * pageSize,
    ...(namespace && { namespace }),
    ...(search && { search }),
    ...(sortBy && { sortBy }),
    ...(sortOrder && { sortOrder })
  };
  
  return apiGet(endpoint, params);
};

// 搜索工具
export const createSearchQuery = (endpoint, searchTerm, options = {}) => {
  const {
    namespace = '',
    limit = 1000
  } = options;
  
  const params = {
    search: searchTerm.trim(),
    limit,
    ...(namespace && { namespace })
  };
  
  return apiGet(endpoint, params);
};

// 资源详情查询
export const getResourceDetail = (endpoint, name, namespace = '') => {
  const params = { name };
  if (namespace) {
    params.namespace = namespace;
  }
  
  return apiGet(endpoint, params);
};

// 错误处理工具
export const handleApiError = (error, fallbackValue = null) => {
  console.error('API Error:', error);
  
  if (error.status === 401) {
    // 认证失败，清除token
    localStorage.removeItem('token');
    window.location.reload();
    return fallbackValue;
  }
  
  if (error.status === 403) {
    console.warn('Access forbidden');
    return fallbackValue;
  }
  
  if (error.status >= 500) {
    console.error('Server error');
    return fallbackValue;
  }
  
  return fallbackValue;
};

// 响应数据标准化
export const normalizeApiResponse = (response) => {
  if (!response || !response.success) {
    return { data: [], page: {}, error: response?.error };
  }
  
  const { data, page } = response.data;
  
  return {
    data: Array.isArray(data) ? data : [],
    page: page || {},
    error: null
  };
};

