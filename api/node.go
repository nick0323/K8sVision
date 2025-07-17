package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/nick0323/K8sVision/model"
	"github.com/nick0323/K8sVision/service"

	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterNode 注册 Node 相关路由，包括列表和详情接口
// @Summary 获取 Node 列表
// @Description 获取集群节点列表
// @Tags Node
// @Security BearerAuth
// @Success 200 {object} model.APIResponse
// @Router /nodes [get]
//
// @Summary 获取 Node 详情
// @Description 获取指定节点详情
// @Tags Node
// @Security BearerAuth
// @Param name path string true "Node 名称"
// @Success 200 {object} model.APIResponse
// @Router /nodes/{name} [get]
func RegisterNode(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listPodsWithRaw func(context.Context, *kubernetes.Clientset, model.PodMetricsMap, string) ([]model.PodStatus, *v1.PodList, error),
	listNodes func(context.Context, *kubernetes.Clientset, *v1.PodList, model.NodeMetricsMap) ([]model.NodeStatus, error),
) {
	r.GET("/nodes", func(c *gin.Context) {
		clientset, metricsClient, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		podMetricsList, _ := metricsClient.MetricsV1beta1().PodMetricses("").List(ctx, metav1.ListOptions{})
		podMetricsMap := make(model.PodMetricsMap)
		if podMetricsList != nil {
			// 这里应根据实际类型断言和处理
		}
		_, podList, _ := listPodsWithRaw(ctx, clientset, podMetricsMap, "")
		metricsList, _ := metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
		nodeMetricsMap := make(model.NodeMetricsMap)
		if metricsList != nil {
			for _, m := range metricsList.Items {
				cpu := m.Usage.Cpu().String()
				mem := m.Usage.Memory().String()
				nodeMetricsMap[m.Name] = model.NodeMetrics{CPU: cpu, Mem: mem}
			}
		}
		nodeStatuses, err := listNodes(ctx, clientset, podList, nodeMetricsMap)
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		total := len(nodeStatuses)
		paged := Paginate(nodeStatuses, offset, limit)
		ResponseOK(c, paged, "success", &model.PageMeta{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		})
	})

	r.GET("/nodes/:name", func(c *gin.Context) {
		clientset, metricsClient, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		name := c.Param("name")
		node, err := clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			ResponseError(c, logger, err, http.StatusNotFound)
			return
		}
		pods, _ := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
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
		nodeMetricsMap := make(model.NodeMetricsMap)
		metricsList, _ := metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
		if metricsList != nil {
			for _, m := range metricsList.Items {
				cpu := m.Usage.Cpu().String()
				mem := m.Usage.Memory().String()
				nodeMetricsMap[m.Name] = model.NodeMetrics{CPU: cpu, Mem: mem}
			}
		}
		metric := nodeMetricsMap[node.Name]
		cpuUsed := service.ParseCPU(metric.CPU)
		memUsed := service.ParseMemory(metric.Mem)
		cpuTotal := service.GetNodeAllocatableCPU(*node)
		memTotal := service.GetNodeAllocatableMemory(*node)
		cpuPercent := 0.0
		memPercent := 0.0
		if cpuTotal > 0 {
			cpuPercent = cpuUsed / cpuTotal * 100
		}
		if memTotal > 0 {
			memPercent = memUsed / memTotal * 100
		}
		nodeStatus := model.NodeStatus{
			Name:         node.Name,
			IP:           ip,
			Status:       status,
			CPUUsage:     math.Round(cpuPercent*10) / 10,
			MemoryUsage:  math.Round(memPercent*10) / 10,
			Role:         roles,
			PodsUsed:     podsUsed,
			PodsCapacity: podsCapacity,
		}
		ResponseOK(c, nodeStatus, "success", nil)
	})
}
