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

// RegisterEvent 注册 Event 相关路由
func RegisterEvent(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listEvents func(context.Context, *kubernetes.Clientset, string) ([]model.EventStatus, error),
) {
	r.GET("/events", getEventList(logger, getK8sClient, listEvents))
	r.GET("/events/:namespace/:name", getEventDetail(logger, getK8sClient))
}

// getEventList 获取Event列表的处理函数
// @Summary 获取 Event 列表
// @Description 获取Event列表，支持分页
// @Tags Event
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /events [get]
func getEventList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listEvents func(context.Context, *kubernetes.Clientset, string) ([]model.EventStatus, error),
) gin.HandlerFunc {
	return GenericListHandler(logger, getK8sClient, listEvents)
}

// getEventDetail 获取Event详情的处理函数
// @Summary 获取 Event 详情
// @Description 获取指定命名空间下的Event详情
// @Tags Event
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "Event 名称"
// @Success 200 {object} model.APIResponse
// @Router /events/{namespace}/{name} [get]
func getEventDetail(
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
		event, err := clientset.CoreV1().Events(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		eventDetail := model.EventDetail{
			Namespace:   event.Namespace,
			Name:        event.Name,
			Reason:      event.Reason,
			Message:     event.Message,
			Type:        event.Type,
			Count:       event.Count,
			FirstSeen:   event.FirstTimestamp.Format("2006-01-02 15:04:05"),
			LastSeen:    event.LastTimestamp.Format("2006-01-02 15:04:05"),
			Duration:    event.LastTimestamp.Sub(event.FirstTimestamp.Time).String(),
			Labels:      event.Labels,
			Annotations: event.Annotations,
		}
		middleware.ResponseSuccess(c, eventDetail, "success", nil)
	}
}
