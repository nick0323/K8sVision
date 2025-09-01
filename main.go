// Package main K8sVision 主程序
// @title K8sVision API
// @version 1.0
// @description K8sVision 后端 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入 "Bearer " 加上JWT token，例如: "Bearer abcde12345"

package main

import (
	"flag"
	"os"
	"time"

	"github.com/nick0323/K8sVision/api"
	"github.com/nick0323/K8sVision/api/middleware"
	"github.com/nick0323/K8sVision/cache"
	"github.com/nick0323/K8sVision/config"
	"github.com/nick0323/K8sVision/monitor"
	"github.com/nick0323/K8sVision/service"

	"github.com/nick0323/K8sVision/model"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "github.com/nick0323/K8sVision/docs"
)

var (
	configFile = flag.String("config", "", "配置文件路径")
	logger     *zap.Logger
	configMgr  *config.Manager
	cacheMgr   *cache.Manager
	monitorMgr *monitor.Monitor
)

func initLogger(cfg *model.Config) (*zap.Logger, error) {
	var zapConfig zap.Config
	
	if cfg.IsDevelopment() {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}
	
	switch cfg.Log.Level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	
	if cfg.Log.Format == "console" {
		zapConfig.Encoding = "console"
	} else {
		zapConfig.Encoding = "json"
	}
	
	return zapConfig.Build()
}

func main() {
	flag.Parse()
	
	tempLogger, _ := zap.NewProduction()
	
	configMgr = config.NewManager(tempLogger)
	
	if err := configMgr.Load(*configFile); err != nil {
		tempLogger.Fatal("加载配置失败", zap.Error(err))
	}
	
	cfg := configMgr.GetConfig()
	
	var err error
	logger, err = initLogger(cfg)
	if err != nil {
		tempLogger.Fatal("初始化日志失败", zap.Error(err))
	}
	defer logger.Sync()
	
	configMgr = config.NewManager(logger)
	if err := configMgr.Load(*configFile); err != nil {
		logger.Fatal("重新加载配置失败", zap.Error(err))
	}
	
	cacheMgr = cache.NewManager(&cfg.Cache, logger)
	defer cacheMgr.Close()
	
	monitorMgr = monitor.NewMonitor(logger)
	defer monitorMgr.Close()
	
	service.SetConfigManager(configMgr)
	service.SetCacheManager(cacheMgr)
	
	api.InitJWTSecret()
	
	if err := configMgr.Watch(); err != nil {
		logger.Warn("启动配置监听失败", zap.Error(err))
	}
	defer configMgr.Close()

	if cfg.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	
	r.Use(middleware.ErrorHandler(logger))
	r.Use(middleware.TraceMiddleware())
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.MetricsMiddleware(monitorMgr.GetMetrics()))
	
	if cfg.Auth.EnableRateLimit {
		r.Use(middleware.ConcurrencyMiddleware(cfg.Auth.RateLimit))
	}
	
	r.POST("/api/login", api.LoginHandler)

	if os.Getenv("SWAGGER_ENABLE") == "true" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "ok")
	})
	
	r.GET("/metrics", func(c *gin.Context) {
		stats := monitorMgr.GetMetrics().GetStats()
		c.JSON(200, stats)
	})
	
	r.GET("/cache/stats", func(c *gin.Context) {
		stats := cacheMgr.GetAllStats()
		c.JSON(200, stats)
	})
	
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.JWTAuthMiddleware(logger))
	
	if cfg.Cache.Enabled {
		apiGroup.Use(middleware.CacheMiddleware(cacheMgr, cfg.Cache.TTL))
	}

	api.RegisterOverview(apiGroup, logger, func(limit, offset int) (*model.OverviewStatus, string, error) {
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
	api.RegisterIngress(apiGroup, logger, service.GetK8sClient, service.ListIngresses)

	api.RegisterCronJob(apiGroup, logger, service.GetK8sClient, service.ListCronJobs)
	api.RegisterJob(apiGroup, logger, service.GetK8sClient, service.ListJobs)

	api.RegisterNamespace(apiGroup, logger, service.GetK8sClient, service.ListNamespaces)
	api.RegisterEvent(apiGroup, logger, service.GetK8sClient, service.ListEvents)

	api.RegisterPVC(apiGroup, logger, service.GetK8sClient, service.ListPVCs)
	api.RegisterPV(apiGroup, logger, service.GetK8sClient, service.ListPVs)
	api.RegisterStorageClass(apiGroup, logger, service.GetK8sClient, service.ListStorageClasses)

	api.RegisterConfigMap(apiGroup, logger, service.GetK8sClient, service.ListConfigMaps)
	api.RegisterSecret(apiGroup, logger, service.GetK8sClient, service.ListSecrets)

	monitorMgr.StartPeriodicLogging(5 * time.Minute)

	serverAddr := cfg.GetServerAddress()
	logger.Info("服务器启动", 
		zap.String("address", serverAddr),
		zap.Bool("cacheEnabled", cfg.Cache.Enabled),
		zap.Bool("rateLimitEnabled", cfg.Auth.EnableRateLimit),
	)
	r.Run(serverAddr)
}
