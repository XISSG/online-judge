package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/repository/redis"
	"go.uber.org/ratelimit"
	"time"
)

func InvokeLimit(limiter ratelimit.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter.Take()
		c.Next()
	}

}
func InvokeCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		path := c.Request.URL.Path
		method := c.Request.Method

		redisConfig := config.LoadConfig().Redis
		rdb := redis.NewRedisClient(redisConfig)

		key := fmt.Sprintf("invoke:%s:%s", method, path)
		rdb.HIncrBy(key, "count", 1)
		rdb.HIncrBy(key, "invoke_time_cost", duration.Milliseconds())
	}
}
