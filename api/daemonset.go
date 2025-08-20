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

// RegisterDaemonSet 注册 DaemonSet 相关路由
func RegisterDaemonSet(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listDaemonSets func(context.Context, *kubernetes.Clientset, string) ([]model.DaemonSetStatus, error),
) {
	r.GET("/daemonsets", getDaemonSetList(logger, getK8sClient, listDaemonSets))
	r.GET("/daemonsets/:namespace/:name", getDaemonSetDetail(logger, getK8sClient))
}

// getDaemonSetList 获取DaemonSet列表的处理函数
// @Summary 获取 DaemonSet 列表
// @Description 获取DaemonSet列表，支持分页
// @Tags DaemonSet
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /daemonsets [get]
func getDaemonSetList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listDaemonSets func(context.Context, *kubernetes.Clientset, string) ([]model.DaemonSetStatus, error),
) gin.HandlerFunc {
	return GenericListHandler(logger, getK8sClient, listDaemonSets)
}

// getDaemonSetDetail 获取DaemonSet详情的处理函数
// @Summary 获取 DaemonSet 详情
// @Description 获取指定命名空间下的DaemonSet详情
// @Tags DaemonSet
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "DaemonSet 名称"
// @Success 200 {object} model.APIResponse
// @Router /daemonsets/{namespace}/{name} [get]
func getDaemonSetDetail(
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
			Namespace:   ds.Namespace,
			Name:        ds.Name,
			Available:   ds.Status.NumberReady,
			Desired:     ds.Status.DesiredNumberScheduled,
			Status:      status,
			Labels:      ds.Labels,
			Annotations: ds.Annotations,
			Selector:    ds.Spec.Selector.MatchLabels,
			Image:       image,
		}
		middleware.ResponseSuccess(c, daemonSetDetail, "success", nil)
	}
}
