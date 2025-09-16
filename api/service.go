package api

import (
	"context"
	"fmt"

	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func RegisterService(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listServices func(context.Context, *kubernetes.Clientset, string) ([]model.ServiceStatus, error),
) {
	r.GET("/services", getServiceList(logger, getK8sClient, listServices))
	r.GET("/services/:namespace/:name", getServiceDetail(logger, getK8sClient))
}

func getServiceList(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listServices func(context.Context, *kubernetes.Clientset, string) ([]model.ServiceStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		HandleListWithPagination(c, logger, func(ctx context.Context, params PaginationParams) ([]model.ServiceStatus, error) {
			clientset, _, err := getK8sClient()
			if err != nil {
				return nil, err
			}
			return listServices(ctx, clientset, params.Namespace)
		}, ListSuccessMessage)
	}
}

func getServiceDetail(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		HandleDetailWithK8s(c, logger, getK8sClient,
			func(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) (model.ServiceDetail, error) {
				svc, err := clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
				if err != nil {
					return model.ServiceDetail{}, err
				}

				ports := make([]string, 0, len(svc.Spec.Ports))
				for _, p := range svc.Spec.Ports {
					ports = append(ports, fmt.Sprintf("%d/%s", p.Port, p.Protocol))
				}

				return model.ServiceDetail{
					CommonResourceFields: model.CommonResourceFields{
						Namespace: svc.Namespace,
						Name:      svc.Name,
						Status:    "Active",
						BaseMetadata: model.BaseMetadata{
							Labels:      svc.Labels,
							Annotations: svc.Annotations,
						},
					},
					Type:      string(svc.Spec.Type),
					ClusterIP: svc.Spec.ClusterIP,
					Ports:     ports,
					Selector:  svc.Spec.Selector,
				}, nil
			}, DetailSuccessMessage)
	}
}
