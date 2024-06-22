package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Elasticsearch ElasticsearchConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	RabbitMQ      RabbitMQConfig
	Server        ServerConfig
	Images        ImageConfig
	AI            AIConfig
	Log           LogConfig
	Jwt           JWTConfig
	Picture       PictureConfig
}

type ElasticsearchConfig struct {
	Addresses []string `mapstructure:"addresses"`
	Username  string   `mapstructure:"username"`
	Password  string   `mapstructure:"password"`
}

type DatabaseConfig struct {
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	DBName          string        `mapstructure:"db_name"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RabbitMQConfig struct {
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	PublishType  string `mapstructure:"publish_type"`
	ExchangeName string `mapstructure:"exchange_name"`
	RoutingKey   string `mapstructure:"routing_key"`
	QueueName    string `mapstructure:"queue_name"`
	ConsumerTag  string `mapstructure:"consumer_tag"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// AIConfig
//
// Spark3.5 Max请求地址，对应的domain参数为generalv3.5：
// wss://spark-api.xf-yun.com/v3.5/chat
//
// Spark Pro请求地址，对应的domain参数为generalv3：
// wss://spark-api.xf-yun.com/v3.1/chat
//
// Spark V2.0请求地址，对应的domain参数为generalv2：
// wss://spark-api.xf-yun.com/v2.1/chat
//
// Spark Lite请求地址，对应的domain参数为general：
// wss://spark-api.xf-yun.com/v1.1/chat
type AIConfig struct {
	HostUrl     string  `mapstructure:"host_url"` //访问的地址：
	AppId       string  `mapstructure:"app_id"`
	ApiSecret   string  `mapstructure:"api_secret"`
	ApiKey      string  `mapstructure:"api_key"`
	Flexibility int     `mapstructure:"flexibility"` //灵活性和topk同义
	Randomness  float64 `mapstructure:"randomness"`  //随机性和temperature同义
	MaxTokens   int     `mapstructure:"max_tokens"`  //限制最大token数
	Domanin     string  `mapstructure:"domain"`      //设置访问的ai版本,有general,generalv2,generalv3,generalv3.5四个版本
}

type ImageConfig map[string]string

type LogConfig struct {
	OutputPaths      []string `mapstructure:"output_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type PictureConfig struct {
	Path string `mapstructure:"path"`
}

func LoadConfig() Config {
	//main执行的路径
	viper.AddConfigPath("internal/config/")
	//service层执行的路径
	//viper.AddConfigPath("../config/")
	//repository层执行的路径
	//viper.AddConfigPath("../../config/")
	//viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var appConfig Config
	if err := viper.Unmarshal(&appConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return appConfig
}
