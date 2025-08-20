package api

import (
	"context"
	"net/http"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterConfigMap 注册 ConfigMap 相关路由
func RegisterConfigMap(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listConfigMaps func(context.Context, *kubernetes.Clientset, string) ([]model.ConfigMapStatus, error),
) {
	r.GET("/configmaps", getConfigMapList(logger, getK8sClient, listConfigMaps))
	r.GET("/configmaps/:namespace/:name", getConfigMapDetail(logger, getK8sClient))
}

// getConfigMapList 获取ConfigMap列表的处理函数
// @Summary 获取 ConfigMap 列表
// @Description 获取ConfigMap列表，支持分页
// @Tags ConfigMap
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /configmaps [get]
func getConfigMapList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listConfigMaps func(context.Context, *kubernetes.Clientset, string) ([]model.ConfigMapStatus, error),
) gin.HandlerFunc {
	return GenericListHandler(logger, getK8sClient, listConfigMaps)
}

// getConfigMapDetail 获取ConfigMap详情的处理函数
// @Summary 获取 ConfigMap 详情
// @Description 获取指定命名空间下的ConfigMap详情
// @Tags ConfigMap
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "ConfigMap 名称"
// @Success 200 {object} model.APIResponse
// @Router /configmaps/{namespace}/{name} [get]
func getConfigMapDetail(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		ns := c.Param("namespace")
		name := c.Param("name")
		configMap, err := clientset.CoreV1().ConfigMaps(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		keys := make([]string, 0, len(configMap.Data))
		for key := range configMap.Data {
			keys = append(keys, key)
		}

		configMapDetail := model.ConfigMapDetail{
			Namespace:   configMap.Namespace,
			Name:        configMap.Name,
			DataCount:   len(configMap.Data),
			Keys:        keys,
			Labels:      configMap.Labels,
			Annotations: configMap.Annotations,
			Data:        configMap.Data,
		}
		middleware.ResponseSuccess(c, configMapDetail, "success", nil)
	}
}
