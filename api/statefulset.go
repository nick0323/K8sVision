package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"
	"github.com/nick0323/K8sVision/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterStatefulSet 注册 StatefulSet 相关路由
func RegisterStatefulSet(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listStatefulSets func(context.Context, *kubernetes.Clientset, string) ([]model.StatefulSetStatus, error),
) {
	r.GET("/statefulsets", getStatefulSetList(logger, getK8sClient, listStatefulSets))
	r.GET("/statefulsets/:namespace/:name", getStatefulSetDetail(logger, getK8sClient))
}

// getStatefulSetList 获取StatefulSet列表的处理函数
// @Summary 获取 StatefulSet 列表
// @Description 获取StatefulSet列表，支持分页和搜索
// @Tags StatefulSet
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Param search query string false "搜索关键词（支持名称、命名空间、状态等字段搜索）"
// @Success 200 {object} model.APIResponse
// @Router /statefulsets [get]
func getStatefulSetList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listStatefulSets func(context.Context, *kubernetes.Clientset, string) ([]model.StatefulSetStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		namespace := c.DefaultQuery("namespace", "")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		search := c.DefaultQuery("search", "") // 新增：搜索关键词

		statefulSets, err := listStatefulSets(ctx, clientset, namespace)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 新增：如果提供了搜索关键词，先进行搜索过滤
		var filteredStatefulSets []model.StatefulSetStatus
		if search != "" {
			filteredStatefulSets = filterStatefulSetsBySearch(statefulSets, search)
		} else {
			filteredStatefulSets = statefulSets
		}

		// 对过滤后的数据进行分页
		paged := Paginate(filteredStatefulSets, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(filteredStatefulSets), // 使用过滤后的总数
			Limit:  limit,
			Offset: offset,
		})
	}
}

// filterStatefulSetsBySearch 根据搜索关键词过滤StatefulSet
func filterStatefulSetsBySearch(statefulSets []model.StatefulSetStatus, search string) []model.StatefulSetStatus {
	if search == "" {
		return statefulSets
	}
	searchLower := strings.ToLower(search)
	var filtered []model.StatefulSetStatus
	for _, statefulSet := range statefulSets {
		if strings.Contains(strings.ToLower(statefulSet.Name), searchLower) ||
			strings.Contains(strings.ToLower(statefulSet.Namespace), searchLower) ||
			strings.Contains(strings.ToLower(statefulSet.Status), searchLower) ||
			strings.Contains(strings.ToLower(strconv.Itoa(int(statefulSet.Available))), searchLower) ||
			strings.Contains(strings.ToLower(strconv.Itoa(int(statefulSet.Desired))), searchLower) {
			filtered = append(filtered, statefulSet)
		}
	}
	return filtered
}

// getStatefulSetDetail 获取StatefulSet详情的处理函数
// @Summary 获取 StatefulSet 详情
// @Description 获取指定命名空间下的StatefulSet详情
// @Tags StatefulSet
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "StatefulSet 名称"
// @Success 200 {object} model.APIResponse
// @Router /statefulsets/{namespace}/{name} [get]
func getStatefulSetDetail(
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
		sts, err := clientset.AppsV1().StatefulSets(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		image := ""
		if len(sts.Spec.Template.Spec.Containers) > 0 {
			image = sts.Spec.Template.Spec.Containers[0].Image
		}

		statefulSetDetail := model.StatefulSetDetail{
			Namespace:   sts.Namespace,
			Name:        sts.Name,
			Replicas:    *sts.Spec.Replicas,
			Available:   sts.Status.AvailableReplicas,
			Desired:     *sts.Spec.Replicas,
			Status:      service.GetWorkloadStatus(sts.Status.AvailableReplicas, *sts.Spec.Replicas),
			Labels:      sts.Labels,
			Annotations: sts.Annotations,
			Selector:    sts.Spec.Selector.MatchLabels,
			ServiceName: sts.Spec.ServiceName,
			Image:       image,
		}
		middleware.ResponseSuccess(c, statefulSetDetail, "success", nil)
	}
}
