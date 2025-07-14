package service

import (
	"os"

	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

// GetK8sConfig 获取 Kubernetes 连接配置，优先环境变量，其次配置文件，最后集群内
func GetK8sConfig() (*rest.Config, error) {
	// 优先环境变量
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		// 其次读取 viper 配置
		kubeconfig = viper.GetString("kubeconfig")
	}
	if kubeconfig != "" {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err == nil {
			config.Insecure = true
			config.TLSClientConfig.CAFile = ""
			config.TLSClientConfig.CAData = nil
			return config, nil
		}
		return nil, err
	}
	// fallback: 集群内
	return rest.InClusterConfig()
}

// GetK8sClient 获取 Kubernetes clientset 和 metrics client
func GetK8sClient() (*kubernetes.Clientset, *metrics.Clientset, error) {
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
	return clientset, metricsClient, nil
}
