package service

import (
	"context"

	"github.com/nick0323/K8sVision/backend/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListCronJobs 采集指定 context 下的 CronJob 信息，返回 CronJobStatus 列表
func ListCronJobs(ctx context.Context, clientset *kubernetes.Clientset) ([]model.CronJobStatus, error) {
	cronjobs, err := clientset.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]model.CronJobStatus, 0, len(cronjobs.Items))
	for _, c := range cronjobs.Items {
		status := ""
		if len(c.Status.Active) > 0 {
			status = string(c.Status.Active[0].Kind)
		}
		lastSchedule := model.FormatTime(c.Status.LastScheduleTime)
		result = append(result, model.CronJobStatus{
			Namespace:        c.Namespace,
			Name:             c.Name,
			Schedule:         c.Spec.Schedule,
			LastScheduleTime: lastSchedule,
			Status:           status,
		})
	}
	return result, nil
}
