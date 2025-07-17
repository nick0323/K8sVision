package service

import (
	"context"

	"github.com/nick0323/K8sVision/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListStatefulSets 采集指定命名空间（或全部）下的 StatefulSet 信息，返回 StatefulSetStatus 列表
// namespace 为空时返回所有命名空间
func ListStatefulSets(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]model.StatefulSetStatus, error) {
	stsList, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]model.StatefulSetStatus, 0, len(stsList.Items))
	for _, s := range stsList.Items {
		status := "Unknown"
		if s.Status.ReadyReplicas == s.Status.Replicas {
			status = "Healthy"
		} else if s.Status.ReadyReplicas > 0 {
			status = "PartialAvailable"
		} else {
			status = "Abnormal"
		}
		result = append(result, model.StatefulSetStatus{
			Namespace: s.Namespace,
			Name:      s.Name,
			Available: s.Status.ReadyReplicas,
			Desired:   s.Status.Replicas,
			Status:    status,
		})
	}
	return result, nil
}
