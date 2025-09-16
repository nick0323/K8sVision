package api

import (
	"context"

	"github.com/nick0323/K8sVision/model"
	"github.com/nick0323/K8sVision/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func RegisterDeployment(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listDeployments func(context.Context, *kubernetes.Clientset, string) ([]model.DeploymentStatus, error),
) {
	r.GET("/deployments", getDeploymentList(logger, getK8sClient, listDeployments))
	r.GET("/deployments/:namespace/:name", getDeploymentDetail(logger, getK8sClient))
}

func getDeploymentList(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listDeployments func(context.Context, *kubernetes.Clientset, string) ([]model.DeploymentStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		HandleListWithPagination(c, logger, func(ctx context.Context, params PaginationParams) ([]model.DeploymentStatus, error) {
			clientset, _, err := getK8sClient()
			if err != nil {
				return nil, err
			}
			return listDeployments(ctx, clientset, params.Namespace)
		}, ListSuccessMessage)
	}
}

func getDeploymentDetail(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		HandleDetailWithK8s(c, logger, getK8sClient,
			func(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) (model.DeploymentDetail, error) {
				dep, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
				if err != nil {
					return model.DeploymentDetail{}, err
				}

				image := ""
				if len(dep.Spec.Template.Spec.Containers) > 0 {
					image = dep.Spec.Template.Spec.Containers[0].Image
				}

				return model.DeploymentDetail{
					WorkloadCommonFields: model.WorkloadCommonFields{
						CommonResourceFields: model.CommonResourceFields{
							Namespace: dep.Namespace,
							Name:      dep.Name,
							Status:    service.GetWorkloadStatus(dep.Status.AvailableReplicas, *dep.Spec.Replicas),
							BaseMetadata: model.BaseMetadata{
								Labels:      dep.Labels,
								Annotations: dep.Annotations,
							},
						},
						Available: dep.Status.AvailableReplicas,
						Desired:   *dep.Spec.Replicas,
						Selector:  dep.Spec.Selector.MatchLabels,
						Image:     image,
					},
					Replicas: *dep.Spec.Replicas,
					Strategy: string(dep.Spec.Strategy.Type),
				}, nil
			}, DetailSuccessMessage)
	}
}
