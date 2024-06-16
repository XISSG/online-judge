package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(data interface{}) Response {
	return Response{
		Code:    200,
		Message: "Success",
		Data:    data,
	}
}

func ErrorResponse(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

func RecoveryMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer logger.Sync()
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("%v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse(500, "Internal Server Error"))
				c.Abort()
			}
		}()
		c.Next()
	}
}
func ResponseMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		defer logger.Sync()
		// 处理请求后的响应
		if c.Writer.Status() == http.StatusOK {
			if filePath, exists := c.Get("file"); exists {
				c.File(filePath.(string))
				// 返回成功响应
			} else {
				c.JSON(http.StatusOK, SuccessResponse(c.Keys["data"]))
			}
		}
	}
}
