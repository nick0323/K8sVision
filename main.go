// @title K8sVision API
// @version 1.0
// @description K8sVision 后端 API 文档
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"github.com/nick0323/K8sVision/api"
	"github.com/nick0323/K8sVision/service"

	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	// 注册 swagger 文档
	"os"

	_ "github.com/nick0323/K8sVision/docs"
)

var logger *zap.Logger

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

// main 入口，负责依赖注入和路由注册
func main() {
	initConfig()
	api.InitJWTSecret()
	logger, _ = zap.NewProduction()
	defer logger.Sync()

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		traceId := uuid.New().String()
		c.Set("traceId", traceId)
		c.Next()
	})
	r.POST("/api/login", api.LoginHandler)

	// 通过环境变量控制 swagger 路由注册
	if os.Getenv("SWAGGER_ENABLE") == "true" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "ok")
	})
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}
	apiGroup := r.Group("/api")
	apiGroup.Use(api.JWTAuthMiddleware())
	api.RegisterOverview(apiGroup, logger, func(namespace string, limit, offset int) (*model.OverviewStatus, string, error) {
		clientset, _, err := service.GetK8sClient()
		if err != nil {
			return nil, "k8s client error", err
		}
		overview, err := service.GetOverviewStatus(clientset)
		if err != nil {
			return nil, "overview error", err
		}
		return overview, "", nil
	})
	api.RegisterNode(apiGroup, logger, service.GetK8sClient, service.ListPodsWithRaw, service.ListNodes)
	api.RegisterPod(apiGroup, logger, service.GetK8sClient, service.ListPodsWithRaw)
	api.RegisterDeployment(apiGroup, logger, service.GetK8sClient, service.ListDeployments)
	api.RegisterStatefulSet(apiGroup, logger, service.GetK8sClient, service.ListStatefulSets)
	api.RegisterDaemonSet(apiGroup, logger, service.GetK8sClient, service.ListDaemonSets)
	api.RegisterService(apiGroup, logger, service.GetK8sClient, service.ListServices)
	api.RegisterNamespace(apiGroup, logger, service.GetK8sClient, service.ListNamespaces)
	api.RegisterEvent(apiGroup, logger, service.GetK8sClient, service.ListEvents)
	api.RegisterCronJob(apiGroup, logger, service.GetK8sClient, service.ListCronJobs)
	api.RegisterJob(apiGroup, logger, service.GetK8sClient, service.ListJobs)
	api.RegisterIngress(apiGroup, logger, service.GetK8sClient, service.ListIngresses)
	r.Run(":" + port)
}
