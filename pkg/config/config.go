package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort      string            // 服务端口
	LogConfig       LogConfig         // 日志配置
	MySQL           DBConfig          // mysql配置
	Postgres        DBConfig          // postgres配置
	DBType          string            // 数据库类型: mysql | postgres
	TencentCloud    TencenCloudConfig // 腾讯云配置
	PkgFileDir      string            // 包文件存放目录
	UserInfoAddress string            // 用户信息地址
	Upload          UploadConfig      // 上传配置
	PasswdKey       string            // 密码加密key
	MailConfig      MailConfig        // 邮件配置
	RedisConfig     RedisConfig       // redis配置
	NoRouterFoward  string            // 是否转发没有路由的请求
	ForwardPrefix   []string          // 转发前缀
	ForwardAddress  string            // 转发地址
	ApiPrefix       string            // api前缀
	AllowedOrigins  []string          // CORS 允许的来源（兼容旧配置）
	CORS            CORSConfig        // CORS 详细配置
	AppRateLimit    RateLimitConfig   // 应用级限流配置
}

type UploadConfig struct {
	Dir string
}

type LogConfig struct {
	Level        string // 日志级别
	Format       string // 日志格式
	MaxSize      int    // 最大文件大小（MB）
	MaxAge       int    // 最大文件保留天数
	Compress     bool   // 是否压缩
	Filename     string // 日志文件名
	ResponseSize int    // 字节
	ShowConsole  bool   // 是否显示在控制台
	// 采样与堆栈
	EnableSampling     bool   // 是否启用日志采样
	SamplingInitial    int    // 采样初始条数/秒
	SamplingThereafter int    // 之后每秒采样条数
	StacktraceLevel    string // 输出堆栈的级别（error|warn|panic 等）
}

type DBConfig struct {
	Host     string // 数据库地址
	Port     int    // 数据库端口
	Username string // 数据库用户名
	Password string // 数据库密码
	Database string // 数据库名
	ShowSQL  bool   // 是否显示SQL
}

type TencenCloudConfig struct {
	SecretId    string
	SecretKey   string
	Region      string
	FileAddress string
}

// 邮件配置
type MailConfig struct {
	Host         string // 邮件服务器地址
	Port         int    // 邮件服务器端口
	Username     string // 邮件服务器用户名
	Password     string // 邮件服务器密码
	RegisterTile string // 注册邮件标题
}

type SmartGptConfig struct {
	Address string
	Key     string
}

type RedisConfig struct {
	Address  string // 地址, 多个使用逗号(,)分隔
	Password string
	Database int
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int // seconds
}

type RateLimitConfig struct {
	R int // 每秒令牌数
	B int // 桶容量
}

var (
	config *Config
	v      *viper.Viper
	// extensions holds registered config extension callbacks
	extensions = make(map[string]struct {
		fn     func(v *viper.Viper, cfg *Config) error
		strict bool
	})
)

func InitConfig() {
	bindEnvs()
	loadConfigFile()
}

// InitConfigWithFile sets the CONFIG_FILE environment variable and then
// initializes configuration. This is a convenience helper for callers who
// want to provide an explicit config file path.
func InitConfigWithFile(path string) {
	if path == "" {
		InitConfig()
		return
	}
	os.Setenv("CONFIG_FILE", path)
	bindEnvs()
	loadConfigFile()
}

