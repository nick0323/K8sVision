package service

import (
	"context"

	"github.com/nick0323/K8sVision/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListDaemonSets 采集 DaemonSet 信息，支持命名空间过滤
func ListDaemonSets(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]model.DaemonSetStatus, error) {
	dsList, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]model.DaemonSetStatus, 0, len(dsList.Items))
	for _, ds := range dsList.Items {
		status := GetResourceStatus(ds.Status.NumberReady, ds.Status.DesiredNumberScheduled)
		result = append(result, model.DaemonSetStatus{
			Namespace: ds.Namespace,
			Name:      ds.Name,
			Available: ds.Status.NumberReady,
			Desired:   ds.Status.DesiredNumberScheduled,
			Status:    status,
		})
	}
	return result, nil
}
