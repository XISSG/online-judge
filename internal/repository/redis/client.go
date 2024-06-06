package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/xissg/online-judge/internal/config"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(cfg config.RedisConfig) *RedisClient {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password, // 没有密码，默认值
		DB:       cfg.DB,       // 默认DB 0
	})
	return &RedisClient{
		rdb,
	}
}
