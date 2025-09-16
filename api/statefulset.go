package api

import (
	"context"
	"net/http"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"
	"github.com/nick0323/K8sVision/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RegisterStatefulSet 注册 StatefulSet 相关路由
func RegisterStatefulSet(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient K8sClientProvider, // 使用类型别名简化签名
	listStatefulSets func(context.Context, *kubernetes.Clientset, string) ([]model.StatefulSetStatus, error),
) {
	r.GET("/statefulsets", getStatefulSetList(logger, getK8sClient, listStatefulSets))
	r.GET("/statefulsets/:namespace/:name", getStatefulSetDetail(logger, getK8sClient))
}

// getStatefulSetList 获取StatefulSet列表的处理函数
func getStatefulSetList(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listStatefulSets func(context.Context, *kubernetes.Clientset, string) ([]model.StatefulSetStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		HandleListWithPagination(c, logger, func(ctx context.Context, params PaginationParams) ([]model.StatefulSetStatus, error) {
			clientset, _, err := getK8sClient()
			if err != nil {
				return nil, err
			}
			return listStatefulSets(ctx, clientset, params.Namespace)
		}, ListSuccessMessage) // 使用常量替代硬编码
	}
}

// getStatefulSetDetail 获取StatefulSet详情的处理函数
func getStatefulSetDetail(
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
		sts, err := clientset.AppsV1().StatefulSets(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		image := ""
		if len(sts.Spec.Template.Spec.Containers) > 0 {
			image = sts.Spec.Template.Spec.Containers[0].Image
		}

		statefulSetDetail := model.StatefulSetDetail{
			WorkloadCommonFields: model.WorkloadCommonFields{
				CommonResourceFields: model.CommonResourceFields{
					Namespace: sts.Namespace,
					Name:      sts.Name,
					Status:    service.GetWorkloadStatus(sts.Status.AvailableReplicas, *sts.Spec.Replicas),
					BaseMetadata: model.BaseMetadata{
						Labels:      sts.Labels,
						Annotations: sts.Annotations,
					},
				},
				Available: sts.Status.AvailableReplicas,
				Desired:   *sts.Spec.Replicas,
				Selector:  sts.Spec.Selector.MatchLabels,
				Image:     image,
			},
			Replicas:    *sts.Spec.Replicas,
			ServiceName: sts.Spec.ServiceName,
		}
		middleware.ResponseSuccess(c, statefulSetDetail, DetailSuccessMessage, nil)
	}
}
