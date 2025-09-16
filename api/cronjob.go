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
)

func RegisterCronJob(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listCronJobs func(context.Context, *kubernetes.Clientset, string) ([]model.CronJobStatus, error),
) {
	r.GET("/cronjobs", getCronJobList(logger, getK8sClient, listCronJobs))
	r.GET("/cronjobs/:namespace/:name", getCronJobDetail(logger, getK8sClient))
}

func getCronJobList(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listCronJobs func(context.Context, *kubernetes.Clientset, string) ([]model.CronJobStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		HandleListWithPagination(c, logger, func(ctx context.Context, params PaginationParams) ([]model.CronJobStatus, error) {
			clientset, _, err := getK8sClient()
			if err != nil {
				return nil, err
			}
			return listCronJobs(ctx, clientset, params.Namespace)
		}, ListSuccessMessage)
	}
}

func getCronJobDetail(
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

		lastScheduleTime := ""
		if cronjob.Status.LastScheduleTime != nil {
			lastScheduleTime = cronjob.Status.LastScheduleTime.Format("2006-01-02 15:04:05")
		}

		cronJobDetail := model.CronJobDetail{
			CommonResourceFields: model.CommonResourceFields{
				Namespace: cronjob.Namespace,
				Name:      cronjob.Name,
				Status:    status,
				BaseMetadata: model.BaseMetadata{
					Labels:      cronjob.Labels,
					Annotations: cronjob.Annotations,
				},
			},
			Schedule:         cronjob.Spec.Schedule,
			Suspend:          *cronjob.Spec.Suspend,
			Active:           len(cronjob.Status.Active),
			LastScheduleTime: lastScheduleTime,
			Image:            image,
		}
		middleware.ResponseSuccess(c, cronJobDetail, DetailSuccessMessage, nil)
	}
}
