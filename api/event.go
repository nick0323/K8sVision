package api

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// RegisterEvent 注册 Event 相关路由
func RegisterEvent(
	r *gin.RouterGroup,
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listEvents func(context.Context, *kubernetes.Clientset, string) ([]model.EventStatus, error),
) {
	r.GET("/events", getEventList(logger, getK8sClient, listEvents))
	r.GET("/events/:namespace/:name", getEventDetail(logger, getK8sClient))
}

// getEventList 获取Event列表的处理函数
// @Summary 获取 Event 列表
// @Description 获取Event列表，支持分页和搜索
// @Tags Event
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Param search query string false "搜索关键词（支持名称、命名空间、原因等字段搜索）"
// @Success 200 {object} model.APIResponse
// @Router /events [get]
func getEventList(
	logger *zap.Logger,
	getK8sClient func() (*kubernetes.Clientset, *versioned.Clientset, error),
	listEvents func(context.Context, *kubernetes.Clientset, string) ([]model.EventStatus, error),
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

		events, err := listEvents(ctx, clientset, namespace)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}

		// 新增：如果提供了搜索关键词，先进行搜索过滤
		var filteredEvents []model.EventStatus
		if search != "" {
			filteredEvents = filterEventsBySearch(events, search)
		} else {
			filteredEvents = events
		}

		// 按LastSeen时间倒序排列（最新事件在前）
		sortEventsByLastSeen(filteredEvents)

		// 对过滤后的数据进行分页
		paged := Paginate(filteredEvents, offset, limit)
		middleware.ResponseSuccess(c, paged, "success", &model.PageMeta{
			Total:  len(filteredEvents), // 使用过滤后的总数
			Limit:  limit,
			Offset: offset,
		})
	}
}

// filterEventsBySearch 根据搜索关键词过滤Event
func filterEventsBySearch(events []model.EventStatus, search string) []model.EventStatus {
	if search == "" {
		return events
	}
	searchLower := strings.ToLower(search)
	var filtered []model.EventStatus
	for _, event := range events {
		if strings.Contains(strings.ToLower(event.Name), searchLower) ||
			strings.Contains(strings.ToLower(event.Namespace), searchLower) ||
			strings.Contains(strings.ToLower(event.Reason), searchLower) ||
			strings.Contains(strings.ToLower(event.Message), searchLower) ||
			strings.Contains(strings.ToLower(event.Type), searchLower) ||
			strings.Contains(strings.ToLower(event.FirstSeen), searchLower) ||
			strings.Contains(strings.ToLower(event.LastSeen), searchLower) ||
			strings.Contains(strings.ToLower(event.Duration), searchLower) ||
			strings.Contains(strings.ToLower(strconv.Itoa(int(event.Count))), searchLower) {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// getEventDetail 获取Event详情的处理函数
// @Summary 获取 Event 详情
// @Description 获取指定命名空间下的Event详情
// @Tags Event
// @Security BearerAuth
// @Param namespace path string true "命名空间"
// @Param name path string true "Event 名称"
// @Success 200 {object} model.APIResponse
// @Router /events/{namespace}/{name} [get]
func getEventDetail(
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
		event, err := clientset.CoreV1().Events(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusNotFound)
			return
		}

		eventDetail := model.EventDetail{
			Namespace:   event.Namespace,
			Name:        event.Name,
			Reason:      event.Reason,
			Message:     event.Message,
			Type:        event.Type,
			Count:       event.Count,
			FirstSeen:   event.FirstTimestamp.Format("2006-01-02 15:04:05"),
			LastSeen:    event.LastTimestamp.Format("2006-01-02 15:04:05"),
			Duration:    event.LastTimestamp.Sub(event.FirstTimestamp.Time).String(),
			Labels:      event.Labels,
			Annotations: event.Annotations,
		}
		middleware.ResponseSuccess(c, eventDetail, "success", nil)
	}
}

// sortEventsByLastSeen 按LastSeen时间倒序排列Events（最新事件在前）
func sortEventsByLastSeen(events []model.EventStatus) {
	sort.Slice(events, func(i, j int) bool {
		// 解析时间字符串进行比较
		timeI, errI := time.Parse("2006-01-02 15:04:05", events[i].LastSeen)
		timeJ, errJ := time.Parse("2006-01-02 15:04:05", events[j].LastSeen)
		
		// 如果解析失败，保持原始顺序
		if errI != nil || errJ != nil {
			return false
		}
		
		// 倒序排列：最新的时间在前
		return timeI.After(timeJ)
	})
}
