package service

import (
	"context"

	"github.com/nick0323/K8sVision/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListJobs 采集指定 context 下的 Job 信息，返回 JobStatus 列表
func ListJobs(ctx context.Context, clientset *kubernetes.Clientset) ([]model.JobStatus, error) {
	jobs, err := clientset.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]model.JobStatus, 0, len(jobs.Items))
	for _, j := range jobs.Items {
		status := "Unknown"
		if j.Status.Succeeded > 0 {
			status = "Succeeded"
		} else if j.Status.Failed > 0 {
			status = "Failed"
		} else if j.Status.Active > 0 {
			status = "Active"
		}
		startTime := model.FormatTime(j.Status.StartTime)
		completionTime := model.FormatTime(j.Status.CompletionTime)
		result = append(result, model.JobStatus{
			Namespace:      j.Namespace,
			Name:           j.Name,
			Completions:    *j.Spec.Completions,
			Succeeded:      j.Status.Succeeded,
			Failed:         j.Status.Failed,
			StartTime:      startTime,
			CompletionTime: completionTime,
			Status:         status,
		})
	}
	return result, nil
}
