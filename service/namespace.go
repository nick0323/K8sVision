package service

import (
	"context"

	"github.com/nick0323/K8sVision/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListNamespaces 采集指定 context 下的所有命名空间详细信息
func ListNamespaces(ctx context.Context, clientset *kubernetes.Clientset) ([]model.NamespaceDetail, error) {
	nsList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]model.NamespaceDetail, 0, len(nsList.Items))
	for _, ns := range nsList.Items {
		result = append(result, model.NamespaceDetail{
			Name:        ns.Name,
			Status:      string(ns.Status.Phase),
			Labels:      ns.Labels,
			Annotations: ns.Annotations,
		})
	}
	return result, nil
}
