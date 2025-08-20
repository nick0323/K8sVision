package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterService 注册 Service 相关路由，包括列表和详情接口
func RegisterService(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listServices func(context.Context, *kubernetes.Clientset, string) ([]model.ServiceStatus, error),
) {
	r.GET("/services", getServiceList(logger, getK8sClient, listServices))
	r.GET("/services/:namespace/:name", getServiceDetail(logger, getK8sClient))
}

// getServiceList 获取Service列表的处理函数
// @Summary 获取 Service 列表
// @Description 获取Service列表，支持分页
// @Tags Service
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /services [get]
func getServiceList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listServices func(context.Context, *kubernetes.Clientset, string) ([]model.ServiceStatus, error),
) gin.HandlerFunc {
	return GenericListHandler(logger, getK8sClient, listServices)
}

// getServiceDetail 获取Service详情的处理函数
// @Summary 获取 Service 详情
// @Description 获取指定命名空间下的Service详情
// @Tags Service
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "Service 名称"
// @Success 200 {object} model.APIResponse
// @Router /services/{namespace}/{name} [get]
func getServiceDetail(
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
		svc, err := clientset.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}
		ports := make([]string, 0, len(svc.Spec.Ports))
		for _, p := range svc.Spec.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", p.Port, p.Protocol))
		}
		svcDetail := model.ServiceDetail{
			Namespace:   svc.Namespace,
			Name:        svc.Name,
			Type:        string(svc.Spec.Type),
			ClusterIP:   svc.Spec.ClusterIP,
			Ports:       ports,
			Labels:      svc.Labels,
			Annotations: svc.Annotations,
			Selector:    svc.Spec.Selector,
		}
		middleware.ResponseSuccess(c, svcDetail, "success", nil)
	}
}
