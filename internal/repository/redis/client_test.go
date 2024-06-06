package redis

import (
	"github.com/xissg/online-judge/internal/config"
	"testing"
)

func TestRedisClient(t *testing.T) {
	config.LoadConfig()
	NewRedisClient(config.AppConfig.Redis)
}
