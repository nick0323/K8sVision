package service

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/nick0323/K8sVision/config"
	"github.com/nick0323/K8sVision/cache"
	"github.com/nick0323/K8sVision/model"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

var (
	configManager *config.Manager
	cacheManager  *cache.Manager
	clientsCache  *cache.MemoryCache
	clientsMutex  sync.RWMutex
)

// SetConfigManager 设置配置管理器
func SetConfigManager(cm *config.Manager) {
	configManager = cm
}

// SetCacheManager 设置缓存管理器
func SetCacheManager(cm *cache.Manager) {
	cacheManager = cm
	// 创建客户端缓存
	clientsCache = cache.NewMemoryCache(&model.CacheConfig{
		Enabled:         true,
		Type:           "memory",
		TTL:            30 * time.Minute, // 客户端缓存30分钟
		MaxSize:        10,               // 最多缓存10个客户端
		CleanupInterval: 5 * time.Minute,
	}, cm.GetLogger()) // 从缓存管理器获取logger
}

// GetK8sConfig 获取 Kubernetes 连接配置，优先环境变量，其次配置文件，最后集群内
func GetK8sConfig() (*rest.Config, error) {
	var k8sConfig *model.KubernetesConfig
	
	// 优先使用配置管理器
	if configManager != nil {
		cfg := configManager.GetConfig()
		k8sConfig = &cfg.Kubernetes
	} else {
		// 兼容旧版本，使用环境变量和viper
		k8sConfig = &model.KubernetesConfig{
			Kubeconfig: os.Getenv("KUBECONFIG"),
			Timeout:    30 * time.Second,
			QPS:        100,
			Burst:      200,
			Insecure:   true,
		}
	}

	// 优先环境变量
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = k8sConfig.Kubeconfig
	}

	// 如果指定了kubeconfig文件
	if kubeconfig != "" {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		
		// 应用配置
		applyK8sConfig(config, k8sConfig)
		return config, nil
	}

	// 如果指定了API服务器和Token
	if k8sConfig.APIServer != "" && k8sConfig.Token != "" {
		config := &rest.Config{
			Host:        k8sConfig.APIServer,
			BearerToken: k8sConfig.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: k8sConfig.Insecure,
			},
		}
		
		// 应用配置
		applyK8sConfig(config, k8sConfig)
		return config, nil
	}

	// fallback: 集群内配置
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	
	// 应用配置
	applyK8sConfig(config, k8sConfig)
	return config, nil
}

// applyK8sConfig 应用Kubernetes配置
func applyK8sConfig(config *rest.Config, k8sConfig *model.KubernetesConfig) {
	// 设置超时
	if k8sConfig.Timeout > 0 {
		config.Timeout = k8sConfig.Timeout
	}
	
	// 设置QPS和Burst
	if k8sConfig.QPS > 0 {
		config.QPS = k8sConfig.QPS
	}
	if k8sConfig.Burst > 0 {
		config.Burst = k8sConfig.Burst
	}
	
	// 设置TLS配置
	if k8sConfig.Insecure {
		config.Insecure = true
		config.TLSClientConfig.CAFile = ""
		config.TLSClientConfig.CAData = nil
	}
	
	// 设置证书文件
	if k8sConfig.CertFile != "" {
		config.TLSClientConfig.CertFile = k8sConfig.CertFile
	}
	if k8sConfig.KeyFile != "" {
		config.TLSClientConfig.KeyFile = k8sConfig.KeyFile
	}
	if k8sConfig.CAFile != "" {
		config.TLSClientConfig.CAFile = k8sConfig.CAFile
	}
}

// GetK8sClient 获取 Kubernetes clientset 和 metrics client（带缓存）
func GetK8sClient() (*kubernetes.Clientset, *metrics.Clientset, error) {
	// 生成缓存键
	cacheKey := "k8s_clients"
	
	// 尝试从缓存获取
	if clientsCache != nil {
		if cached, exists := clientsCache.Get(cacheKey); exists {
			if clients, ok := cached.([]interface{}); ok && len(clients) == 2 {
				if k8sClient, ok := clients[0].(*kubernetes.Clientset); ok {
					if metricsClient, ok := clients[1].(*metrics.Clientset); ok {
						return k8sClient, metricsClient, nil
					}
				}
			}
		}
	}
	
	// 缓存未命中，创建新的客户端
	config, err := GetK8sConfig()
	if err != nil {
		return nil, nil, err
	}
	
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	
	metricsClient, err := metrics.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	
	// 缓存客户端
	if clientsCache != nil {
		clientsCache.Set(cacheKey, []interface{}{clientset, metricsClient})
	}
	
	return clientset, metricsClient, nil
}

// GetK8sClientWithContext 获取带上下文的Kubernetes客户端
func GetK8sClientWithContext(ctx context.Context) (*kubernetes.Clientset, *metrics.Clientset, error) {
	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
	}
	
	return GetK8sClient()
}

// ClearK8sClientCache 清除Kubernetes客户端缓存
func ClearK8sClientCache() {
	if clientsCache != nil {
		clientsCache.Clear()
	}
}

// GetCachedResource 获取缓存的资源数据
func GetCachedResource[T any](key string, ttl time.Duration, loader func() (T, error)) (T, error) {
	var zero T
	
	if cacheManager == nil {
		// 缓存未启用，直接加载
		return loader()
	}
	
	// 尝试从缓存获取
	if cached, exists := cacheManager.Get(key); exists {
		if value, ok := cached.(T); ok {
			return value, nil
		}
	}
	
	// 缓存未命中，加载数据
	value, err := loader()
	if err != nil {
		return zero, err
	}
	
	// 设置缓存
	cacheManager.SetWithTTL(key, value, ttl)
	return value, nil
}

// GetCachedResourceWithCache 从指定缓存获取资源数据
func GetCachedResourceWithCache[T any](cacheName, key string, ttl time.Duration, loader func() (T, error)) (T, error) {
	var zero T
	
	if cacheManager == nil {
		// 缓存未启用，直接加载
		return loader()
	}
	
	// 尝试从指定缓存获取
	if value, err := cacheManager.GetOrSetWithCache(cacheName, key, ttl, func() (interface{}, error) {
		return loader()
	}); err != nil {
		return zero, err
	} else if result, ok := value.(T); ok {
		return result, nil
	}
	
	return zero, nil
}
