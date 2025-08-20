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

// RegisterCronJob 注册 CronJob 相关路由
func RegisterCronJob(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listCronJobs func(context.Context, *kubernetes.Clientset, string) ([]model.CronJobStatus, error),
) {
	r.GET("/cronjobs", getCronJobList(logger, getK8sClient, listCronJobs))
	r.GET("/cronjobs/:namespace/:name", getCronJobDetail(logger, getK8sClient))
}

// getCronJobList 获取CronJob列表的处理函数
// @Summary 获取 CronJob 列表
// @Description 获取CronJob列表，支持分页
// @Tags CronJob
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /cronjobs [get]
func getCronJobList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listCronJobs func(context.Context, *kubernetes.Clientset, string) ([]model.CronJobStatus, error),
) gin.HandlerFunc {
	return GenericListHandler(logger, getK8sClient, listCronJobs)
}

// getCronJobDetail 获取CronJob详情的处理函数
// @Summary 获取 CronJob 详情
// @Description 获取指定命名空间下的CronJob详情
// @Tags CronJob
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "CronJob 名称"
// @Success 200 {object} model.APIResponse
// @Router /cronjobs/{namespace}/{name} [get]
func getCronJobDetail(
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
		cronjob, err := clientset.BatchV1().CronJobs(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		status := "Unknown"
		if cronjob.Status.Active != nil && len(cronjob.Status.Active) > 0 {
			status = "Running"
		} else if cronjob.Status.LastSuccessfulTime != nil {
			status = "Succeeded"
		} else {
			status = "Pending"
		}

		image := ""
		if len(cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers) > 0 {
			image = cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image
		}

		// 处理时间字段
		lastScheduleTime := ""
		if cronjob.Status.LastScheduleTime != nil {
			lastScheduleTime = cronjob.Status.LastScheduleTime.Format("2006-01-02 15:04:05")
		}

		cronJobDetail := model.CronJobDetail{
			Namespace:        cronjob.Namespace,
			Name:             cronjob.Name,
			Schedule:         cronjob.Spec.Schedule,
			Suspend:          *cronjob.Spec.Suspend,
			Active:           len(cronjob.Status.Active),
			LastScheduleTime: lastScheduleTime,
			Status:           status,
			Labels:           cronjob.Labels,
			Annotations:      cronjob.Annotations,
			Image:            image,
		}
		middleware.ResponseSuccess(c, cronJobDetail, "success", nil)
	}
}
