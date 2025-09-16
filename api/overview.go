package api

import (
	"net/http"

	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterOverview(
	r *gin.RouterGroup,
	logger *zap.Logger,
	overviewFunc func(limit, offset int) (*model.OverviewStatus, string, error),
) {
	r.GET("/overview", func(c *gin.Context) {
		params := ParsePaginationParams(c)
		overview, msg, err := overviewFunc(params.Limit, params.Offset)
		if err != nil {
			middleware.ResponseError(c, logger, err, http.StatusInternalServerError)
			return
		}
		middleware.ResponseSuccess(c, overview, msg, &model.PageMeta{
			Total:  1,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	})
}