func loadConfigFile() {
	v = viper.New()

	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.yaml"
	}

	v.SetConfigFile(configFile)

	if err := v.ReadInConfig(); err != nil {
		// 配置文件不存在时，回退到仅使用环境变量（容器化部署常见）
		log.Printf("failed to read config file: %v, falling back to environment variables", err)
	}

	config = &Config{}

	err := v.Unmarshal(&config)
	if err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}

	// viper 的 AutomaticEnv 不会注入到 Unmarshal，所以在文件缺失时手动从环境变量回填
	if config.ServerPort == "" {
		if v := os.Getenv("SERVER_PORT"); v != "" {
			config.ServerPort = v
		}
	}
	if config.DBType == "" {
		if v := os.Getenv("DB_TYPE"); v != "" {
			config.DBType = v
		}
	}
	if config.Postgres.Host == "" {
		if v := os.Getenv("POSTGRES_HOST"); v != "" {
			config.Postgres.Host = v
		}
	}
	if config.Postgres.Port == 0 {
		if v := os.Getenv("POSTGRES_PORT"); v != "" {
			fmt.Sscanf(v, "%d", &config.Postgres.Port)
		}
	}
	if config.Postgres.Username == "" {
		config.Postgres.Username = os.Getenv("POSTGRES_USER")
	}
	if config.Postgres.Password == "" {
		config.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	}
	if config.Postgres.Database == "" {
		config.Postgres.Database = os.Getenv("POSTGRES_DB")
	}
	if config.RedisConfig.Address == "" {
		config.RedisConfig.Address = os.Getenv("REDIS_ADDR")
	}
	if config.RedisConfig.Password == "" {
		config.RedisConfig.Password = os.Getenv("REDIS_PASSWORD")
	}
	if config.Upload.Dir == "" {
		config.Upload.Dir = os.Getenv("UPLOAD_DIR")
	}

	// 兼容从环境变量注入的逗号分隔形式的 AllowedOrigins（旧字段）
	if len(config.AllowedOrigins) == 1 && strings.Contains(config.AllowedOrigins[0], ",") {
		config.AllowedOrigins = splitAndTrim(config.AllowedOrigins[0])
	}
	// 同步到新字段（若新字段未配置）
	if len(config.CORS.AllowedOrigins) == 0 && len(config.AllowedOrigins) > 0 {
		config.CORS.AllowedOrigins = append([]string(nil), config.AllowedOrigins...)
	}

	// 兼容从环境变量注入的逗号分隔形式的 CORS 其他字段
	if len(config.CORS.AllowMethods) == 1 && strings.Contains(config.CORS.AllowMethods[0], ",") {
		config.CORS.AllowMethods = splitAndTrim(config.CORS.AllowMethods[0])
	}
	if len(config.CORS.AllowHeaders) == 1 && strings.Contains(config.CORS.AllowHeaders[0], ",") {
		config.CORS.AllowHeaders = splitAndTrim(config.CORS.AllowHeaders[0])
	}
	if len(config.CORS.ExposeHeaders) == 1 && strings.Contains(config.CORS.ExposeHeaders[0], ",") {
		config.CORS.ExposeHeaders = splitAndTrim(config.CORS.ExposeHeaders[0])
	}

	// 启动时做最小化配置校验与提示
	if err := config.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	// 执行已注册的扩展回调（插件式解析）
	for name, ext := range extensions {
		if ext.fn == nil {
			continue
		}
		if err := ext.fn(v, config); err != nil {
			if ext.strict {
				log.Fatalf("extension %s init failed: %v", name, err)
			} else {
				log.Printf("warning: extension %s init failed: %v", name, err)
			}
		}
	}
}

// RegisterExtension allows external packages to register a callback that will be
// executed after the main configuration file is loaded. The callback receives the
// underlying *viper.Viper and the decoded *Config. If strict is true, a non-nil
// error returned from the callback will cause application startup to fail.
func RegisterExtension(name string, fn func(v *viper.Viper, cfg *Config) error, strict bool) {
	if name == "" || fn == nil {
		return
	}
	extensions[name] = struct {
		fn     func(v *viper.Viper, cfg *Config) error
		strict bool
	}{fn: fn, strict: strict}
}

