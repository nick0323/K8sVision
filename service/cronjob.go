package service

import (
	"context"

	"github.com/nick0323/K8sVision/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListCronJobs 采集 CronJob 信息，支持命名空间过滤
func ListCronJobs(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]model.CronJobStatus, error) {
	cronjobs, err := clientset.BatchV1().CronJobs(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]model.CronJobStatus, 0, len(cronjobs.Items))
	for _, cj := range cronjobs.Items {
		status := GetCronJobStatus(len(cj.Status.Active), cj.Status.LastSuccessfulTime)
		
		// 处理时间字段
		lastScheduleTime := ""
		if cj.Status.LastScheduleTime != nil {
			lastScheduleTime = cj.Status.LastScheduleTime.Format("2006-01-02 15:04:05")
		}
		
		result = append(result, model.CronJobStatus{
			Namespace:        cj.Namespace,
			Name:             cj.Name,
			Schedule:         cj.Spec.Schedule,
			Suspend:          SafeBoolPtr(cj.Spec.Suspend, false),
			Active:           len(cj.Status.Active),
			LastScheduleTime: lastScheduleTime,
			Status:           status,
		})
	}
	return result, nil
}
