package service

import (
	"context"
	"strings"

	"github.com/nick0323/K8sVision/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListPVCs 获取 PVC 列表
func ListPVCs(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]model.PVCStatus, error) {
	var pvcList *v1.PersistentVolumeClaimList
	var err error

	if namespace == "" {
		pvcList, err = clientset.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	} else {
		pvcList, err = clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	}

	if err != nil {
		return nil, err
	}

	var pvcStatuses []model.PVCStatus
	for _, pvc := range pvcList.Items {
		status := "Pending"
		if pvc.Status.Phase == v1.ClaimBound {
			status = "Bound"
		} else if pvc.Status.Phase == v1.ClaimLost {
			status = "Lost"
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

		volumeName := ""
		if pvc.Spec.VolumeName != "" {
			volumeName = pvc.Spec.VolumeName
		}

		pvcStatus := model.PVCStatus{
			Namespace:    pvc.Namespace,
			Name:         pvc.Name,
			Status:       status,
			Capacity:     capacity,
			AccessMode:   strings.Join(accessModes, ","),
			StorageClass: storageClass,
			VolumeName:   volumeName,
		}
		pvcStatuses = append(pvcStatuses, pvcStatus)
	}

	return pvcStatuses, nil
}