func bindEnvs() {
	viper.AutomaticEnv()

	viper.BindEnv("ServerPort", "SERVER_PORT")
	viper.BindEnv("LogConfig.Level", "LOG_LEVEL")
	viper.BindEnv("LogConfig.Format", "LOG_FORMAT")
	viper.BindEnv("LogConfig.Filename", "LOG_FILE")
	viper.BindEnv("LogConfig.MaxSize", "LOG_MAX_SIZE")
	viper.BindEnv("LogConfig.MaxAge", "LOG_MAX_AGE")
	viper.BindEnv("LogConfig.Compress", "LOG_COMPRESS")
	viper.BindEnv("LogConfig.EnableSampling", "LOG_SAMPLING_ENABLE")
	viper.BindEnv("LogConfig.SamplingInitial", "LOG_SAMPLING_INITIAL")
	viper.BindEnv("LogConfig.SamplingThereafter", "LOG_SAMPLING_THEREAFTER")
	viper.BindEnv("LogConfig.StacktraceLevel", "LOG_STACKTRACE_LEVEL")
	viper.BindEnv("MySQL.Host", "MYSQL_HOST")
	viper.BindEnv("MySQL.Port", "MYSQL_PORT")
	viper.BindEnv("Postgres.Host", "POSTGRES_HOST")
	viper.BindEnv("Postgres.Port", "POSTGRES_PORT")
	viper.BindEnv("Postgres.Username", "POSTGRES_USER")
	viper.BindEnv("Postgres.Password", "POSTGRES_PASSWORD")
	viper.BindEnv("Postgres.Database", "POSTGRES_DB")
	viper.BindEnv("DBType", "DB_TYPE")
	viper.BindEnv("Redis.Address", "REDIS_ADDR")
	viper.BindEnv("Redis.Password", "REDIS_PASSWORD")
	viper.BindEnv("Redis.Database", "REDIS_DB")
	viper.BindEnv("AllowedOrigins", "ALLOWED_ORIGINS")
	viper.BindEnv("PasswdKey", "PASSWD_KEY")
	viper.BindEnv("Upload.Dir", "UPLOAD_DIR")
	// CORS 详细配置
	viper.BindEnv("CORS.AllowedOrigins", "CORS_ALLOWED_ORIGINS")
	viper.BindEnv("CORS.AllowMethods", "CORS_ALLOW_METHODS")
	viper.BindEnv("CORS.AllowHeaders", "CORS_ALLOW_HEADERS")
	viper.BindEnv("CORS.ExposeHeaders", "CORS_EXPOSE_HEADERS")
	viper.BindEnv("CORS.AllowCredentials", "CORS_ALLOW_CREDENTIALS")
	viper.BindEnv("CORS.MaxAge", "CORS_MAX_AGE")
	// 限流配置
	viper.BindEnv("AppRateLimit.R", "APP_RATE_LIMIT_R")
	viper.BindEnv("AppRateLimit.B", "APP_RATE_LIMIT_B")
}

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

// 基础配置校验（最小化约束，避免误伤现有用法）
func (c *Config) Validate() error {
	if c.ServerPort == "" {
		c.ServerPort = "8080"
	}
	if c.PasswdKey == "" {
		// 警告：未配置 PasswdKey，将回退到代码中的默认值，不建议在生产环境使用
		log.Printf("warning: PasswdKey is empty, fallback key will be used. Please set PASSWD_KEY in production.")
	}
	if c.AppRateLimit.R <= 0 {
		c.AppRateLimit.R = 100
	}
	if c.AppRateLimit.B <= 0 {
		c.AppRateLimit.B = 200
	}
	// 允许不配置 DB/Redis；JWT 密钥若未配置则在 utils 里回退旧值
	return nil
}

func GetConfig() *Config {
	return config
}

// UnmarshalKey decodes a specific key from the underlying viper into out.
// Returns an error if config hasn't been initialized.
func UnmarshalKey(key string, out interface{}) error {
	if v == nil {
		return errors.New("config not initialized")
	}
	return v.UnmarshalKey(key, out)
}

// Raw returns the underlying viper instance (read-only use recommended).
func Raw() *viper.Viper {
	return v
}
