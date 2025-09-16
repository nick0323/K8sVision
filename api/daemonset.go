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

// RegisterDaemonSet 注册 DaemonSet 相关路由
func RegisterDaemonSet(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient K8sClientProvider, // 使用类型别名简化签名
	listDaemonSets func(context.Context, *kubernetes.Clientset, string) ([]model.DaemonSetStatus, error),
) {
	r.GET("/daemonsets", getDaemonSetList(logger, getK8sClient, listDaemonSets))
	r.GET("/daemonsets/:namespace/:name", getDaemonSetDetail(logger, getK8sClient))
}

// getDaemonSetList 获取DaemonSet列表的处理函数
func getDaemonSetList(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listDaemonSets func(context.Context, *kubernetes.Clientset, string) ([]model.DaemonSetStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		HandleListWithPagination(c, logger, func(ctx context.Context, params PaginationParams) ([]model.DaemonSetStatus, error) {
			clientset, _, err := getK8sClient()
			if err != nil {
				return nil, err
			}
			return listDaemonSets(ctx, clientset, params.Namespace)
		}, ListSuccessMessage) // 使用常量替代硬编码
	}
}

// getDaemonSetDetail 获取DaemonSet详情的处理函数
func getDaemonSetDetail(
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
		ds, err := clientset.AppsV1().DaemonSets(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		status := "Unknown"
		if ds.Status.NumberReady == ds.Status.DesiredNumberScheduled && ds.Status.DesiredNumberScheduled > 0 {
			status = "Ready"
		} else if ds.Status.NumberReady > 0 {
			status = "PartialAvailable"
		} else {
			status = "Not Ready"
		}

		image := ""
		if len(ds.Spec.Template.Spec.Containers) > 0 {
			image = ds.Spec.Template.Spec.Containers[0].Image
		}

		daemonSetDetail := model.DaemonSetDetail{
			WorkloadCommonFields: model.WorkloadCommonFields{
				CommonResourceFields: model.CommonResourceFields{
					Namespace: ds.Namespace,
					Name:      ds.Name,
					Status:    status,
					BaseMetadata: model.BaseMetadata{
						Labels:      ds.Labels,
						Annotations: ds.Annotations,
					},
				},
				Available: ds.Status.NumberReady,
				Desired:   ds.Status.DesiredNumberScheduled,
				Selector:  ds.Spec.Selector.MatchLabels,
				Image:     image,
			},
		}
		middleware.ResponseSuccess(c, daemonSetDetail, DetailSuccessMessage, nil)
	}
}
