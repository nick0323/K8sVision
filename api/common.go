package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/nick0323/K8sVision/model"
	"github.com/nick0323/K8sVision/api/middleware"

	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

var jwtSecret []byte

// InitJWTSecret 初始化 JWT 密钥，优先环境变量，其次配置管理器，最后默认
func InitJWTSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// 这里可以通过配置管理器获取，暂时保持兼容性
		secret = "k8svision-secret-key"
	}
	jwtSecret = []byte(secret)
	// 同时初始化中间件的JWT密钥
	middleware.InitJWTSecret(secret)
}

// GetTraceID 获取请求的追踪ID
func GetTraceID(c *gin.Context) string {
	tid := c.GetHeader("X-Trace-ID")
	if tid == "" {
		tid = c.GetString("traceId")
		if tid == "" {
			return ""
		}
	}
	return tid
}

// Paginate 对任意切片类型进行分页，返回 offset-limit 范围内的子切片
func Paginate[T any](list []T, offset, limit int) []T {
	if offset > len(list) {
		return []T{}
	}
	end := offset + limit
	if end > len(list) {
		end = len(list)
	}
	return list[offset:end]
}

// GenericListHandler 通用列表处理函数
func GenericListHandler[T any](
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listFunc func(context.Context, *kubernetes.Clientset, string) ([]T, error),
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

		items, err := listFunc(ctx, clientset, namespace)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		paged := Paginate(items, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(items),
			Limit:  limit,
			Offset: offset,
		})
	}
}

// GenericDetailHandler 通用详情处理函数
func GenericDetailHandler[T any](
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	getFunc func(context.Context, *kubernetes.Clientset, string, string) (T, error),
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

		item, err := getFunc(ctx, clientset, ns, name)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		middleware.ResponseSuccess(c, item, "success", nil)
	}
}
