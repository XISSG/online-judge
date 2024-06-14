package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORS() gin.HandlerFunc {

	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.GetHeader("Origin")
		c.Header("Access-Control-Allow-Origin", origin)                                                                                                                  // 不能配置为通配符“*”号
		c.Header("Access-Control-Allow-Credentials", "true")                                                                                                             // 必须设定为 true
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Cookie, Origin, X-Requested-With, Content-Type, Accept, Authorization, Token, Timestamp") // 自定义的header字段都需要在此声明
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,cache-control")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		// 处理请求
		c.Next()
	}
}
