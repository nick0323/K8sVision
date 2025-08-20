package service

import (
	"context"

	"github.com/nick0323/K8sVision/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListDeployments 采集指定命名空间（或全部）下的 Deployment 信息，返回 DeploymentStatus 列表
// namespace 为空时返回所有命名空间
func ListDeployments(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]model.DeploymentStatus, error) {
	depList, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]model.DeploymentStatus, 0, len(depList.Items))
	for _, d := range depList.Items {
		status := GetWorkloadStatus(d.Status.ReadyReplicas, d.Status.Replicas)
		result = append(result, model.DeploymentStatus{
			Namespace: d.Namespace,
			Name:      d.Name,
			Available: d.Status.ReadyReplicas,
			Desired:   d.Status.Replicas,
			Status:    status,
		})
	}
	return result, nil
}
