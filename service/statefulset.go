package service

import (
	"context"

	"github.com/nick0323/K8sVision/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListStatefulSets(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]model.StatefulSetStatus, error) {
	stsList, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]model.StatefulSetStatus, 0, len(stsList.Items))
	for _, s := range stsList.Items {
		status := GetWorkloadStatus(s.Status.ReadyReplicas, s.Status.Replicas)
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
