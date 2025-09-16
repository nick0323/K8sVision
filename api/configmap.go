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
)

func RegisterConfigMap(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listConfigMaps func(context.Context, *kubernetes.Clientset, string) ([]model.ConfigMapStatus, error),
) {
	r.GET("/configmaps", getConfigMapList(logger, getK8sClient, listConfigMaps))
	r.GET("/configmaps/:namespace/:name", getConfigMapDetail(logger, getK8sClient))
}

func getConfigMapList(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listConfigMaps func(context.Context, *kubernetes.Clientset, string) ([]model.ConfigMapStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		HandleListWithPagination(c, logger, func(ctx context.Context, params PaginationParams) ([]model.ConfigMapStatus, error) {
			clientset, _, err := getK8sClient()
			if err != nil {
				return nil, err
			}
			return listConfigMaps(ctx, clientset, params.Namespace)
		}, ListSuccessMessage)
	}
}

func getConfigMapDetail(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := GetRequestContext(c)
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
			CommonResourceFields: model.CommonResourceFields{
				Namespace: configMap.Namespace,
				Name:      configMap.Name,
				Status:    "Active",
				BaseMetadata: model.BaseMetadata{
					Labels:      configMap.Labels,
					Annotations: configMap.Annotations,
				},
			},
			DataCount: len(configMap.Data),
			Keys:      keys,
			Data:      configMap.Data,
		}
		middleware.ResponseSuccess(c, configMapDetail, DetailSuccessMessage, nil)
	}
}
