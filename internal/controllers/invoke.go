package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xissg/online-judge/internal/middlewares"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/repository/redis"
	"go.uber.org/zap"
	"net/http"
)

type InvokeHandler struct {
	rdb       *redis.RedisClient
	logger    *zap.SugaredLogger
	validator *validator.Validate
}

func NewInvokeHandler(rdb *redis.RedisClient, logger *zap.SugaredLogger) *InvokeHandler {
	return &InvokeHandler{
		rdb:       rdb,
		logger:    logger,
		validator: validator.New(),
	}
}

// GetInvokeCount
//
//	@Summary		get api invoke count
//	@Description	get api invoke count
//	@Tags			invoke
//	@Accept			json
//	@Produce		json
//	@Param			invokeRequest	body		request.Invoke			true	"invoke request"
//	@Success		200				{object}	middlewares.Response	"ok"
//	@Failure		400				{object}	middlewares.Response	"bad request"
//	@Failure		500				{object}	middlewares.Response	"Internal Server Error"
//	@Router			/admin/invoke/count [post]
func (h *InvokeHandler) GetInvokeCount(ctx *gin.Context) {
	//获取请求数据
	invokeInfo := request.Invoke{}
	err := ctx.ShouldBindJSON(&invokeInfo)
	if err != nil {
		h.logger.Infof("invalid invoke information %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid invoke information"))
		return
	}

	//校验数据
	err = h.validator.Struct(invokeInfo)
	if err != nil {
		h.logger.Infof("invalid invoke information %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid invoke information"))
		return
	}

	//查询调用结果
	key := fmt.Sprintf("invoke:%s:%s", invokeInfo.Method, invokeInfo.Path)
	data, err := h.rdb.HGetAll(key)
	if err != nil {
		h.logger.Infof("failed to get invoke count %v", err)
		ctx.JSON(http.StatusInternalServerError, middlewares.ErrorResponse(http.StatusInternalServerError, "failed to get invoke count"))
		return
	}

	//返回结果
	ctx.Set("data", data)
}
