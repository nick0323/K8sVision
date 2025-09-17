package api

import (
	"context"
	"net/http"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RegisterSecret 注册 Secret 相关路由
func RegisterSecret(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient K8sClientProvider, // 使用类型别名简化签名
	listSecrets func(context.Context, *kubernetes.Clientset, string) ([]model.SecretStatus, error),
) {
	r.GET("/secrets", getSecretList(logger, getK8sClient, listSecrets))
	r.GET("/secrets/:namespace/:name", getSecretDetail(logger, getK8sClient))
}

// getSecretList 获取Secret列表的处理函数
func getSecretList(
	logger *zap.Logger,
	getK8sClient K8sClientProvider,
	listSecrets func(context.Context, *kubernetes.Clientset, string) ([]model.SecretStatus, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		HandleListWithPagination(c, logger, func(ctx context.Context, params PaginationParams) ([]model.SecretStatus, error) {
			clientset, _, err := getK8sClient()
			if err != nil {
				return nil, err
			}
			return listSecrets(ctx, clientset, params.Namespace)
		}, ListSuccessMessage)
	}
}

// getSecretDetail 获取Secret详情的处理函数
func getSecretDetail(
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
		secret, err := clientset.CoreV1().Secrets(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		keys := make([]string, 0, len(secret.Data))
		for key := range secret.Data {
			keys = append(keys, key)
		}

		// 将base64编码的数据转换为字符串
		data := make(map[string]string)
		for key, value := range secret.Data {
			data[key] = string(value)
		}

		secretDetail := model.SecretDetail{
			CommonResourceFields: model.CommonResourceFields{
				Namespace: secret.Namespace,
				Name:      secret.Name,
				Status:    "Active",
				BaseMetadata: model.BaseMetadata{
					Labels:      secret.Labels,
					Annotations: secret.Annotations,
				},
			},
			Type:      string(secret.Type),
			DataCount: len(secret.Data),
			Keys:      keys,
			Data:      data,
		}
		middleware.ResponseSuccess(c, secretDetail, "success", nil)
	}
}
