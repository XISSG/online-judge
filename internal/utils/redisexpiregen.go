package utils

import (
	"github.com/xissg/online-judge/internal/constant"
	"math/rand"
	"time"
)

func RandomExpireTime() time.Duration {
	rand.Seed(time.Now().UnixNano())
	minExpire := constant.MinExpire
	maxExpire := constant.MaxExpire
	expire := minExpire + time.Duration(rand.Int63n(int64(maxExpire-minExpire)))
	return expire
}
