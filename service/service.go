package service

import (
	"context"
	"fmt"

	"github.com/nick0323/K8sVision/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListServices 采集 Service 信息，返回 ServiceStatus 列表
func ListServices(ctx context.Context, clientset *kubernetes.Clientset) ([]model.ServiceStatus, error) {
	svcs, err := clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	svcStatuses := make([]model.ServiceStatus, 0, len(svcs.Items))
	for _, svc := range svcs.Items {
		ports := make([]string, 0, len(svc.Spec.Ports))
		for _, p := range svc.Spec.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", p.Port, p.Protocol))
		}
		svcStatuses = append(svcStatuses, model.ServiceStatus{
			Namespace: svc.Namespace,
			Name:      svc.Name,
			Type:      string(svc.Spec.Type),
			ClusterIP: svc.Spec.ClusterIP,
			Ports:     ports,
		})
	}
	return svcStatuses, nil
}
