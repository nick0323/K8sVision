package model

import (
	"fmt"
	"time"
)

// Config 应用配置结构体
type Config struct {
	Server     ServerConfig     `mapstructure:"server" json:"server"`
	Kubernetes KubernetesConfig `mapstructure:"kubernetes" json:"kubernetes"`
	JWT        JWTConfig        `mapstructure:"jwt" json:"jwt"`
	Log        LogConfig        `mapstructure:"log" json:"log"`
	Auth       AuthConfig       `mapstructure:"auth" json:"auth"`
	Cache      CacheConfig      `mapstructure:"cache" json:"cache"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port           string        `mapstructure:"port" json:"port"`
	Host           string        `mapstructure:"host" json:"host"`
	ReadTimeout    time.Duration `mapstructure:"readTimeout" json:"readTimeout"`
	WriteTimeout   time.Duration `mapstructure:"writeTimeout" json:"writeTimeout"`
	IdleTimeout    time.Duration `mapstructure:"idleTimeout" json:"idleTimeout"`
	MaxHeaderBytes int           `mapstructure:"maxHeaderBytes" json:"maxHeaderBytes"`
}

// KubernetesConfig Kubernetes配置
type KubernetesConfig struct {
	Kubeconfig string        `mapstructure:"kubeconfig" json:"kubeconfig"`
	Context    string        `mapstructure:"context" json:"context"`
	Timeout    time.Duration `mapstructure:"timeout" json:"timeout"`
	QPS        float32       `mapstructure:"qps" json:"qps"`
	Burst      int           `mapstructure:"burst" json:"burst"`
	Insecure   bool          `mapstructure:"insecure" json:"insecure"`
	CAFile     string        `mapstructure:"caFile" json:"caFile"`
	CertFile   string        `mapstructure:"certFile" json:"certFile"`
	KeyFile    string        `mapstructure:"keyFile" json:"keyFile"`
	Token      string        `mapstructure:"token" json:"token"`
	APIServer  string        `mapstructure:"apiServer" json:"apiServer"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `mapstructure:"secret" json:"secret"`
	Expiration time.Duration `mapstructure:"expiration" json:"expiration"`
	Issuer     string        `mapstructure:"issuer" json:"issuer"`
	Audience   string        `mapstructure:"audience" json:"audience"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level" json:"level"`
	Format     string `mapstructure:"format" json:"format"`
	Output     string `mapstructure:"output" json:"output"`
	MaxSize    int    `mapstructure:"maxSize" json:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups" json:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge" json:"maxAge"`
	Compress   bool   `mapstructure:"compress" json:"compress"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Username        string        `mapstructure:"username" json:"username"`
	Password        string        `mapstructure:"password" json:"password"`
	MaxLoginFail    int           `mapstructure:"maxLoginFail" json:"maxLoginFail"`
	LockDuration    time.Duration `mapstructure:"lockDuration" json:"lockDuration"`
	SessionTimeout  time.Duration `mapstructure:"sessionTimeout" json:"sessionTimeout"`
	EnableRateLimit bool          `mapstructure:"enableRateLimit" json:"enableRateLimit"`
	RateLimit       int           `mapstructure:"rateLimit" json:"rateLimit"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled         bool          `mapstructure:"enabled" json:"enabled"`
	Type            string        `mapstructure:"type" json:"type"`
	TTL             time.Duration `mapstructure:"ttl" json:"ttl"`
	MaxSize         int           `mapstructure:"maxSize" json:"maxSize"`
	CleanupInterval time.Duration `mapstructure:"cleanupInterval" json:"cleanupInterval"`
}

