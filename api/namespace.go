package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterNamespace 注册 Namespace 相关路由，包括列表和详情接口
func RegisterNamespace(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listNamespaces func(context.Context, *kubernetes.Clientset) ([]model.NamespaceDetail, error),
) {
	r.GET("/namespaces", getNamespaceList(logger, getK8sClient, listNamespaces))
	r.GET("/namespaces/:name", getNamespaceDetail(logger, getK8sClient))
}

// getNamespaceList 获取Namespace列表的处理函数
// @Summary 获取 Namespace 列表
// @Description 获取所有命名空间列表，支持分页
// @Tags Namespace
// @Security BearerAuth
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /namespaces [get]
func getNamespaceList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listNamespaces func(context.Context, *kubernetes.Clientset) ([]model.NamespaceDetail, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		namespaces, err := listNamespaces(ctx, clientset)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		total := len(namespaces)
		var paged []model.NamespaceDetail
		if len(namespaces) > 0 {
			paged = Paginate(namespaces, offset, limit)
		} else {
			paged = []model.NamespaceDetail{}
		}
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		})
	}
}

// getNamespaceDetail 获取Namespace详情的处理函数
// @Summary 获取 Namespace 详情
// @Description 获取指定命名空间的详细信息
// @Tags Namespace
// @Security BearerAuth
// @Param name path string true "命名空间名称"
// @Success 200 {object} model.APIResponse
// @Router /namespaces/{name} [get]
func getNamespaceDetail(
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
		name := c.Param("name")
		ns, err := clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		namespaceDetail := model.NamespaceDetail{
			Name:        ns.Name,
			Status:      string(ns.Status.Phase),
			Labels:      ns.Labels,
			Annotations: ns.Annotations,
		}
		middleware.ResponseSuccess(c, namespaceDetail, "success", nil)
	}
}
