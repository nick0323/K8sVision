package model

import (
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
	Port         string        `mapstructure:"port" json:"port"`
	Host         string        `mapstructure:"host" json:"host"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout" json:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout" json:"writeTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout" json:"idleTimeout"`
	MaxHeaderBytes int         `mapstructure:"maxHeaderBytes" json:"maxHeaderBytes"`
}

// KubernetesConfig Kubernetes配置
type KubernetesConfig struct {
	Kubeconfig     string        `mapstructure:"kubeconfig" json:"kubeconfig"`
	Context        string        `mapstructure:"context" json:"context"`
	Timeout        time.Duration `mapstructure:"timeout" json:"timeout"`
	QPS            float32       `mapstructure:"qps" json:"qps"`
	Burst          int           `mapstructure:"burst" json:"burst"`
	Insecure       bool          `mapstructure:"insecure" json:"insecure"`
	CAFile         string        `mapstructure:"caFile" json:"caFile"`
	CertFile       string        `mapstructure:"certFile" json:"certFile"`
	KeyFile        string        `mapstructure:"keyFile" json:"keyFile"`
	Token          string        `mapstructure:"token" json:"token"`
	APIServer      string        `mapstructure:"apiServer" json:"apiServer"`
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
	Enabled     bool          `mapstructure:"enabled" json:"enabled"`
	Type        string        `mapstructure:"type" json:"type"`
	TTL         time.Duration `mapstructure:"ttl" json:"ttl"`
	MaxSize     int           `mapstructure:"maxSize" json:"maxSize"`
	CleanupInterval time.Duration `mapstructure:"cleanupInterval" json:"cleanupInterval"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         "8080",
			Host:         "0.0.0.0",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
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
			Secret:     "k8svision-secret-key",
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
			Username:        "admin",
			Password:        "admin",
			MaxLoginFail:    5,
			LockDuration:    10 * time.Minute,
			SessionTimeout:  24 * time.Hour,
			EnableRateLimit: true,
			RateLimit:       100,
		},
		Cache: CacheConfig{
			Enabled:         true,
			Type:           "memory",
			TTL:            5 * time.Minute,
			MaxSize:        1000,
			CleanupInterval: 10 * time.Minute,
		},
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 这里可以添加配置验证逻辑
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