// DefaultConfig 返回系统默认配置
// 包含服务器、Kubernetes、JWT、日志、认证和缓存的默认设置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:           "8080",
			Host:           "0.0.0.0",
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1MB
		},
		Kubernetes: KubernetesConfig{
			Kubeconfig: "",
			Context:    "",
			Timeout:    30 * time.Second,
			QPS:        100,
			Burst:      200,
			Insecure:   true,
		},
		JWT: JWTConfig{
			Secret:     "k8svision-default-jwt-secret-key-32-chars", // 默认密钥，生产环境请设置环境变量 K8SVISION_JWT_SECRET
			Expiration: 24 * time.Hour,
			Issuer:     "k8svision",
			Audience:   "k8svision-client",
		},
		Log: LogConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		},
		Auth: AuthConfig{
			Username:        "admin",     // 默认用户名，生产环境请设置环境变量 K8SVISION_AUTH_USERNAME
			Password:        "admin123!", // 默认密码，生产环境请设置环境变量 K8SVISION_AUTH_PASSWORD
			MaxLoginFail:    5,
			LockDuration:    10 * time.Minute,
			SessionTimeout:  24 * time.Hour,
			EnableRateLimit: true,
			RateLimit:       100,
		},
		Cache: CacheConfig{
			Enabled:         true,
			Type:            "memory",
			TTL:             5 * time.Minute,
			MaxSize:         1000,
			CleanupInterval: 10 * time.Minute,
		},
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证服务器配置
	if c.Server.Port == "" {
		return fmt.Errorf("服务器端口不能为空")
	}
	if c.Server.Host == "" {
		return fmt.Errorf("服务器主机不能为空")
	}
	if c.Server.ReadTimeout <= 0 {
		return fmt.Errorf("读取超时时间必须大于0")
	}
	if c.Server.WriteTimeout <= 0 {
		return fmt.Errorf("写入超时时间必须大于0")
	}

	// 验证JWT配置
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空，请设置环境变量 K8SVISION_JWT_SECRET 或在配置文件中设置 jwt.secret")
	}
	if len(c.JWT.Secret) < 16 {
		return fmt.Errorf("JWT密钥长度至少16位字符，当前长度: %d", len(c.JWT.Secret))
	}
	if c.JWT.Expiration <= 0 {
		return fmt.Errorf("JWT过期时间必须大于0")
	}

	// 验证认证配置
	if c.Auth.Username == "" {
		return fmt.Errorf("认证用户名不能为空，请设置环境变量 LOGIN_USERNAME 或在配置文件中设置 auth.username")
	}
	if c.Auth.Password == "" {
		return fmt.Errorf("认证密码不能为空，请设置环境变量 LOGIN_PASSWORD 或在配置文件中设置 auth.password")
	}
	if c.Auth.MaxLoginFail <= 0 {
		return fmt.Errorf("最大登录失败次数必须大于0")
	}
	if c.Auth.LockDuration <= 0 {
		return fmt.Errorf("锁定时间必须大于0")
	}

	// 验证日志配置
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[c.Log.Level] {
		return fmt.Errorf("无效的日志级别: %s", c.Log.Level)
	}

	validLogFormats := map[string]bool{
		"json": true, "console": true,
	}
	if !validLogFormats[c.Log.Format] {
		return fmt.Errorf("无效的日志格式: %s", c.Log.Format)
	}

	// 验证缓存配置
	if c.Cache.Enabled {
		if c.Cache.TTL <= 0 {
			return fmt.Errorf("缓存TTL必须大于0")
		}
		if c.Cache.MaxSize <= 0 {
			return fmt.Errorf("缓存最大大小必须大于0")
		}
	}

	// 验证Kubernetes配置
	if c.Kubernetes.QPS <= 0 {
		return fmt.Errorf("kubernetes QPS必须大于0")
	}
	if c.Kubernetes.Burst <= 0 {
		return fmt.Errorf("kubernetes Burst必须大于0")
	}
	if c.Kubernetes.Timeout <= 0 {
		return fmt.Errorf("Kubernetes超时时间必须大于0")
	}

	return nil
}

// GetServerAddress 获取服务器地址
func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

// IsDevelopment 判断是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.Log.Level == "debug"
}

// IsProduction 判断是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Log.Level == "info" || c.Log.Level == "warn" || c.Log.Level == "error"
}
