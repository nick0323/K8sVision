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

// RegisterSecret 注册 Secret 相关路由
func RegisterSecret(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listSecrets func(context.Context, *kubernetes.Clientset, string) ([]model.SecretStatus, error),
) {
	r.GET("/secrets", getSecretList(logger, getK8sClient, listSecrets))
	r.GET("/secrets/:namespace/:name", getSecretDetail(logger, getK8sClient))
}

// getSecretList 获取Secret列表的处理函数
// @Summary 获取 Secret 列表
// @Description 获取Secret列表，支持分页和搜索
// @Tags Secret
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Param search query string false "搜索关键词（支持名称、命名空间、类型等字段搜索）"
// @Success 200 {object} model.APIResponse
// @Router /secrets [get]
func getSecretList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listSecrets func(context.Context, *kubernetes.Clientset, string) ([]model.SecretStatus, error),
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

		secrets, err := listSecrets(ctx, clientset, namespace)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 新增：如果提供了搜索关键词，先进行搜索过滤
		var filteredSecrets []model.SecretStatus
		if search != "" {
			filteredSecrets = filterSecretsBySearch(secrets, search)
		} else {
			filteredSecrets = secrets
		}

		// 对过滤后的数据进行分页
		paged := Paginate(filteredSecrets, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(filteredSecrets), // 使用过滤后的总数
			Limit:  limit,
			Offset: offset,
		})
	}
}

// filterSecretsBySearch 根据搜索关键词过滤Secret
func filterSecretsBySearch(secrets []model.SecretStatus, search string) []model.SecretStatus {
	if search == "" {
		return secrets
	}
	searchLower := strings.ToLower(search)
	var filtered []model.SecretStatus
	for _, secret := range secrets {
		if strings.Contains(strings.ToLower(secret.Name), searchLower) ||
			strings.Contains(strings.ToLower(secret.Namespace), searchLower) ||
			strings.Contains(strings.ToLower(secret.Type), searchLower) ||
			strings.Contains(strings.ToLower(strconv.Itoa(secret.DataCount)), searchLower) ||
			strings.Contains(strings.Join(secret.Keys, ","), searchLower) {
			filtered = append(filtered, secret)
		}
	}
	return filtered
}

// getSecretDetail 获取Secret详情的处理函数
// @Summary 获取 Secret 详情
// @Description 获取指定命名空间下的Secret详情
// @Tags Secret
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "Secret 名称"
// @Success 200 {object} model.APIResponse
// @Router /secrets/{namespace}/{name} [get]
func getSecretDetail(
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
			Namespace:   secret.Namespace,
			Name:        secret.Name,
			Type:        string(secret.Type),
			DataCount:   len(secret.Data),
			Keys:        keys,
			Labels:      secret.Labels,
			Annotations: secret.Annotations,
			Data:        data,
		}
		middleware.ResponseSuccess(c, secretDetail, "success", nil)
	}
}
