package service

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/nick0323/K8sVision/backend/model"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ParseCPU: "123456n" -> float64(核)
func ParseCPU(cpuStr string) float64 {
	if cpuStr == "" {
		return 0
	}
	if strings.HasSuffix(cpuStr, "n") {
		n, _ := strconv.ParseFloat(strings.TrimSuffix(cpuStr, "n"), 64)
		return n / 1e9
	}
	if strings.HasSuffix(cpuStr, "m") {
		n, _ := strconv.ParseFloat(strings.TrimSuffix(cpuStr, "m"), 64)
		return n / 1000
	}
	n, _ := strconv.ParseFloat(cpuStr, 64)
	return n
}

// ParseMemory: "123456Ki" -> float64(GiB)
func ParseMemory(memStr string) float64 {
	if memStr == "" {
		return 0
	}
	if strings.HasSuffix(memStr, "Ki") {
		n, _ := strconv.ParseFloat(strings.TrimSuffix(memStr, "Ki"), 64)
		return n / 1048576
	}
	if strings.HasSuffix(memStr, "Mi") {
		n, _ := strconv.ParseFloat(strings.TrimSuffix(memStr, "Mi"), 64)
		return n / 1024
	}
	if strings.HasSuffix(memStr, "Gi") {
		n, _ := strconv.ParseFloat(strings.TrimSuffix(memStr, "Gi"), 64)
		return n
	}
	n, _ := strconv.ParseFloat(memStr, 64)
	return n
}

// 获取节点可分配 CPU（核）
func GetNodeAllocatableCPU(node v1.Node) float64 {
	if v, ok := node.Status.Allocatable["cpu"]; ok {
		return float64(v.MilliValue()) / 1000
	}
	return 0
}

// 获取节点可分配内存（GiB）
func GetNodeAllocatableMemory(node v1.Node) float64 {
	if v, ok := node.Status.Allocatable["memory"]; ok {
		return float64(v.Value()) / 1024 / 1024 / 1024
	}
	return 0
}

// ListNodes 采集 Node 信息，返回 NodeStatus 列表
func ListNodes(ctx context.Context, clientset *kubernetes.Clientset, pods *v1.PodList, nodeMetricsMap model.NodeMetricsMap) ([]model.NodeStatus, error) {
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodeStatuses := make([]model.NodeStatus, 0, len(nodes.Items))
	for _, node := range nodes.Items {
		status := "Unknown"
		ip := ""
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				ip = addr.Address
				break
			}
		}
		for _, cond := range node.Status.Conditions {
			if cond.Type == "Ready" && cond.Status == "True" {
				status = "Active"
				break
			}
		}
		roles := []string{}
		for k := range node.Labels {
			if strings.HasPrefix(k, "node-role.kubernetes.io/") {
				role := strings.TrimPrefix(k, "node-role.kubernetes.io/")
				if role == "" {
					role = "worker"
				}
				roles = append(roles, role)
			}
		}
		if len(roles) == 0 {
			roles = append(roles, "-")
		}
		podsUsed := 0
		for _, pod := range pods.Items {
			if pod.Spec.NodeName == node.Name {
				podsUsed++
			}
		}
		podsCapacity := 0
		if v, ok := node.Status.Allocatable["pods"]; ok {
			podsCapacity = int(v.Value())
		}
		metric := nodeMetricsMap[node.Name]
		cpuUsed := ParseCPU(metric.CPU)
		memUsed := ParseMemory(metric.Mem)
		cpuTotal := GetNodeAllocatableCPU(node)
		memTotal := GetNodeAllocatableMemory(node)
		cpuPercent := 0.0
		memPercent := 0.0
		if cpuTotal > 0 {
			cpuPercent = cpuUsed / cpuTotal * 100
		}
		if memTotal > 0 {
			memPercent = memUsed / memTotal * 100
		}
		nodeStatuses = append(nodeStatuses, model.NodeStatus{
			Name:         node.Name,
			IP:           ip,
			Status:       status,
			CPUUsage:     math.Round(cpuPercent*10) / 10,
			MemoryUsage:  math.Round(memPercent*10) / 10,
			Role:         roles,
			PodsUsed:     podsUsed,
			PodsCapacity: podsCapacity,
		})
	}
	return nodeStatuses, nil
}
