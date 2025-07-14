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

// RegisterCronJob 注册 CronJob 相关路由，包括列表和详情接口
// @Summary 获取 CronJob 列表
// @Description 支持分页
// @Tags CronJob
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /cronjobs [get]
//
// @Summary 获取 CronJob 详情
// @Description 获取指定命名空间下的 CronJob 详情
// @Tags CronJob
// @Param namespace path string true "命名空间"
// @Param name path string true "CronJob 名称"
// @Success 200 {object} model.APIResponse
// @Router /cronjobs/{namespace}/{name} [get]
func RegisterCronJob(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listCronJobs func(context.Context, *kubernetes.Clientset) ([]model.CronJobStatus, error),
) {
	r.GET("/cronjobs", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		cronjobs, err := listCronJobs(ctx, clientset)
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		paged := Paginate(cronjobs, offset, limit)
		ResponseOK(c, paged, "success", &model.PageMeta{
			Total:  len(cronjobs),
			Limit:  limit,
			Offset: offset,
		})
	})

	r.GET("/cronjobs/:namespace/:name", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		ns := c.Param("namespace")
		name := c.Param("name")
		cronjob, err := clientset.BatchV1().CronJobs(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			ResponseError(c, logger, err, http.StatusNotFound)
			return
		}
		ResponseOK(c, cronjob, "success", nil)
	})
}
