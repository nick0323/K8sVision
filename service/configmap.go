package service

import (
	"context"

	"github.com/nick0323/K8sVision/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListConfigMaps 获取 ConfigMap 列表
func ListConfigMaps(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]model.ConfigMapStatus, error) {
	var cmList *v1.ConfigMapList
	var err error

	if namespace == "" {
		cmList, err = clientset.CoreV1().ConfigMaps("").List(ctx, metav1.ListOptions{})
	} else {
		cmList, err = clientset.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	}

	if err != nil {
		return nil, err
	}

	var cmStatuses []model.ConfigMapStatus
	for _, cm := range cmList.Items {
		cmStatus := model.ConfigMapStatus{
			Namespace: cm.Namespace,
			Name:      cm.Name,
			DataCount: len(cm.Data),
			Keys:      ExtractKeys(cm.Data),
		}
		cmStatuses = append(cmStatuses, cmStatus)
	}

	return cmStatuses, nil
}
