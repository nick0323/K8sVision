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

// RegisterPV 注册 PV 相关路由
func RegisterPV(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listPVs func(context.Context, *kubernetes.Clientset) ([]model.PVStatus, error),
) {
	r.GET("/pvs", getPVList(logger, getK8sClient, listPVs))
	r.GET("/pvs/:name", getPVDetail(logger, getK8sClient))
}

// getPVList 获取PV列表的处理函数
// @Summary 获取 PV 列表
// @Description 获取PV列表，支持分页
// @Tags PV
// @Security BearerAuth
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /pvs [get]
func getPVList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listPVs func(context.Context, *kubernetes.Clientset) ([]model.PVStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		search := c.DefaultQuery("search", "") // 新增：搜索关键词

		pvStatuses, err := listPVs(ctx, clientset)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 新增：如果提供了搜索关键词，先进行搜索过滤
		var filteredPVs []model.PVStatus
		if search != "" {
			filteredPVs = filterPVsBySearch(pvStatuses, search)
		} else {
			filteredPVs = pvStatuses
		}

		// 对过滤后的数据进行分页
		paged := Paginate(filteredPVs, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(filteredPVs), // 使用过滤后的总数
			Limit:  limit,
			Offset: offset,
		})
	}
}

// getPVDetail 获取PV详情的处理函数
// @Summary 获取 PV 详情
// @Description 获取指定PV的详细信息
// @Tags PV
// @Security BearerAuth
// @Param name path string true "PV 名称"
// @Success 200 {object} model.APIResponse
// @Router /pvs/{name} [get]
func getPVDetail(
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
		name := c.Param("name")
		pv, err := clientset.CoreV1().PersistentVolumes().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		capacity := ""
		if pv.Spec.Capacity != nil {
			if storage, ok := pv.Spec.Capacity[v1.ResourceStorage]; ok {
				capacity = storage.String()
			}
		}

		accessModes := make([]string, 0)
		for _, mode := range pv.Spec.AccessModes {
			accessModes = append(accessModes, string(mode))
		}

		storageClass := ""
		if pv.Spec.StorageClassName != "" {
			storageClass = pv.Spec.StorageClassName
		}

		claimRef := ""
		if pv.Spec.ClaimRef != nil {
			claimRef = pv.Spec.ClaimRef.Namespace + "/" + pv.Spec.ClaimRef.Name
		}

		pvDetail := model.PVDetail{
			Name:          pv.Name,
			Status:        string(pv.Status.Phase),
			Capacity:      capacity,
			AccessMode:    accessModes,
			StorageClass:  storageClass,
			ClaimRef:      claimRef,
			ReclaimPolicy: string(pv.Spec.PersistentVolumeReclaimPolicy),
			Labels:        pv.Labels,
			Annotations:   pv.Annotations,
		}
		middleware.ResponseSuccess(c, pvDetail, "success", nil)
	}
}

// filterPVsBySearch 根据搜索关键词过滤PV
func filterPVsBySearch(pvs []model.PVStatus, search string) []model.PVStatus {
	if search == "" {
		return pvs
	}

	searchLower := strings.ToLower(search)
	var filtered []model.PVStatus

	for _, pv := range pvs {
		// 检查PV的各个字段是否匹配搜索关键词
		if strings.Contains(strings.ToLower(pv.Name), searchLower) ||
			strings.Contains(strings.ToLower(pv.Status), searchLower) ||
			strings.Contains(strings.ToLower(pv.Capacity), searchLower) ||
			strings.Contains(strings.ToLower(pv.StorageClass), searchLower) {
			filtered = append(filtered, pv)
		}
	}

	return filtered
}
