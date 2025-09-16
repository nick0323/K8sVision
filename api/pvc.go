package api

import (
	"context"
	"net/http"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func RegisterPVC(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listPVCs func(context.Context, *kubernetes.Clientset, string) ([]model.PVCStatus, error),
) {
	r.GET("/pvcs", getPVCList(logger, getK8sClient, listPVCs))
	r.GET("/pvcs/:namespace/:name", getPVCDetail(logger, getK8sClient))
}

func getPVCList(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listPVCs func(context.Context, *kubernetes.Clientset, string) ([]model.PVCStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		HandleListWithPagination(c, logger, func(ctx context.Context, params PaginationParams) ([]model.PVCStatus, error) {
			clientset, _, err := getK8sClient()
			if err != nil {
				return nil, err
			}
			return listPVCs(ctx, clientset, params.Namespace)
		}, ListSuccessMessage)
	}
}

func getPVCDetail(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientset, _, err := getK8sClient()
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ctx := GetRequestContext(c)
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
			CommonResourceFields: model.CommonResourceFields{
				Namespace: pvc.Namespace,
				Name:      pvc.Name,
				Status:    string(pvc.Status.Phase),
				BaseMetadata: model.BaseMetadata{
					Labels:      pvc.Labels,
					Annotations: pvc.Annotations,
				},
			},
			Capacity:     capacity,
			AccessMode:   accessModes,
			StorageClass: storageClass,
			VolumeName:   pvc.Spec.VolumeName,
		}
		middleware.ResponseSuccess(c, pvcDetail, DetailSuccessMessage, nil)
	}
}
