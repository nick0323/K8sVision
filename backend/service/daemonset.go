package service

import (
	"context"

	"github.com/nick0323/K8sVision/backend/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListDaemonSets 采集指定命名空间（或全部）下的 DaemonSet 信息，返回 DaemonSetStatus 列表
// namespace 为空时返回所有命名空间
func ListDaemonSets(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]model.DaemonSetStatus, error) {
	daemonList, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]model.DaemonSetStatus, 0, len(daemonList.Items))
	for _, d := range daemonList.Items {
		status := "Unknown"
		if d.Status.NumberReady == d.Status.DesiredNumberScheduled {
			status = "Healthy"
		} else if d.Status.NumberReady > 0 {
			status = "PartialAvailable"
		} else {
			status = "Abnormal"
		}
		result = append(result, model.DaemonSetStatus{
			Namespace: d.Namespace,
			Name:      d.Name,
			Available: d.Status.NumberReady,
			Desired:   d.Status.DesiredNumberScheduled,
			Status:    status,
		})
	}
	return result, nil
}
