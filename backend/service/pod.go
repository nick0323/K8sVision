package service

import (
	"context"
	"fmt"

	"github.com/nick0323/K8sVision/backend/model"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListPodsWithRaw 采集 Pod 信息，返回 PodStatus 列表和原始 PodList
func ListPodsWithRaw(ctx context.Context, clientset *kubernetes.Clientset, podMetricsMap model.PodMetricsMap, namespace string) ([]model.PodStatus, *v1.PodList, error) {
	var pods *v1.PodList
	var err error
	if namespace == "" {
		pods, err = clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	} else {
		pods, err = clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	}
	if err != nil {
		return nil, nil, err
	}
	podStatuses := make([]model.PodStatus, 0, len(pods.Items))
	for _, pod := range pods.Items {
		cpuVal, memVal := "-", "-"
		if m, ok := podMetricsMap[pod.Namespace+"/"+pod.Name]; ok {
			cpuVal = fmt.Sprintf("%.2f mCPU", float64(m.CPU))
			memVal = fmt.Sprintf("%.2f MiB", float64(m.Mem)/(1024*1024))
		}
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
