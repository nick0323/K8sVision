package api

import (
	"context"
	"net/http"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"
	"github.com/nick0323/K8sVision/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func RegisterJob(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listJobs func(context.Context, *kubernetes.Clientset, string) ([]model.JobStatus, error),
) {
	r.GET("/jobs", getJobList(logger, getK8sClient, listJobs))
	r.GET("/jobs/:namespace/:name", getJobDetail(logger, getK8sClient))
}

func getJobList(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listJobs func(context.Context, *kubernetes.Clientset, string) ([]model.JobStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		HandleListWithPagination(c, logger, func(ctx context.Context, params PaginationParams) ([]model.JobStatus, error) {
			clientset, _, err := getK8sClient()
			if err != nil {
				return nil, err
			}
			return listJobs(ctx, clientset, params.Namespace)
		}, ListSuccessMessage)
	}
}

func getJobDetail(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := GetRequestContext(c)
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

		startTime := ""
		if job.Status.StartTime != nil {
			startTime = job.Status.StartTime.Format("2006-01-02 15:04:05")
		}

		completionTime := ""
		if job.Status.CompletionTime != nil {
			completionTime = job.Status.CompletionTime.Format("2006-01-02 15:04:05")
		}

		jobDetail := model.JobDetail{
			CommonResourceFields: model.CommonResourceFields{
				Namespace: job.Namespace,
				Name:      job.Name,
				Status:    service.GetJobStatus(job.Status.Succeeded, job.Status.Failed, job.Status.Active),
				BaseMetadata: model.BaseMetadata{
					Labels:      job.Labels,
					Annotations: job.Annotations,
				},
			},
			Completions:    *job.Spec.Completions,
			Succeeded:      job.Status.Succeeded,
			Failed:         job.Status.Failed,
			StartTime:      startTime,
			CompletionTime: completionTime,
			Image:          image,
		}
		middleware.ResponseSuccess(c, jobDetail, DetailSuccessMessage, nil)
	}
}
