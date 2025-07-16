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

// RegisterDaemonSet 注册 DaemonSet 相关路由，包括列表和详情接口
// @Summary 获取 DaemonSet 列表
// @Description 支持分页
// @Tags DaemonSet
// @Security BearerAuth
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /daemonsets [get]
//
// @Summary 获取 DaemonSet 详情
// @Description 获取指定命名空间下的 DaemonSet 详情
// @Tags DaemonSet
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "DaemonSet 名称"
// @Success 200 {object} model.APIResponse
// @Router /daemonsets/{namespace}/{name} [get]
func RegisterDaemonSet(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listDaemonSets func(context.Context, *kubernetes.Clientset, string) ([]model.DaemonSetStatus, error),
) {
	r.GET("/daemonsets", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		namespace := c.DefaultQuery("namespace", "")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		daemonsets, err := listDaemonSets(ctx, clientset, namespace)
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		paged := Paginate(daemonsets, offset, limit)
		ResponseOK(c, paged, "success", &model.PageMeta{
			Total:  len(daemonsets),
			Limit:  limit,
			Offset: offset,
		})
	})

	r.GET("/daemonsets/:namespace/:name", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		ns := c.Param("namespace")
		name := c.Param("name")
		ds, err := clientset.AppsV1().DaemonSets(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			ResponseError(c, logger, err, http.StatusNotFound)
			return
		}
		ResponseOK(c, ds, "success", nil)
	})
}
