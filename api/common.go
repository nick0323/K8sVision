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
	"strings"
	"reflect"
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
		search := c.DefaultQuery("search", "") // 新增：搜索关键词

		items, err := listFunc(ctx, clientset, namespace)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 新增：如果提供了搜索关键词，先进行搜索过滤
		var filteredItems []T
		if search != "" {
			filteredItems = filterItemsBySearch(items, search)
		} else {
			filteredItems = items
		}

		// 对过滤后的数据进行分页
		paged := Paginate(filteredItems, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(filteredItems), // 使用过滤后的总数
			Limit:  limit,
			Offset: offset,
		})
	}
}

// filterItemsBySearch 根据搜索关键词过滤项目
// 这是一个通用的搜索过滤函数，支持常见的字段搜索
func filterItemsBySearch[T any](items []T, search string) []T {
	if search == "" {
		return items
	}

	searchLower := strings.ToLower(search)
	var filtered []T

	for _, item := range items {
		// 使用反射获取结构体字段值进行搜索
		if matchesSearch(item, searchLower) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// matchesSearch 检查单个项目是否匹配搜索关键词
func matchesSearch[T any](item T, searchLower string) bool {
	val := reflect.ValueOf(item)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return false
	}

	// 定义要搜索的字段名（这些字段通常是字符串类型且对用户有意义）
	searchableFields := []string{"Name", "Namespace", "Status", "PodIP", "NodeName", "Image"}
	
	for _, fieldName := range searchableFields {
		field := val.FieldByName(fieldName)
		if field.IsValid() && field.CanInterface() {
			fieldValue := field.Interface()
			if str, ok := fieldValue.(string); ok {
				if strings.Contains(strings.ToLower(str), searchLower) {
					return true
				}
			}
		}
	}

	return false
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
