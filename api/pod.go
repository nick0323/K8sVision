package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterPod 注册 Pod 相关路由，包括列表和详情接口
// @Summary 获取 Pod 列表
// @Description 支持分页
// @Tags Pod
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /pods [get]
//
// @Summary 获取 Pod 详情
// @Description 获取指定命名空间下的 Pod 详情
// @Tags Pod
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "Pod 名称"
// @Success 200 {object} model.APIResponse
// @Router /pods/{namespace}/{name} [get]
func RegisterPod(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listPodsWithRaw func(context.Context, *kubernetes.Clientset, model.PodMetricsMap, string) ([]model.PodStatus, *v1.PodList, error),
) {
	r.GET("/pods", func(c *gin.Context) {
		clientset, metricsClient, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		namespace := c.DefaultQuery("namespace", "")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
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
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		paged := Paginate(podStatuses, offset, limit)
		ResponseOK(c, paged, "success", &model.PageMeta{
			Total:  len(podStatuses),
			Limit:  limit,
			Offset: offset,
		})
	})

	r.GET("/pods/:namespace/:name", func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		ns := c.Param("namespace")
		name := c.Param("name")
		pod, err := clientset.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			ResponseError(c, logger, err, http.StatusNotFound)
			return
		}
		containers := make([]string, 0, len(pod.Spec.Containers))
		for _, ctn := range pod.Spec.Containers {
			containers = append(containers, ctn.Name)
		}
		podDetail := model.PodDetail{
			Namespace:   pod.Namespace,
			Name:        pod.Name,
			Status:      string(pod.Status.Phase),
			PodIP:       pod.Status.PodIP,
			NodeName:    pod.Spec.NodeName,
			StartTime:   model.FormatTime(pod.Status.StartTime),
			Labels:      pod.Labels,
			Annotations: pod.Annotations,
			Containers:  containers,
		}
		ResponseOK(c, podDetail, "success", nil)
	})
}
