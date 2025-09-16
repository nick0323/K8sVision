package service

import (
	"context"

	"github.com/nick0323/K8sVision/model"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListPodsWithRaw(ctx context.Context, clientset *kubernetes.Clientset, podMetricsMap model.PodMetricsMap, namespace string) ([]model.PodStatus, *v1.PodList, error) {
	pods, err := ListResourcesWithNamespace(ctx, namespace,
		func() (*v1.PodList, error) {
			return clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
		},
		func(ns string) (*v1.PodList, error) {
			return clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
		},
	)
	if err != nil {
		return nil, nil, err
	}

	podStatuses := make([]model.PodStatus, 0, len(pods.Items))
	for _, pod := range pods.Items {
		cpuVal, memVal := FormatPodResourceUsage(podMetricsMap, pod.Namespace, pod.Name)

		podStatuses = append(podStatuses, model.PodStatus{
			Namespace:   pod.Namespace,
			Name:        pod.Name,
			Status:      string(pod.Status.Phase),
			CPUUsage:    cpuVal,
			MemoryUsage: memVal,
			PodIP:       pod.Status.PodIP,
			NodeName:    pod.Spec.NodeName,
		})
	}
	return podStatuses, pods, nil
}
