package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterService 注册 Service 相关路由，包括列表和详情接口
// @Summary 获取 Service 列表
// @Description 支持分页
// @Tags Service
// @Security BearerAuth
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /services [get]
//
// @Summary 获取 Service 详情
// @Description 获取指定命名空间下的 Service 详情
// @Tags Service
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "Service 名称"
// @Success 200 {object} model.APIResponse
// @Router /services/{namespace}/{name} [get]
func RegisterService(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listServices func(context.Context, *kubernetes.Clientset) ([]model.ServiceStatus, error),
) {
	r.GET("/services", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		services, err := listServices(ctx, clientset)
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		paged := Paginate(services, offset, limit)
		ResponseOK(c, paged, "success", &model.PageMeta{
			Total:  len(services),
			Limit:  limit,
			Offset: offset,
		})
	})

	r.GET("/services/:namespace/:name", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		ns := c.Param("namespace")
		name := c.Param("name")
		svc, err := clientset.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			ResponseError(c, logger, err, http.StatusNotFound)
			return
		}
		ports := make([]string, 0, len(svc.Spec.Ports))
		for _, p := range svc.Spec.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", p.Port, p.Protocol))
		}
		svcDetail := model.ServiceStatus{
			Namespace: svc.Namespace,
			Name:      svc.Name,
			Type:      string(svc.Spec.Type),
			ClusterIP: svc.Spec.ClusterIP,
			Ports:     ports,
		}
		ResponseOK(c, svcDetail, "success", nil)
	})
}
