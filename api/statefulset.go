package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/nick0323/K8sVision/model"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterStatefulSet 注册 StatefulSet 相关路由，包括列表和详情接口
// @Summary 获取 StatefulSet 列表
// @Description 支持分页
// @Tags StatefulSet
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /statefulsets [get]
//
// @Summary 获取 StatefulSet 详情
// @Description 获取指定命名空间下的 StatefulSet 详情
// @Tags StatefulSet
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "StatefulSet 名称"
// @Success 200 {object} model.APIResponse
// @Router /statefulsets/{namespace}/{name} [get]
func RegisterStatefulSet(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listStatefulSets func(context.Context, *kubernetes.Clientset, string) ([]model.StatefulSetStatus, error),
) {
	r.GET("/statefulsets", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		namespace := c.DefaultQuery("namespace", "")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		sets, err := listStatefulSets(ctx, clientset, namespace)
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		paged := Paginate(sets, offset, limit)
		ResponseOK(c, paged, "success", &model.PageMeta{
			Total:  len(sets),
			Limit:  limit,
			Offset: offset,
		})
	})

	r.GET("/statefulsets/:namespace/:name", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		ns := c.Param("namespace")
		name := c.Param("name")
		sts, err := clientset.AppsV1().StatefulSets(ns).Get(ctx, name, v1.GetOptions{})
		if err != nil {
			ResponseError(c, logger, err, http.StatusNotFound)
			return
		}
		ResponseOK(c, sts, "success", nil)
	})
}
