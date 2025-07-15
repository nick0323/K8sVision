package service

import (
	"context"

	"github.com/nick0323/K8sVision/backend/model"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListCronJobs 采集指定 context 下的 CronJob 信息，返回 CronJobStatus 列表
func ListCronJobs(ctx context.Context, clientset *kubernetes.Clientset) ([]model.CronJobStatus, error) {
	cronjobs, err := clientset.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	// 查询所有 Job 以便判断 CronJob 最近一次 Job 状态
	jobs, _ := clientset.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	result := make([]model.CronJobStatus, 0, len(cronjobs.Items))
	for _, c := range cronjobs.Items {
		status := "Unknown"
		if c.Spec.Suspend != nil && *c.Spec.Suspend {
			status = "Suspended"
		} else if len(c.Status.Active) > 0 {
			status = "Active"
		} else {
			// 找到该 CronJob 关联的 Job（ownerReference 指向该 CronJob）
			var lastJob *batchv1.Job
			for i, job := range jobs.Items {
				for _, owner := range job.OwnerReferences {
					if owner.Kind == "CronJob" && owner.Name == c.Name && job.Namespace == c.Namespace {
						if lastJob == nil || job.Status.StartTime != nil && lastJob.Status.StartTime != nil && job.Status.StartTime.After(lastJob.Status.StartTime.Time) {
							lastJob = &jobs.Items[i]
						}
					}
				}
			}
			if lastJob != nil {
				if lastJob.Status.Failed > 0 {
					status = "Failed"
				} else if lastJob.Status.Succeeded > 0 {
					status = "Succeeded"
				} else {
					status = "Pending"
				}
			} else if c.Status.LastScheduleTime == nil {
				status = "Pending"
			}
		}
		lastSchedule := model.FormatTime(c.Status.LastScheduleTime)
		suspend := false
		if c.Spec.Suspend != nil {
			suspend = *c.Spec.Suspend
		}
		result = append(result, model.CronJobStatus{
			Namespace:        c.Namespace,
			Name:             c.Name,
			Schedule:         c.Spec.Schedule,
			Suspend:          suspend,
			Active:           len(c.Status.Active),
			LastScheduleTime: lastSchedule,
			Status:           status,
		})
	}
	return result, nil
}
