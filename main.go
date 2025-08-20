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

	// 注册 swagger 文档
	_ "github.com/nick0323/K8sVision/docs"
)

var (
	configFile = flag.String("config", "", "配置文件路径")
	logger     *zap.Logger
	configMgr  *config.Manager
	cacheMgr   *cache.Manager
	monitorMgr *monitor.Monitor
)

// initLogger 初始化日志
func initLogger(cfg *model.Config) (*zap.Logger, error) {
	var zapConfig zap.Config
	
	if cfg.IsDevelopment() {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}
	
	// 设置日志级别
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
	
	// 设置日志格式
	if cfg.Log.Format == "console" {
		zapConfig.Encoding = "console"
	} else {
		zapConfig.Encoding = "json"
	}
	
	return zapConfig.Build()
}

// main 入口，负责依赖注入和路由注册
func main() {
	flag.Parse()
	
	// 创建临时logger用于配置加载
	tempLogger, _ := zap.NewProduction()
	
	// 初始化配置管理器
	configMgr = config.NewManager(tempLogger)
	
	// 加载配置
	if err := configMgr.Load(*configFile); err != nil {
		tempLogger.Fatal("加载配置失败", zap.Error(err))
	}
	
	// 获取配置
	cfg := configMgr.GetConfig()
	
	// 初始化日志
	var err error
	logger, err = initLogger(cfg)
	if err != nil {
		tempLogger.Fatal("初始化日志失败", zap.Error(err))
	}
	defer logger.Sync()
	
	// 更新配置管理器的logger
	configMgr = config.NewManager(logger)
	if err := configMgr.Load(*configFile); err != nil {
		logger.Fatal("重新加载配置失败", zap.Error(err))
	}
	
	// 初始化缓存管理器
	cacheMgr = cache.NewManager(&cfg.Cache, logger)
	defer cacheMgr.Close()
	
	// 初始化性能监控器
	monitorMgr = monitor.NewMonitor(logger)
	defer monitorMgr.Close()
	
	// 设置配置管理器到service包
	service.SetConfigManager(configMgr)
	service.SetCacheManager(cacheMgr)
	
	// 初始化JWT密钥
	api.InitJWTSecret()
	
	// 启动配置监听
	if err := configMgr.Watch(); err != nil {
		logger.Warn("启动配置监听失败", zap.Error(err))
	}
	defer configMgr.Close()

	// 设置Gin模式
	if cfg.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New() // 使用gin.New()而不是gin.Default()，以便自定义中间件
	
	// 添加自定义中间件
	r.Use(middleware.ErrorHandler(logger))     // 错误处理中间件
	r.Use(middleware.TraceMiddleware())        // 请求追踪中间件
	r.Use(middleware.LoggingMiddleware(logger)) // 日志记录中间件
	r.Use(middleware.MetricsMiddleware(monitorMgr.GetMetrics())) // 性能监控中间件
	
	// 添加并发控制中间件（如果启用）
	if cfg.Auth.EnableRateLimit {
		r.Use(middleware.ConcurrencyMiddleware(cfg.Auth.RateLimit))
	}
	
	r.POST("/api/login", api.LoginHandler)

	// 通过环境变量控制 swagger 路由注册
	if os.Getenv("SWAGGER_ENABLE") == "true" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "ok")
	})
	
	// 添加性能指标API
	r.GET("/metrics", func(c *gin.Context) {
		stats := monitorMgr.GetMetrics().GetStats()
		c.JSON(200, stats)
	})
	
	// 添加缓存统计API
	r.GET("/cache/stats", func(c *gin.Context) {
		stats := cacheMgr.GetAllStats()
		c.JSON(200, stats)
	})
	
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.JWTAuthMiddleware(logger))
	
	// 为API组添加缓存中间件（如果启用）
	if cfg.Cache.Enabled {
		apiGroup.Use(middleware.CacheMiddleware(cacheMgr, cfg.Cache.TTL))
	}

	// Overview API
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

	// 计算资源API注册
	api.RegisterNode(apiGroup, logger, service.GetK8sClient, service.ListPodsWithRaw, service.ListNodes)
	api.RegisterPod(apiGroup, logger, service.GetK8sClient, service.ListPodsWithRaw)
	api.RegisterDeployment(apiGroup, logger, service.GetK8sClient, service.ListDeployments)
	api.RegisterStatefulSet(apiGroup, logger, service.GetK8sClient, service.ListStatefulSets)
	api.RegisterDaemonSet(apiGroup, logger, service.GetK8sClient, service.ListDaemonSets)

	// 网络资源API注册
	api.RegisterService(apiGroup, logger, service.GetK8sClient, service.ListServices)
	api.RegisterIngress(apiGroup, logger, service.GetK8sClient, service.ListIngresses)

	// 工作负载资源API注册
	api.RegisterCronJob(apiGroup, logger, service.GetK8sClient, service.ListCronJobs)
	api.RegisterJob(apiGroup, logger, service.GetK8sClient, service.ListJobs)

	// 集群资源API注册
	api.RegisterNamespace(apiGroup, logger, service.GetK8sClient, service.ListNamespaces)
	api.RegisterEvent(apiGroup, logger, service.GetK8sClient, service.ListEvents)

	// 存储资源API注册
	api.RegisterPVC(apiGroup, logger, service.GetK8sClient, service.ListPVCs)
	api.RegisterPV(apiGroup, logger, service.GetK8sClient, service.ListPVs)
	api.RegisterStorageClass(apiGroup, logger, service.GetK8sClient, service.ListStorageClasses)

	// 配置资源API注册
	api.RegisterConfigMap(apiGroup, logger, service.GetK8sClient, service.ListConfigMaps)
	api.RegisterSecret(apiGroup, logger, service.GetK8sClient, service.ListSecrets)

	// 启动性能监控定期日志记录
	monitorMgr.StartPeriodicLogging(5 * time.Minute)

	// 启动服务器
	serverAddr := cfg.GetServerAddress()
	logger.Info("服务器启动", 
		zap.String("address", serverAddr),
		zap.Bool("cacheEnabled", cfg.Cache.Enabled),
		zap.Bool("rateLimitEnabled", cfg.Auth.EnableRateLimit),
	)
	r.Run(serverAddr)
}
