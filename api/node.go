package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterNode 注册 Node 相关路由，包括列表和详情接口
func RegisterNode(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listPodsWithRaw func(context.Context, *kubernetes.Clientset, model.PodMetricsMap, string) ([]model.PodStatus, *v1.PodList, error),
	listNodes func(context.Context, *kubernetes.Clientset, *v1.PodList, model.NodeMetricsMap) ([]model.NodeStatus, error),
) {
	r.GET("/nodes", getNodeList(logger, getK8sClient, listPodsWithRaw, listNodes))
	r.GET("/nodes/:name", getNodeDetail(logger, getK8sClient))
}

// getNodeList 获取Node列表的处理函数
// @Summary 获取 Node 列表
// @Description 获取集群节点列表，支持分页
// @Tags Node
// @Security BearerAuth
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /nodes [get]
func getNodeList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listPodsWithRaw func(context.Context, *kubernetes.Clientset, model.PodMetricsMap, string) ([]model.PodStatus, *v1.PodList, error),
	listNodes func(context.Context, *kubernetes.Clientset, *v1.PodList, model.NodeMetricsMap) ([]model.NodeStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientset, metricsClient, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		search := c.DefaultQuery("search", "") // 新增：搜索关键词
		
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
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 新增：如果提供了搜索关键词，先进行搜索过滤
		var filteredNodes []model.NodeStatus
		if search != "" {
			filteredNodes = filterNodesBySearch(nodeStatuses, search)
		} else {
			filteredNodes = nodeStatuses
		}

		// 对过滤后的数据进行分页
		total := len(filteredNodes)
		paged := Paginate(filteredNodes, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		})
	}
}

// getNodeDetail 获取Node详情的处理函数
// @Summary 获取 Node 详情
// @Description 获取指定节点的详细信息
// @Tags Node
// @Security BearerAuth
// @Param name path string true "Node 名称"
// @Success 200 {object} model.APIResponse
// @Router /nodes/{name} [get]
func getNodeDetail(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientset, metricsClient, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		name := c.Param("name")
		node, err := clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
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
		ip := ""
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				ip = addr.Address
				break
			}
		}
		roles := make([]string, 0)
		for key := range node.Labels {
			if strings.HasPrefix(key, "node-role.kubernetes.io/") {
				role := strings.TrimPrefix(key, "node-role.kubernetes.io/")
				if role == "" {
					role = "worker"
				}
				roles = append(roles, role)
			}
		}
		// 如果没有找到角色标签，默认为worker
		if len(roles) == 0 {
			roles = append(roles, "worker")
		}
		metricsList, _ := metricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, name, metav1.GetOptions{})
		var cpuPercent, memPercent float64
		if metricsList != nil {
			cpuUsed := metricsList.Usage.Cpu().MilliValue()
			memUsed := metricsList.Usage.Memory().Value()
			cpuTotal := node.Status.Allocatable.Cpu().MilliValue()
			memTotal := node.Status.Allocatable.Memory().Value()
			if cpuTotal > 0 {
				cpuPercent = float64(cpuUsed) / float64(cpuTotal) * 100
			}
			if memTotal > 0 {
				memPercent = float64(memUsed) / float64(memTotal) * 100
			}
		}
		nodeDetail := model.NodeDetail{
			Name:         node.Name,
			IP:           ip,
			Status:       string(node.Status.Conditions[len(node.Status.Conditions)-1].Type),
			CPUUsage:     cpuPercent,
			MemoryUsage:  memPercent,
			Role:         roles,
			PodsUsed:     podsUsed,
			PodsCapacity: podsCapacity,
			Labels:       node.Labels,
			Annotations:  node.Annotations,
		}
		middleware.ResponseSuccess(c, nodeDetail, "success", nil)
	}
}

// filterNodesBySearch 根据搜索关键词过滤Node
func filterNodesBySearch(nodes []model.NodeStatus, search string) []model.NodeStatus {
	if search == "" {
		return nodes
	}

	searchLower := strings.ToLower(search)
	var filtered []model.NodeStatus

	for _, node := range nodes {
		// 检查Node的各个字段是否匹配搜索关键词
		if strings.Contains(strings.ToLower(node.Name), searchLower) ||
			strings.Contains(strings.ToLower(node.Status), searchLower) ||
			strings.Contains(strings.ToLower(node.IP), searchLower) {
			filtered = append(filtered, node)
		}
	}

	return filtered
}
