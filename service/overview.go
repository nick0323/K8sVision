package service

import (
	"context"
	"math"
	"sort"
	"sync"

	"github.com/nick0323/K8sVision/model"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 保留一位小数
func round1(val float64) float64 {
	return math.Round(val*10) / 10
}

// 资源用量统计（节点Capacity、Pod Requests/Limits）
func calcResourceUsage(nodesRaw *corev1.NodeList, podList *corev1.PodList) (cpuCap, memCap, cpuReq, cpuLim, memReq, memLim float64) {
	for _, node := range nodesRaw.Items {
		if c, ok := node.Status.Capacity["cpu"]; ok {
			cpuCap += float64(c.MilliValue()) / 1000.0
		}
		if m, ok := node.Status.Capacity["memory"]; ok {
			memCap += float64(m.Value()) / (1024 * 1024 * 1024)
		}
	}
	if podList != nil {
		for _, pod := range podList.Items {
			for _, c := range pod.Spec.Containers {
				if c.Resources.Requests != nil {
					if v, ok := c.Resources.Requests["cpu"]; ok {
						cpuReq += float64(v.MilliValue()) / 1000.0
					}
					if v, ok := c.Resources.Requests["memory"]; ok {
						memReq += float64(v.Value()) / (1024 * 1024 * 1024)
					}
				}
				if c.Resources.Limits != nil {
					if v, ok := c.Resources.Limits["cpu"]; ok {
						cpuLim += float64(v.MilliValue()) / 1000.0
					}
					if v, ok := c.Resources.Limits["memory"]; ok {
						memLim += float64(v.Value()) / (1024 * 1024 * 1024)
					}
				}
			}
		}
	}
	return round1(cpuCap), round1(memCap), round1(cpuReq), round1(cpuLim), round1(memReq), round1(memLim)
}

// GetOverviewStatus 采集集群总览信息，复用各资源采集逻辑，并发采集、减少API调用
func GetOverviewStatus(clientset *kubernetes.Clientset) (*model.OverviewStatus, error) {
	overview := &model.OverviewStatus{}
	ctx := context.Background()

	var (
		pods    []model.PodStatus
		podList *corev1.PodList
		podsErr error

		nodesRaw *corev1.NodeList
		nodesErr error

		nsList []model.NamespaceDetail
		nsErr  error

		services    []model.ServiceStatus
		servicesErr error
	)

	wg := sync.WaitGroup{}
	wg.Add(4)

	// 并发采集 Pod
	go func() {
		defer wg.Done()
		pods, podList, podsErr = ListPodsWithRaw(ctx, clientset, nil, "")
	}()

	// 并发采集节点原始数据
	go func() {
		defer wg.Done()
		nodesRaw, nodesErr = clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{})
	}()

	// 并发采集 Namespace
	go func() {
		defer wg.Done()
		nsList, nsErr = ListNamespaces(ctx, clientset)
	}()

	// 并发采集 Service
	go func() {
		defer wg.Done()
		services, servicesErr = ListServices(ctx, clientset, "")
	}()

	wg.Wait()

	// 节点状态统计（只用一次API结果）
	var nodes []model.NodeStatus
	var nodesStatusErr error
	if nodesErr == nil && podList != nil {
		nodes, nodesStatusErr = ListNodes(ctx, clientset, podList, nil)
	}

	// Pod统计
	if podsErr == nil {
		overview.PodCount = len(pods)
		podNotReady := 0
		for _, p := range pods {
			if p.Status != "Running" && p.Status != "Succeeded" {
				podNotReady++
			}
		}
		overview.PodNotReady = podNotReady
	}

	// Node统计
	if nodesStatusErr == nil {
		overview.NodeCount = len(nodes)
		nodeReady := 0
		for _, n := range nodes {
			if n.Status == "Active" {
				nodeReady++
			}
		}
		overview.NodeReady = nodeReady
	}

	// Namespace统计
	if nsErr == nil {
		overview.NamespaceCount = len(nsList)
	}

	// Service统计
	if servicesErr == nil {
		overview.ServiceCount = len(services)
	}

	// 资源用量统计
	if nodesRaw != nil && podList != nil {
		cpuCap, memCap, cpuReq, cpuLim, memReq, memLim := calcResourceUsage(nodesRaw, podList)
		overview.CPUCapacity = cpuCap
		overview.MemoryCapacity = memCap
		overview.CPURequests = cpuReq
		overview.CPULimits = cpuLim
		overview.MemoryRequests = memReq
		overview.MemoryLimits = memLim
	}

	// recent events 聚合
	eventList, err := clientset.CoreV1().Events("").List(ctx, v1.ListOptions{})
	if err == nil && eventList != nil {
		events := eventList.Items
		sort.Slice(events, func(i, j int) bool {
			return events[i].LastTimestamp.Time.After(events[j].LastTimestamp.Time)
		})
		var recentEvents []model.EventStatus
		for i, e := range events {
			if i >= 5 {
				break
			}
			recentEvents = append(recentEvents, model.EventStatus{
				Namespace: e.Namespace,
				Name:      e.Name,
				Reason:    e.Reason,
				Message:   e.Message,
				Type:      e.Type,
				Count:     e.Count,
				FirstSeen: model.FormatTime(&e.FirstTimestamp),
				LastSeen:  model.FormatTime(&e.LastTimestamp),
			})
		}
		overview.Events = recentEvents
	}

	return overview, nil
}
