package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server        ServerConfig        `mapstructure:"server"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Redis         RedisConfig         `mapstructure:"redis"`
	JWT           JWTConfig           `mapstructure:"jwt"`
	SMS           SMSConfig           `mapstructure:"sms"`
	OSS           OSSConfig           `mapstructure:"oss"`
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
	RabbitMQ      RabbitMQConfig      `mapstructure:"rabbitmq"`
	Wallet        WalletConfig        `mapstructure:"wallet"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type SMSConfig struct {
	Provider        string `mapstructure:"provider"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	SignName        string `mapstructure:"sign_name"`
	TemplateCode    string `mapstructure:"template_code"`
}

type OSSConfig struct {
	Provider        string `mapstructure:"provider"`
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	Bucket          string `mapstructure:"bucket"`
}

type ElasticsearchConfig struct {
	Addresses []string `mapstructure:"addresses"`
}

type RabbitMQConfig struct {
	URL string `mapstructure:"url"`
}

type WalletConfig struct {
	InitialReward float64 `mapstructure:"initial_reward"`
}

var Global *Config

// Load 加载配置。优先读配置文件（若存在），环境变量可覆盖：
// 前缀 NEOBARTER_，嵌套键用下划线，如 NEOBARTER_DATABASE_HOST=postgres。
// 配置文件不存在时，仅靠环境变量 + 默认值也能启动（容器场景）。
func Load(path string) (*Config, error) {
	setDefaults()

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("NEOBARTER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 配置文件可选：找不到文件不报错，用环境变量 + 默认值兜底
	if err := viper.ReadInConfig(); err != nil {
		if _, notFound := err.(viper.ConfigFileNotFoundError); !notFound && !os.IsNotExist(err) {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// 切片/特殊类型的环境变量单独处理（viper AutomaticEnv 对 []string 支持不佳）
	// NEOBARTER_ELASTICSEARCH_ADDRESSES 支持逗号分隔多地址
	if v := os.Getenv("NEOBARTER_ELASTICSEARCH_ADDRESSES"); v != "" {
		cfg.Elasticsearch.Addresses = strings.Split(v, ",")
	}
	if v := os.Getenv("NEOBARTER_RABBITMQ_URL"); v != "" {
		cfg.RabbitMQ.URL = v
	}

	Global = &cfg
	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "neobarter")
	viper.SetDefault("database.password", "neobarter123")
	viper.SetDefault("database.dbname", "neobarter")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("jwt.secret", "change-me-in-production")
	viper.SetDefault("jwt.expire_hours", 168)
	viper.SetDefault("oss.provider", "local")
	viper.SetDefault("wallet.initial_reward", 100.0)
}
