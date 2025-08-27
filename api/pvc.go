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

// RegisterPVC 注册 PVC 相关路由
func RegisterPVC(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listPVCs func(context.Context, *kubernetes.Clientset, string) ([]model.PVCStatus, error),
) {
	r.GET("/pvcs", getPVCList(logger, getK8sClient, listPVCs))
	r.GET("/pvcs/:namespace/:name", getPVCDetail(logger, getK8sClient))
}

// getPVCList 获取PVC列表的处理函数
// @Summary 获取 PVC 列表
// @Description 获取PVC列表，支持分页
// @Tags PVC
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /pvcs [get]
func getPVCList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listPVCs func(context.Context, *kubernetes.Clientset, string) ([]model.PVCStatus, error),
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

		pvcStatuses, err := listPVCs(ctx, clientset, namespace)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 新增：如果提供了搜索关键词，先进行搜索过滤
		var filteredPVCs []model.PVCStatus
		if search != "" {
			filteredPVCs = filterPVCsBySearch(pvcStatuses, search)
		} else {
			filteredPVCs = pvcStatuses
		}

		// 对过滤后的数据进行分页
		paged := Paginate(filteredPVCs, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(filteredPVCs), // 使用过滤后的总数
			Limit:  limit,
			Offset: offset,
		})
	}
}

// getPVCDetail 获取PVC详情的处理函数
// @Summary 获取 PVC 详情
// @Description 获取指定命名空间下的PVC详情
// @Tags PVC
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "PVC 名称"
// @Success 200 {object} model.APIResponse
// @Router /pvcs/{namespace}/{name} [get]
func getPVCDetail(
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
		pvc, err := clientset.CoreV1().PersistentVolumeClaims(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		capacity := ""
		if pvc.Status.Capacity != nil {
			if storage, ok := pvc.Status.Capacity[v1.ResourceStorage]; ok {
				capacity = storage.String()
			}
		}

		accessModes := make([]string, 0)
		for _, mode := range pvc.Spec.AccessModes {
			accessModes = append(accessModes, string(mode))
		}

		storageClass := ""
		if pvc.Spec.StorageClassName != nil {
			storageClass = *pvc.Spec.StorageClassName
		}

		pvcDetail := model.PVCDetail{
			Namespace:    pvc.Namespace,
			Name:         pvc.Name,
			Status:       string(pvc.Status.Phase),
			Capacity:     capacity,
			AccessMode:   accessModes,
			StorageClass: storageClass,
			VolumeName:   pvc.Spec.VolumeName,
			Labels:       pvc.Labels,
			Annotations:  pvc.Annotations,
		}
		middleware.ResponseSuccess(c, pvcDetail, "success", nil)
	}
}

// filterPVCsBySearch 根据搜索关键词过滤PVC
func filterPVCsBySearch(pvcs []model.PVCStatus, search string) []model.PVCStatus {
	if search == "" {
		return pvcs
	}

	searchLower := strings.ToLower(search)
	var filtered []model.PVCStatus

	for _, pvc := range pvcs {
		// 检查PVC的各个字段是否匹配搜索关键词
		if strings.Contains(strings.ToLower(pvc.Name), searchLower) ||
			strings.Contains(strings.ToLower(pvc.Namespace), searchLower) ||
			strings.Contains(strings.ToLower(pvc.Status), searchLower) ||
			strings.Contains(strings.ToLower(pvc.Capacity), searchLower) ||
			strings.Contains(strings.ToLower(pvc.StorageClass), searchLower) {
			filtered = append(filtered, pvc)
		}
	}

	return filtered
}
