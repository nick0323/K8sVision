package service

import (
	"context"

	"github.com/nick0323/K8sVision/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListSecrets 获取 Secret 列表
func ListSecrets(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]model.SecretStatus, error) {
	var secretList *v1.SecretList
	var err error

	if namespace == "" {
		secretList, err = clientset.CoreV1().Secrets("").List(ctx, metav1.ListOptions{})
	} else {
		secretList, err = clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	}

	if err != nil {
		return nil, err
	}

	var secretStatuses []model.SecretStatus
	for _, secret := range secretList.Items {
		secretStatus := model.SecretStatus{
			Namespace: secret.Namespace,
			Name:      secret.Name,
			Type:      string(secret.Type),
			DataCount: len(secret.Data),
			Keys:      ExtractKeys(secret.Data),
		}
		secretStatuses = append(secretStatuses, secretStatus)
	}

	return secretStatuses, nil
}
