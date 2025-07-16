package api

import (
	"net/http"
	"strconv"

	"github.com/nick0323/K8sVision/backend/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterOverview 注册 overview 接口，overviewFunc 只返回 (*OverviewStatus, string, error)
// @Summary 获取集群资源总览
// @Description 获取集群整体资源状态
// @Tags Overview
// @Security BearerAuth
// @Param namespace query string false "命名空间"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} model.APIResponse
// @Router /overview [get]
func RegisterOverview(
	r *gin.RouterGroup,
	logger *zap.Logger,
	overviewFunc func(namespace string, limit, offset int) (*model.OverviewStatus, string, error),
) {
	r.GET("/overview", func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", "")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		overview, msg, err := overviewFunc(namespace, limit, offset)
		if err != nil {
			ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		ResponseOK(c, overview, msg, &model.PageMeta{
			Total:  1,
			Limit:  limit,
			Offset: offset,
		})
	})
}
