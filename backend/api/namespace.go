package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/nick0323/K8sVision/backend/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterNamespace 注册 Namespace 相关路由，包括列表和详情接口
// @Summary 获取 Namespace 列表
// @Description 获取所有命名空间
// @Tags Namespace
// @Success 200 {object} model.APIResponse
// @Router /namespaces [get]
//
// @Summary 获取 Namespace 详情
// @Description 获取指定命名空间详情
// @Tags Namespace
// @Param name path string true "命名空间"
// @Success 200 {object} model.APIResponse
// @Router /namespaces/{name} [get]
func RegisterNamespace(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listNamespaces func(context.Context, *kubernetes.Clientset) ([]string, error),
) {
	r.GET("/namespaces", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		namespaces, err := listNamespaces(ctx, clientset)
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		total := len(namespaces)
		paged := Paginate(namespaces, offset, limit)
		ResponseOK(c, paged, "success", &model.PageMeta{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		})
	})

	r.GET("/namespaces/:name", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		name := c.Param("name")
		// 这里只做简单查找
		ns, err := clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			ResponseError(c, logger, err, http.StatusNotFound)
			return
		}
		ResponseOK(c, ns, "success", nil)
	})
}
