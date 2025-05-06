package config

import (
	"github.com/Cospk/go-mall/pkg/config"
	"time"
)

// 定义配置结构体
var (
	App              AppConfig
	ServiceDiscovery ServiceDiscoveryConfig
	Services         ServicesConfig
	Auth             AuthConfig
	RateLimit        RateLimitConfig
	Cache            CacheConfig
	Cors             CorsConfig
	Monitor          MonitorConfig
)

// 初始化配置
func InitApiConfig() {
	config.InitConfig("api", map[string]interface{}{
		"app":               &App,
		"service_discovery": &ServiceDiscovery,
		"services":          &Services,
		"auth":              &Auth,
		"rate_limit":        &RateLimit,
		"cache":             &Cache,
		"cors":              &Cors,
		"monitor":           &Monitor,
	})
}

// AppConfig 应用基础配置
type AppConfig struct {
	Name        string `mapstructure:"name"`         // 应用名称
	Version     string `mapstructure:"version"`      // 应用版本
	Mode        string `mapstructure:"mode"`         // 运行模式：dev, test, prod
	Port        int    `mapstructure:"port"`         // 监听端口
	Host        string `mapstructure:"host"`         // 监听地址
	LogLevel    string `mapstructure:"log_level"`    // 日志级别
	LogPath     string `mapstructure:"log_path"`     // 日志路径
	EnablePProf bool   `mapstructure:"enable_pprof"` // 是否启用pprof
}

// ServiceDiscoveryConfig 服务发现配置
type ServiceDiscoveryConfig struct {
	Type      string   `mapstructure:"type"`      // 服务发现类型：etcd, consul, nacos
	Endpoints []string `mapstructure:"endpoints"` // 服务发现地址
	Namespace string   `mapstructure:"namespace"` // 命名空间
	Username  string   `mapstructure:"username"`  // 用户名
	Password  string   `mapstructure:"password"`  // 密码
	TTL       int      `mapstructure:"ttl"`       // 服务注册TTL
}

// ServicesConfig 微服务配置
type ServicesConfig struct {
	User struct {
		Address     string `mapstructure:"address"`      // 用户服务地址
		Timeout     int    `mapstructure:"timeout"`      // 超时时间(ms)
		MaxRetries  int    `mapstructure:"max_retries"`  // 最大重试次数
		PoolSize    int    `mapstructure:"pool_size"`    // 连接池大小
		HealthCheck bool   `mapstructure:"health_check"` // 是否启用健康检查
	} `mapstructure:"user"`

	Product struct {
		Address     string `mapstructure:"address"`      // 商品服务地址
		Timeout     int    `mapstructure:"timeout"`      // 超时时间(ms)
		MaxRetries  int    `mapstructure:"max_retries"`  // 最大重试次数
		PoolSize    int    `mapstructure:"pool_size"`    // 连接池大小
		HealthCheck bool   `mapstructure:"health_check"` // 是否启用健康检查
	} `mapstructure:"product"`

	Order struct {
		Address     string `mapstructure:"address"`      // 订单服务地址
		Timeout     int    `mapstructure:"timeout"`      // 超时时间(ms)
		MaxRetries  int    `mapstructure:"max_retries"`  // 最大重试次数
		PoolSize    int    `mapstructure:"pool_size"`    // 连接池大小
		HealthCheck bool   `mapstructure:"health_check"` // 是否启用健康检查
	} `mapstructure:"order"`

	// 其他微服务...
}

// AuthConfig 认证配置
type AuthConfig struct {
	JWT struct {
		Secret     string        `mapstructure:"secret"`     // JWT密钥
		Issuer     string        `mapstructure:"issuer"`     // 签发者
		Expiration time.Duration `mapstructure:"expiration"` // 过期时间
	} `mapstructure:"jwt"`

	OAuth struct {
		Enabled      bool   `mapstructure:"enabled"`       // 是否启用OAuth
		ClientID     string `mapstructure:"client_id"`     // 客户端ID
		ClientSecret string `mapstructure:"client_secret"` // 客户端密钥
		RedirectURL  string `mapstructure:"redirect_url"`  // 重定向URL
	} `mapstructure:"oauth"`

	WhiteList []string `mapstructure:"white_list"` // 白名单路径
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled     bool     `mapstructure:"enabled"`       // 是否启用限流
	Type        string   `mapstructure:"type"`          // 限流类型：token_bucket, leaky_bucket
	Rate        int      `mapstructure:"rate"`          // 速率
	Burst       int      `mapstructure:"burst"`         // 突发流量
	RedisKey    string   `mapstructure:"redis_key"`     // Redis键前缀
	IPWhiteList []string `mapstructure:"ip_white_list"` // IP白名单
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Type            string        `mapstructure:"type"`             // 缓存类型：redis, memory
	Expiration      time.Duration `mapstructure:"expiration"`       // 过期时间
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"` // 清理间隔
	MaxSize         int           `mapstructure:"max_size"`         // 最大缓存条目数
}

// CorsConfig 跨域配置
type CorsConfig struct {
	Enabled          bool     `mapstructure:"enabled"`           // 是否启用跨域
	AllowOrigins     []string `mapstructure:"allow_origins"`     // 允许的源
	AllowMethods     []string `mapstructure:"allow_methods"`     // 允许的方法
	AllowHeaders     []string `mapstructure:"allow_headers"`     // 允许的头
	ExposeHeaders    []string `mapstructure:"expose_headers"`    // 暴露的头
	AllowCredentials bool     `mapstructure:"allow_credentials"` // 是否允许凭证
	MaxAge           int      `mapstructure:"max_age"`           // 预检请求缓存时间
}

// MonitorConfig 监控配置
type MonitorConfig struct {
	Prometheus struct {
		Enabled bool   `mapstructure:"enabled"` // 是否启用Prometheus
		Path    string `mapstructure:"path"`    // 指标路径
	} `mapstructure:"prometheus"`

	Tracing struct {
		Enabled    bool    `mapstructure:"enabled"`     // 是否启用链路追踪
		Type       string  `mapstructure:"type"`        // 追踪类型：jaeger, zipkin
		Endpoint   string  `mapstructure:"endpoint"`    // 追踪服务地址
		SampleRate float64 `mapstructure:"sample_rate"` // 采样率
	} `mapstructure:"tracing"`
}
