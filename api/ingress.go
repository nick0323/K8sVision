package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

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
	r.GET("/ingresses", getIngressList(logger, getK8sClient, listIngresses))
	r.GET("/ingresses/:namespace/:name", getIngressDetail(logger, getK8sClient))
}

// getIngressList 获取Ingress列表的处理函数
// @Summary 获取 Ingress 列表
// @Description 获取Ingress列表，支持分页和搜索
// @Tags Ingress
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Param search query string false "搜索关键词（支持名称、命名空间、主机等字段搜索）"
// @Success 200 {object} model.APIResponse
// @Router /ingresses [get]
func getIngressList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listIngresses func(context.Context, *kubernetes.Clientset, string) ([]model.IngressStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		namespace := c.DefaultQuery("namespace", "")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		search := c.DefaultQuery("search", "") // 新增：搜索关键词

		ingresses, err := listIngresses(ctx, clientset, namespace)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 新增：如果提供了搜索关键词，先进行搜索过滤
		var filteredIngresses []model.IngressStatus
		if search != "" {
			filteredIngresses = filterIngressesBySearch(ingresses, search)
		} else {
			filteredIngresses = ingresses
		}

		// 对过滤后的数据进行分页
		paged := Paginate(filteredIngresses, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(filteredIngresses), // 使用过滤后的总数
			Limit:  limit,
			Offset: offset,
		})
	}
}

// filterIngressesBySearch 根据搜索关键词过滤Ingress
func filterIngressesBySearch(ingresses []model.IngressStatus, search string) []model.IngressStatus {
	if search == "" {
		return ingresses
	}
	searchLower := strings.ToLower(search)
	var filtered []model.IngressStatus
	for _, ingress := range ingresses {
		if strings.Contains(strings.ToLower(ingress.Name), searchLower) ||
			strings.Contains(strings.ToLower(ingress.Namespace), searchLower) ||
			strings.Contains(strings.Join(ingress.Hosts, ","), searchLower) ||
			strings.Contains(strings.ToLower(ingress.Address), searchLower) ||
			strings.Contains(strings.ToLower(ingress.Class), searchLower) ||
			strings.Contains(strings.ToLower(ingress.Status), searchLower) ||
			strings.Contains(strings.Join(ingress.Ports, ","), searchLower) ||
			strings.Contains(strings.Join(ingress.Path, ","), searchLower) ||
			strings.Contains(strings.Join(ingress.TargetService, ","), searchLower) {
			filtered = append(filtered, ingress)
		}
	}
	return filtered
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
