package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort      string            // 服务端口
	LogConfig       LogConfig         // 日志配置
	MySQL           MySQLConfig       // mysql配置
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
}

type MySQLConfig struct {
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

var (
	config *Config
)

func InitConfig() {
	bindEnvs()
	loadConfigFile()
}

func loadConfigFile() {
	v := viper.New()

	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.yaml"
	}

	v.SetConfigFile(configFile)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	config = &Config{}

	err := v.Unmarshal(&config)
	if err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}
}

func bindEnvs() {
	viper.AutomaticEnv()

	viper.BindEnv("ServerPort", "SERVER_PORT")
	viper.BindEnv("LogConfig.Level", "LOG_LEVEL")
	viper.BindEnv("LogConfig.Format", "LOG_FORMAT")
	viper.BindEnv("LogConfig.File", "LOG_FILE")
	viper.BindEnv("LogConfig.MaxSize", "LOG_MAX_SIZE")
	viper.BindEnv("LogConfig.MaxAge", "LOG_MAX_AGE")
	viper.BindEnv("LogConfig.Compress", "LOG_COMPRESS")
	viper.BindEnv("MySQL.Host", "MYSQL_HOST")
	viper.BindEnv("MySQL.Port", "MYSQL_PORT")
}

func GetConfig() *Config {
	return config
}
