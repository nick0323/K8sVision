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
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterPod 注册 Pod 相关路由，包括列表和详情接口
func RegisterPod(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listPodsWithRaw func(context.Context, *kubernetes.Clientset, model.PodMetricsMap, string) ([]model.PodStatus, *v1.PodList, error),
) {
	r.GET("/pods", getPodList(logger, getK8sClient, listPodsWithRaw))
	r.GET("/pods/:namespace/:name", getPodDetail(logger, getK8sClient))
}

// getPodList 获取Pod列表的处理函数
// @Summary 获取 Pod 列表
// @Description 获取Pod列表，支持分页和搜索
// @Tags Pod
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Param search query string false "搜索关键词（支持名称、命名空间、状态、PodIP、节点等字段搜索）"
// @Success 200 {object} model.APIResponse
// @Router /pods [get]
func getPodList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listPodsWithRaw func(context.Context, *kubernetes.Clientset, model.PodMetricsMap, string) ([]model.PodStatus, *v1.PodList, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientset, metricsClient, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		namespace := c.DefaultQuery("namespace", "")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		search := c.DefaultQuery("search", "") // 新增：搜索关键词
		
		metricsList, _ := metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
		podMetricsMap := make(model.PodMetricsMap)
		if metricsList != nil {
			for _, m := range metricsList.Items {
				var cpuSum, memSum int64
				for _, ctn := range m.Containers {
					cpuSum += ctn.Usage.Cpu().MilliValue()
					memSum += ctn.Usage.Memory().Value()
				}
				podMetricsMap[m.Namespace+"/"+m.Name] = model.PodMetrics{CPU: cpuSum, Mem: memSum}
			}
		}
		podStatuses, _, err := listPodsWithRaw(ctx, clientset, podMetricsMap, namespace)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 新增：如果提供了搜索关键词，先进行搜索过滤
		var filteredPods []model.PodStatus
		if search != "" {
			filteredPods = filterPodsBySearch(podStatuses, search)
		} else {
			filteredPods = podStatuses
		}

		// 对过滤后的数据进行分页
		paged := Paginate(filteredPods, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(filteredPods), // 使用过滤后的总数
			Limit:  limit,
			Offset: offset,
		})
	}
}

// getPodDetail 获取Pod详情的处理函数
// @Summary 获取 Pod 详情
// @Description 获取指定命名空间下的Pod详情
// @Tags Pod
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "Pod 名称"
// @Success 200 {object} model.APIResponse
// @Router /pods/{namespace}/{name} [get]
func getPodDetail(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		ns := c.Param("namespace")
		name := c.Param("name")
		pod, err := clientset.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}
		containers := make([]string, 0, len(pod.Spec.Containers))
		for _, ctn := range pod.Spec.Containers {
			containers = append(containers, ctn.Name + " (" + ctn.Image + ")")
		}
		podDetail := model.PodDetail{
			Namespace:   pod.Namespace,
			Name:        pod.Name,
			Status:      string(pod.Status.Phase),
			PodIP:       pod.Status.PodIP,
			NodeName:    pod.Spec.NodeName,
			StartTime:   pod.Status.StartTime.Format("2006-01-02 15:04:05"),
			Labels:      pod.Labels,
			Annotations: pod.Annotations,
			Containers:  containers,
		}
		middleware.ResponseSuccess(c, podDetail, "success", nil)
	}
}

// filterPodsBySearch 根据搜索关键词过滤Pod
func filterPodsBySearch(pods []model.PodStatus, search string) []model.PodStatus {
	if search == "" {
		return pods
	}

	searchLower := strings.ToLower(search)
	var filtered []model.PodStatus

	for _, pod := range pods {
		// 检查Pod的各个字段是否匹配搜索关键词
		if strings.Contains(strings.ToLower(pod.Name), searchLower) ||
			strings.Contains(strings.ToLower(pod.Namespace), searchLower) ||
			strings.Contains(strings.ToLower(pod.Status), searchLower) ||
			strings.Contains(strings.ToLower(pod.PodIP), searchLower) ||
			strings.Contains(strings.ToLower(pod.NodeName), searchLower) {
			filtered = append(filtered, pod)
		}
	}

	return filtered
}
