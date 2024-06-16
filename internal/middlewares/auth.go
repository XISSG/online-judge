package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
	"github.com/xissg/online-judge/internal/utils"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func AuthLogin(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, ok := parseAuthInfo(c)
		if !ok {
			logger.Infof("user authentication failed")
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "user authentication failed"))
			c.Abort()
			return
		}
		c.Set("user", data.ID)
		c.Next()
	}
}

func AuthAdmin(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, ok := parseAuthInfo(c)
		if !ok {
			logger.Infof("admin user authentication failed")
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "user authentication failed"))
			c.Abort()
			return
		}

		if data.UserRole != constant.ADMIN {
			logger.Infof("admin user authentication failed")
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "user authentication failed"))
			c.Abort()
			return
		}
		c.Set("user", data.ID)
		c.Next()
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
