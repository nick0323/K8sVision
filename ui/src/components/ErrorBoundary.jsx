import React from 'react';

/**
 * 全局错误边界组件
 * 捕获和处理React组件树中的JavaScript错误
 */
class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
      errorId: null,
    };
  }

  static getDerivedStateFromError(error) {
    // 更新state，使下一次渲染能够显示降级后的UI
    return {
      hasError: true,
      errorId: Date.now().toString(36) + Math.random().toString(36).substr(2),
    };
  }

  componentDidCatch(error, errorInfo) {
    // 记录错误信息
    this.setState({
      error,
      errorInfo,
    });

    // 发送错误报告到监控服务
    this.reportError(error, errorInfo);
  }

  reportError = (error, errorInfo) => {
    // 这里可以集成错误监控服务，如Sentry
    console.error('ErrorBoundary caught an error:', error, errorInfo);
    
    // 发送到后端错误收集端点
    if (process.env.NODE_ENV === 'production') {
      fetch('/api/errors', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          errorId: this.state.errorId,
          message: error.message,
          stack: error.stack,
          componentStack: errorInfo.componentStack,
          timestamp: new Date().toISOString(),
          userAgent: navigator.userAgent,
          url: window.location.href,
        }),
      }).catch(console.error);
    }
  };

  handleRetry = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
      errorId: null,
    });
  };

  handleReload = () => {
    window.location.reload();
  };

  render() {
    if (this.state.hasError) {
      const { error, errorInfo, errorId } = this.state;
      const { fallback: Fallback, showDetails = false } = this.props;

      // 如果有自定义的fallback组件，使用它
      if (Fallback) {
        return <Fallback error={error} errorInfo={errorInfo} onRetry={this.handleRetry} />;
      }

      // 默认错误UI
      return (
        <div className="error-boundary" style={{
          padding: '20px',
          maxWidth: '600px',
          margin: '0 auto',
          fontFamily: 'Arial, sans-serif'
        }}>
          <div style={{
            border: '1px solid #dc3545',
            borderRadius: '8px',
            overflow: 'hidden',
            boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
          }}>
            <div style={{
              backgroundColor: '#dc3545',
              color: 'white',
              padding: '15px 20px',
              fontSize: '18px',
              fontWeight: 'bold'
            }}>
              ⚠️ 应用程序出现错误
            </div>
            <div style={{ padding: '20px' }}>
              <div style={{
                backgroundColor: '#f8d7da',
                color: '#721c24',
                padding: '12px',
                borderRadius: '4px',
                marginBottom: '15px',
                border: '1px solid #f5c6cb'
              }}>
                <strong>错误ID:</strong> {errorId}<br />
                <strong>时间:</strong> {new Date().toLocaleString()}
              </div>

              <p style={{ marginBottom: '15px', lineHeight: '1.5' }}>
                抱歉，应用程序遇到了一个意外错误。请尝试以下操作：
              </p>

              <div style={{ 
                display: 'flex', 
                gap: '10px', 
                marginBottom: '15px',
                flexWrap: 'wrap'
              }}>
                <button 
                  onClick={this.handleRetry}
                  style={{
                    backgroundColor: '#007bff',
                    color: 'white',
                    border: 'none',
                    padding: '8px 16px',
                    borderRadius: '4px',
                    cursor: 'pointer',
                    fontSize: '14px'
                  }}
                >
                  🔄 重试
                </button>
                <button 
                  onClick={this.handleReload}
                  style={{
                    backgroundColor: '#6c757d',
                    color: 'white',
                    border: 'none',
                    padding: '8px 16px',
                    borderRadius: '4px',
                    cursor: 'pointer',
                    fontSize: '14px'
                  }}
                >
                  🔄 重新加载页面
                </button>
              </div>

              {showDetails && error && (
                <details style={{ marginTop: '15px' }}>
                  <summary style={{ 
                    color: '#6c757d', 
                    cursor: 'pointer',
                    marginBottom: '10px'
                  }}>
                    错误详情
                  </summary>
                  <pre style={{
                    marginTop: '10px',
                    padding: '10px',
                    backgroundColor: '#f8f9fa',
                    borderRadius: '4px',
                    fontSize: '12px',
                    overflow: 'auto',
                    border: '1px solid #dee2e6'
                  }}>
                    {error.toString()}
                    {errorInfo && errorInfo.componentStack}
                  </pre>
                </details>
              )}

              <div style={{ 
                marginTop: '15px', 
                color: '#6c757d',
                fontSize: '14px'
              }}>
                如果问题持续存在，请联系技术支持并提供错误ID。
              </div>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

/**
 * 功能组件错误边界
 * 使用React 18的错误边界Hook
 */
export const ErrorBoundaryHook = ({ children, fallback, onError }) => {
  const [hasError, setHasError] = React.useState(false);
  const [error, setError] = React.useState(null);

  React.useEffect(() => {
    const handleError = (event) => {
      setHasError(true);
      setError(event.error);
      if (onError) {
        onError(event.error, event.errorInfo);
      }
    };

    window.addEventListener('error', handleError);
    window.addEventListener('unhandledrejection', handleError);

    return () => {
      window.removeEventListener('error', handleError);
      window.removeEventListener('unhandledrejection', handleError);
    };
  }, [onError]);

  if (hasError) {
    return fallback ? fallback(error) : (
      <div style={{
        padding: '20px',
        textAlign: 'center',
        color: '#dc3545',
        backgroundColor: '#f8d7da',
        border: '1px solid #f5c6cb',
        borderRadius: '4px',
        margin: '20px'
      }}>
        <h3>出现错误</h3>
        <p>{error?.message || '未知错误'}</p>
      </div>
    );
  }

  return children;
};

/**
 * 异步错误边界
 * 专门处理异步操作中的错误
 */
export const AsyncErrorBoundary = ({ children, onAsyncError }) => {
  const [asyncError, setAsyncError] = React.useState(null);

  const handleAsyncError = React.useCallback((error) => {
    setAsyncError(error);
    if (onAsyncError) {
      onAsyncError(error);
    }
  }, [onAsyncError]);

  React.useEffect(() => {
    // 监听未处理的Promise拒绝
    const handleUnhandledRejection = (event) => {
      handleAsyncError(event.reason);
    };

    window.addEventListener('unhandledrejection', handleUnhandledRejection);
    return () => {
      window.removeEventListener('unhandledrejection', handleUnhandledRejection);
    };
  }, [handleAsyncError]);

  if (asyncError) {
    return (
      <div style={{
        backgroundColor: '#f8d7da',
        color: '#721c24',
        padding: '15px',
        borderRadius: '4px',
        border: '1px solid #f5c6cb',
        margin: '10px'
      }}>
        <h5 style={{ margin: '0 0 10px 0' }}>异步操作错误</h5>
        <p style={{ margin: '0 0 10px 0' }}>{asyncError.message}</p>
        <button 
          onClick={() => setAsyncError(null)}
          style={{
            backgroundColor: '#007bff',
            color: 'white',
            border: 'none',
            padding: '6px 12px',
            borderRadius: '4px',
            cursor: 'pointer'
          }}
        >
          重试
        </button>
      </div>
    );
  }

  return children;
};

export default ErrorBoundary;
