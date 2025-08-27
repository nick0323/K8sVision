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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterConfigMap 注册 ConfigMap 相关路由
func RegisterConfigMap(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listConfigMaps func(context.Context, *kubernetes.Clientset, string) ([]model.ConfigMapStatus, error),
) {
	r.GET("/configmaps", getConfigMapList(logger, getK8sClient, listConfigMaps))
	r.GET("/configmaps/:namespace/:name", getConfigMapDetail(logger, getK8sClient))
}

// getConfigMapList 获取ConfigMap列表的处理函数
// @Summary 获取 ConfigMap 列表
// @Description 获取ConfigMap列表，支持分页和搜索
// @Tags ConfigMap
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Param search query string false "搜索关键词（支持名称、命名空间、数据等字段搜索）"
// @Success 200 {object} model.APIResponse
// @Router /configmaps [get]
func getConfigMapList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listConfigMaps func(context.Context, *kubernetes.Clientset, string) ([]model.ConfigMapStatus, error),
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

		configMaps, err := listConfigMaps(ctx, clientset, namespace)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 新增：如果提供了搜索关键词，先进行搜索过滤
		var filteredConfigMaps []model.ConfigMapStatus
		if search != "" {
			filteredConfigMaps = filterConfigMapsBySearch(configMaps, search)
		} else {
			filteredConfigMaps = configMaps
		}

		// 对过滤后的数据进行分页
		paged := Paginate(filteredConfigMaps, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(filteredConfigMaps), // 使用过滤后的总数
			Limit:  limit,
			Offset: offset,
		})
	}
}

// filterConfigMapsBySearch 根据搜索关键词过滤ConfigMap
func filterConfigMapsBySearch(configMaps []model.ConfigMapStatus, search string) []model.ConfigMapStatus {
	if search == "" {
		return configMaps
	}
	searchLower := strings.ToLower(search)
	var filtered []model.ConfigMapStatus
	for _, configMap := range configMaps {
		if strings.Contains(strings.ToLower(configMap.Name), searchLower) ||
			strings.Contains(strings.ToLower(configMap.Namespace), searchLower) ||
			strings.Contains(strings.ToLower(strconv.Itoa(configMap.DataCount)), searchLower) ||
			strings.Contains(strings.Join(configMap.Keys, ","), searchLower) {
			filtered = append(filtered, configMap)
		}
	}
	return filtered
}

// getConfigMapDetail 获取ConfigMap详情的处理函数
// @Summary 获取 ConfigMap 详情
// @Description 获取指定命名空间下的ConfigMap详情
// @Tags ConfigMap
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "ConfigMap 名称"
// @Success 200 {object} model.APIResponse
// @Router /configmaps/{namespace}/{name} [get]
func getConfigMapDetail(
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
		configMap, err := clientset.CoreV1().ConfigMaps(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		keys := make([]string, 0, len(configMap.Data))
		for key := range configMap.Data {
			keys = append(keys, key)
		}

		configMapDetail := model.ConfigMapDetail{
			Namespace:   configMap.Namespace,
			Name:        configMap.Name,
			DataCount:   len(configMap.Data),
			Keys:        keys,
			Labels:      configMap.Labels,
			Annotations: configMap.Annotations,
			Data:        configMap.Data,
		}
		middleware.ResponseSuccess(c, configMapDetail, "success", nil)
	}
}
