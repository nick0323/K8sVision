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

// RegisterIngress 注册 Ingress 相关路由
func RegisterIngress(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listIngresses func(context.Context, *kubernetes.Clientset, string) ([]model.IngressStatus, error),
) {
	// 同时支持单数和复数形式
	r.GET("/ingress", getIngressList(logger, getK8sClient, listIngresses))
	r.GET("/ingress/:namespace/:name", getIngressDetail(logger, getK8sClient))
}

// getIngressList 获取Ingress列表的处理函数
// @Summary 获取 Ingress 列表
// @Description 获取Ingress列表，支持分页
// @Tags Ingress
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /ingresses [get]
func getIngressList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listIngresses func(context.Context, *kubernetes.Clientset, string) ([]model.IngressStatus, error),
) gin.HandlerFunc {
	return GenericListHandler(logger, getK8sClient, listIngresses)
}

// getIngressDetail 获取Ingress详情的处理函数
// @Summary 获取 Ingress 详情
// @Description 获取指定命名空间下的Ingress详情
// @Tags Ingress
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "Ingress 名称"
// @Success 200 {object} model.APIResponse
// @Router /ingresses/{namespace}/{name} [get]
func getIngressDetail(
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
		ingress, err := clientset.NetworkingV1().Ingresses(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		hosts := make([]string, 0, len(ingress.Spec.Rules))
		paths := make([]string, 0)
		targetServices := make([]string, 0)
		
		for _, rule := range ingress.Spec.Rules {
			hosts = append(hosts, rule.Host)
			
			// 解析路径和目标服务
			if rule.HTTP != nil {
				for _, path := range rule.HTTP.Paths {
					// 添加路径信息
					if path.Path != "" {
						paths = append(paths, path.Path)
					} else {
						paths = append(paths, "/")
					}
					
					// 添加目标服务信息
					if path.Backend.Service != nil {
						targetServices = append(targetServices, path.Backend.Service.Name)
					}
				}
			}
		}

		address := ""
		if len(ingress.Status.LoadBalancer.Ingress) > 0 {
			if ingress.Status.LoadBalancer.Ingress[0].IP != "" {
				address = ingress.Status.LoadBalancer.Ingress[0].IP
			} else if ingress.Status.LoadBalancer.Ingress[0].Hostname != "" {
				address = ingress.Status.LoadBalancer.Ingress[0].Hostname
			}
		}

		// 安全获取 IngressClass
		class := ""
		if ingress.Spec.IngressClassName != nil {
			class = *ingress.Spec.IngressClassName
		}

		ingressDetail := model.IngressDetail{
			Namespace:     ingress.Namespace,
			Name:          ingress.Name,
			Hosts:         hosts,
			Address:       address,
			Ports:         []string{"80", "443"},
			Class:         class,
			Status:        "Ready",
			Path:          paths,
			TargetService: targetServices,
			Labels:        ingress.Labels,
			Annotations:   ingress.Annotations,
		}
		middleware.ResponseSuccess(c, ingressDetail, "success", nil)
	}
}
