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
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterDeployment 注册 Deployment 相关路由，包括列表和详情接口
func RegisterDeployment(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listDeployments func(context.Context, *kubernetes.Clientset, string) ([]model.DeploymentStatus, error),
) {
	r.GET("/deployments", getDeploymentList(logger, getK8sClient, listDeployments))
	r.GET("/deployments/:namespace/:name", getDeploymentDetail(logger, getK8sClient))
}

// getDeploymentList 获取Deployment列表的处理函数
// @Summary 获取 Deployment 列表
// @Description 获取Deployment列表，支持分页
// @Tags Deployment
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /deployments [get]
func getDeploymentList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listDeployments func(context.Context, *kubernetes.Clientset, string) ([]model.DeploymentStatus, error),
) gin.HandlerFunc {
	return GenericListHandler(logger, getK8sClient, listDeployments)
}

// getDeploymentDetail 获取Deployment详情的处理函数
// @Summary 获取 Deployment 详情
// @Description 获取指定命名空间下的Deployment详情
// @Tags Deployment
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "Deployment 名称"
// @Success 200 {object} model.APIResponse
// @Router /deployments/{namespace}/{name} [get]
func getDeploymentDetail(
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
		dep, err := clientset.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		image := ""
		if len(dep.Spec.Template.Spec.Containers) > 0 {
			image = dep.Spec.Template.Spec.Containers[0].Image
		}

		deploymentDetail := model.DeploymentDetail{
			Namespace:   dep.Namespace,
			Name:        dep.Name,
			Replicas:    *dep.Spec.Replicas,
			Available:   dep.Status.AvailableReplicas,
			Desired:     *dep.Spec.Replicas,
			Status:      service.GetWorkloadStatus(dep.Status.AvailableReplicas, *dep.Spec.Replicas),
			Labels:      dep.Labels,
			Annotations: dep.Annotations,
			Selector:    dep.Spec.Selector.MatchLabels,
			Strategy:    string(dep.Spec.Strategy.Type),
			Image:       image,
		}
		middleware.ResponseSuccess(c, deploymentDetail, "success", nil)
	}
}
