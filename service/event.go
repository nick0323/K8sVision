package service

import (
	"context"

	"github.com/nick0323/K8sVision/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListEvents 采集指定 context 下的 Event 信息，返回 EventStatus 列表
func ListEvents(ctx context.Context, clientset *kubernetes.Clientset) ([]model.EventStatus, error) {
	events, err := clientset.CoreV1().Events("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]model.EventStatus, 0, len(events.Items))
	for _, e := range events.Items {
		result = append(result, model.EventStatus{
			Namespace: e.Namespace,
			Name:      e.Name,
			Reason:    e.Reason,
			Message:   e.Message,
			Type:      e.Type,
			Count:     e.Count,
			FirstSeen: model.FormatTime(&e.FirstTimestamp),
			LastSeen:  model.FormatTime(&e.LastTimestamp),
			Duration:  e.LastTimestamp.Sub(e.FirstTimestamp.Time).String(),
		})
	}
	return result, nil
}
