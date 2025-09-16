package service

import (
	"context"

	"github.com/nick0323/K8sVision/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListJobs(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]model.JobStatus, error) {
	jobs, err := clientset.BatchV1().Jobs(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]model.JobStatus, 0, len(jobs.Items))
	for _, job := range jobs.Items {
		status := GetJobStatus(job.Status.Succeeded, job.Status.Failed, job.Status.Active)

		startTime := ""
		if job.Status.StartTime != nil {
			startTime = job.Status.StartTime.Format("2006-01-02 15:04:05")
		}

		completionTime := ""
		if job.Status.CompletionTime != nil {
			completionTime = job.Status.CompletionTime.Format("2006-01-02 15:04:05")
		}

		result = append(result, model.JobStatus{
			Namespace:      job.Namespace,
			Name:           job.Name,
			Completions:    SafeInt32Ptr(job.Spec.Completions, 0),
			Succeeded:      job.Status.Succeeded,
			Failed:         job.Status.Failed,
			StartTime:      startTime,
			CompletionTime: completionTime,
			Status:         status,
		})
	}
	return result, nil
}
