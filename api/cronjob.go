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
// @Description 获取CronJob列表，支持分页和搜索
// @Tags CronJob
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Param search query string false "搜索关键词（支持名称、命名空间、状态等字段搜索）"
// @Success 200 {object} model.APIResponse
// @Router /cronjobs [get]
func getCronJobList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listCronJobs func(context.Context, *kubernetes.Clientset, string) ([]model.CronJobStatus, error),
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

		cronJobs, err := listCronJobs(ctx, clientset, namespace)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 新增：如果提供了搜索关键词，先进行搜索过滤
		var filteredCronJobs []model.CronJobStatus
		if search != "" {
			filteredCronJobs = filterCronJobsBySearch(cronJobs, search)
		} else {
			filteredCronJobs = cronJobs
		}

		// 对过滤后的数据进行分页
		paged := Paginate(filteredCronJobs, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(filteredCronJobs), // 使用过滤后的总数
			Limit:  limit,
			Offset: offset,
		})
	}
}

// filterCronJobsBySearch 根据搜索关键词过滤CronJob
func filterCronJobsBySearch(cronJobs []model.CronJobStatus, search string) []model.CronJobStatus {
	if search == "" {
		return cronJobs
	}
	searchLower := strings.ToLower(search)
	var filtered []model.CronJobStatus
	for _, cronJob := range cronJobs {
		if strings.Contains(strings.ToLower(cronJob.Name), searchLower) ||
			strings.Contains(strings.ToLower(cronJob.Namespace), searchLower) ||
			strings.Contains(strings.ToLower(cronJob.Schedule), searchLower) ||
			strings.Contains(strings.ToLower(cronJob.Status), searchLower) ||
			strings.Contains(strings.ToLower(strconv.Itoa(cronJob.Active)), searchLower) ||
			strings.Contains(strings.ToLower(cronJob.LastScheduleTime), searchLower) {
			filtered = append(filtered, cronJob)
		}
	}
	return filtered
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
		if len(cronjob.Status.Active) > 0 {
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
