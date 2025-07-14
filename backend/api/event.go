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
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterEvent 注册 Event 相关路由，包括列表和详情接口
// @Summary 获取 Event 列表
// @Description 支持分页
// @Tags Event
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /events [get]
//
// @Summary 获取 Event 详情
// @Description 获取指定命名空间下的 Event 详情
// @Tags Event
// @Param namespace path string true "命名空间"
// @Param name path string true "Event 名称"
// @Success 200 {object} model.APIResponse
// @Router /events/{namespace}/{name} [get]
func RegisterEvent(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *metrics.Clientset, error),
	listEvents func(context.Context, *kubernetes.Clientset) ([]model.EventStatus, error),
) {
	r.GET("/events", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		events, err := listEvents(ctx, clientset)
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		paged := Paginate(events, offset, limit)
		ResponseOK(c, paged, "success", &model.PageMeta{
			Total:  len(events),
			Limit:  limit,
			Offset: offset,
		})
	})

	r.GET("/events/:namespace/:name", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		ns := c.Param("namespace")
		name := c.Param("name")
		evt, err := clientset.CoreV1().Events(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			ResponseError(c, logger, err, http.StatusNotFound)
			return
		}
		ResponseOK(c, evt, "success", nil)
	})
}
