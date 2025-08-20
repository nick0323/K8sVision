package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterStorageClass 注册 StorageClass 相关路由
func RegisterStorageClass(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listStorageClasses func(context.Context, *kubernetes.Clientset) ([]model.StorageClassStatus, error),
) {
	r.GET("/storageclasses", getStorageClassList(logger, getK8sClient, listStorageClasses))
	r.GET("/storageclasses/:name", getStorageClassDetail(logger, getK8sClient))
}

// @Summary 获取 StorageClass 列表
// @Description 支持分页
// @Tags StorageClass
// @Security BearerAuth
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /storageclasses [get]
func getStorageClassList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listStorageClasses func(context.Context, *kubernetes.Clientset) ([]model.StorageClassStatus, error),
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

		storageClassStatuses, err := listStorageClasses(ctx, clientset)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		paged := Paginate(storageClassStatuses, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(storageClassStatuses),
			Limit:  limit,
			Offset: offset,
		})
	}
}

// @Summary 获取 StorageClass 详情
// @Description 获取指定 StorageClass 详情
// @Tags StorageClass
// @Security BearerAuth
// @Param name path string true "StorageClass 名称"
// @Success 200 {object} model.APIResponse
// @Router /storageclasses/{name} [get]
func getStorageClassDetail(
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
		storageClass, err := clientset.StorageV1().StorageClasses().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		isDefault := false
		if storageClass.Annotations != nil {
			if _, ok := storageClass.Annotations["storageclass.kubernetes.io/is-default-class"]; ok {
				isDefault = true
			}
		}

		storageClassDetail := model.StorageClassDetail{
			Name:              storageClass.Name,
			Provisioner:       storageClass.Provisioner,
			ReclaimPolicy:     string(*storageClass.ReclaimPolicy),
			VolumeBindingMode: string(*storageClass.VolumeBindingMode),
			IsDefault:         isDefault,
			Labels:            storageClass.Labels,
			Annotations:       storageClass.Annotations,
			Parameters:        storageClass.Parameters,
		}
		middleware.ResponseSuccess(c, storageClassDetail, "success", nil)
	}
}
