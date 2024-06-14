package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
	"github.com/xissg/online-judge/internal/utils"
	"time"
)

func AuthLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, ok := parseAuthInfo(c)
		if !ok {
			c.Abort()
			return
		}
		c.Set("user", data.ID)
		c.Next()
	}
}

func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, ok := parseAuthInfo(c)
		if !ok && data.UserRole != constant.ADMIN {
			c.Abort()
			return
		}

		c.Set("user", data.ID)
	}
}

func parseAuthInfo(c *gin.Context) (*entity.JwtData, bool) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return nil, false
	}
	jwtConfig := config.LoadConfig().Jwt
	data, err := utils.Parse(tokenString, jwtConfig)
	if err != nil {
		return nil, false
	}
	if data.Expiration.Unix() < time.Now().Unix() {
		return nil, false
	}
	return data, true
}
