package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"
	"github.com/nick0323/K8sVision/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterJob 注册 Job 相关路由
func RegisterJob(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listJobs func(context.Context, *kubernetes.Clientset, string) ([]model.JobStatus, error),
) {
	r.GET("/jobs", getJobList(logger, getK8sClient, listJobs))
	r.GET("/jobs/:namespace/:name", getJobDetail(logger, getK8sClient))
}

// getJobList 获取Job列表的处理函数
// @Summary 获取 Job 列表
// @Description 获取Job列表，支持分页和搜索
// @Tags Job
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Param search query string false "搜索关键词（支持名称、命名空间、状态等字段搜索）"
// @Success 200 {object} model.APIResponse
// @Router /jobs [get]
func getJobList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listJobs func(context.Context, *kubernetes.Clientset, string) ([]model.JobStatus, error),
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

		jobs, err := listJobs(ctx, clientset, namespace)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 新增：如果提供了搜索关键词，先进行搜索过滤
		var filteredJobs []model.JobStatus
		if search != "" {
			filteredJobs = filterJobsBySearch(jobs, search)
		} else {
			filteredJobs = jobs
		}

		// 对过滤后的数据进行分页
		paged := Paginate(filteredJobs, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(filteredJobs), // 使用过滤后的总数
			Limit:  limit,
			Offset: offset,
		})
	}
}

// filterJobsBySearch 根据搜索关键词过滤Job
func filterJobsBySearch(jobs []model.JobStatus, search string) []model.JobStatus {
	if search == "" {
		return jobs
	}
	searchLower := strings.ToLower(search)
	var filtered []model.JobStatus
	for _, job := range jobs {
		if strings.Contains(strings.ToLower(job.Name), searchLower) ||
			strings.Contains(strings.ToLower(job.Namespace), searchLower) ||
			strings.Contains(strings.ToLower(job.Status), searchLower) ||
			strings.Contains(strings.ToLower(strconv.Itoa(int(job.Completions))), searchLower) ||
			strings.Contains(strings.ToLower(strconv.Itoa(int(job.Succeeded))), searchLower) ||
			strings.Contains(strings.ToLower(strconv.Itoa(int(job.Failed))), searchLower) ||
			strings.Contains(strings.ToLower(job.StartTime), searchLower) ||
			strings.Contains(strings.ToLower(job.CompletionTime), searchLower) {
			filtered = append(filtered, job)
		}
	}
	return filtered
}

// getJobDetail 获取Job详情的处理函数
// @Summary 获取 Job 详情
// @Description 获取指定命名空间下的Job详情
// @Tags Job
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "Job 名称"
// @Success 200 {object} model.APIResponse
// @Router /jobs/{namespace}/{name} [get]
func getJobDetail(
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
		job, err := clientset.BatchV1().Jobs(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		image := ""
		if len(job.Spec.Template.Spec.Containers) > 0 {
			image = job.Spec.Template.Spec.Containers[0].Image
		}

		// 处理时间字段
		startTime := ""
		if job.Status.StartTime != nil {
			startTime = job.Status.StartTime.Format("2006-01-02 15:04:05")
		}
		
		completionTime := ""
		if job.Status.CompletionTime != nil {
			completionTime = job.Status.CompletionTime.Format("2006-01-02 15:04:05")
		}

		jobDetail := model.JobDetail{
			Namespace:      job.Namespace,
			Name:           job.Name,
			Completions:    *job.Spec.Completions,
			Succeeded:      job.Status.Succeeded,
			Failed:         job.Status.Failed,
			StartTime:      startTime,
			CompletionTime: completionTime,
			Status:         service.GetJobStatus(job.Status.Succeeded, job.Status.Failed, job.Status.Active),
			Labels:         job.Labels,
			Annotations:    job.Annotations,
			Image:          image,
		}
		middleware.ResponseSuccess(c, jobDetail, "success", nil)
	}
}
