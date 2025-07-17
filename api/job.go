package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterJob 注册 Job 相关路由，包括列表和详情接口
// @Summary 获取 Job 列表
// @Description 支持分页
// @Tags Job
// @Security BearerAuth
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /jobs [get]
//
// @Summary 获取 Job 详情
// @Description 获取指定命名空间下的 Job 详情
// @Tags Job
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "Job 名称"
// @Success 200 {object} model.APIResponse
// @Router /jobs/{namespace}/{name} [get]
func RegisterJob(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listJobs func(context.Context, *kubernetes.Clientset) ([]model.JobStatus, error),
) {
	r.GET("/jobs", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		jobs, err := listJobs(ctx, clientset)
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		paged := Paginate(jobs, offset, limit)
		ResponseOK(c, paged, "success", &model.PageMeta{
			Total:  len(jobs),
			Limit:  limit,
			Offset: offset,
		})
	})

	r.GET("/jobs/:namespace/:name", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		ns := c.Param("namespace")
		name := c.Param("name")
		job, err := clientset.BatchV1().Jobs(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			ResponseError(c, logger, err, http.StatusNotFound)
			return
		}
		ResponseOK(c, job, "success", nil)
	})
}
