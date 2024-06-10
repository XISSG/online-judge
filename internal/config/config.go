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
	Image         ImageConfig
	AI            AIConfig
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
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type AIConfig struct {
	HostUrl   string `mapstructure:"host_url"`
	AppId     string `mapstructure:"app_id"`
	ApiSecret string `mapstructure:"api_secret"`
	ApiKey    string `mapstructure:"api_key"`
}

type ImageConfig map[string]string

func LoadConfig() Config {
	//main执行的路径
	//viper.AddConfigPath("internal/config/")
	//service层执行的路径
	viper.AddConfigPath("../config/")
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
