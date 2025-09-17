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

func RegisterNamespace(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listNamespaces func(context.Context, *kubernetes.Clientset) ([]model.NamespaceDetail, error),
) {
	r.GET("/namespaces", getNamespaceList(logger, getK8sClient, listNamespaces))
	r.GET("/namespaces/:name", getNamespaceDetail(logger, getK8sClient))
}

func getNamespaceList(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listNamespaces func(context.Context, *kubernetes.Clientset) ([]model.NamespaceDetail, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 命名空间用于下拉选择，返回全集更符合预期（不做分页）
		clientset, _, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := GetRequestContext(c)
		ns, err := listNamespaces(ctx, clientset)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		middleware.ResponseSuccess(c, ns, ListSuccessMessage, nil)
	}
}

func getNamespaceDetail(
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
		name := c.Param("name")
		ns, err := clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		namespaceDetail := model.NamespaceDetail{
			Name:   ns.Name,
			Status: string(ns.Status.Phase),
			BaseMetadata: model.BaseMetadata{
				Labels:      ns.Labels,
				Annotations: ns.Annotations,
			},
		}
		middleware.ResponseSuccess(c, namespaceDetail, DetailSuccessMessage, nil)
	}
}
