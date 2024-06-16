package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xissg/online-judge/internal/middlewares"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/service"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type SubmitHandler struct {
	submitService   service.SubmitService
	rabbitMqService service.RabbitMqService
	logger          *zap.SugaredLogger
	validator       *validator.Validate
}

func NewSubmitHandler(submitService service.SubmitService, rabbitMqService service.RabbitMqService, logger *zap.SugaredLogger) *SubmitHandler {
	return &SubmitHandler{
		submitService:   submitService,
		rabbitMqService: rabbitMqService,
		logger:          logger,
		validator:       validator.New(),
	}
}

// CreateSubmit
//
//	@Summary		Create submit
//	@Description	Create submit
//	@Tags			submit
//	@Accept			json
//	@Produce		json
//	@Param			createRequest	body		request.Submit			true	"create submit"
//	@Success		200				{object}	middlewares.Response	"ok"
//	@Failure		400				{object}	middlewares.Response	"bad request"
//	@Failure		500				{object}	middlewares.Response	"Internal Server Error"
//	@Router			/user/submit/create_submit [post]
func (r *SubmitHandler) CreateSubmit(ctx *gin.Context) {
	//判断用户是否已登录
	userId, exist := ctx.Get("user")
	if !exist {
		r.logger.Infof("authorization error")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "you must Login first"))
		return
	}

	//获取请求数据
	req := &request.Submit{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		r.logger.Infof("submit request error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "submit request error"))
		return
	}

	//验证请求数据
	err = r.validator.Struct(req)
	if err != nil {
		r.logger.Infof("submit validate error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "submit validate error"))
		return
	}

	//创建提交信息
	idInt, err := r.submitService.CreateSubmit(req, userId.(int))
	if err != nil {
		r.logger.Infof("submit create error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "submit create error"))
		return
	}

	//将提交信息推送到消息队列
	id := strconv.FormatInt(int64(idInt), 10)
	r.logger.Infof("publish submit id")
	err = r.rabbitMqService.Publish(id)
	if err != nil {
		r.logger.Infof("submit publish error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "submit publish error"))
		return
	}
	ctx.Set("data", "submit success")
}

// GetSubmitList
//
//	@Summary		get submit list
//	@Description	get submit list
//	@Tags			submit
//	@Accept			json
//	@Produce		json
//	@Param			page		query		string					true	"page number"
//	@Param			page_size	query		string					true	"page size"
//	@Success		200			{object}	middlewares.Response	"ok"
//	@Failure		400			{object}	middlewares.Response	"bad request"
//	@Failure		500			{object}	middlewares.Response	"Internal Server Error"
//	@Router			/user/submit/get_submits [get]
func (r *SubmitHandler) GetSubmitList(ctx *gin.Context) {
	//获取请求数据
	page := ctx.Query("page")
	pageSize := ctx.Query("page_size")

	//校验数据
	if page == "" || pageSize == "" {
		r.logger.Infof("invalid submit page or page size")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid submit page or page size"))
		return
	}
	pageInt, err := strconv.Atoi(page)
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		r.logger.Infof("invalid submit page or page size")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid submit page or page size"))
		return
	}

	//获取提交信息
	submits, err := r.submitService.GetSubmitList(pageInt, pageSizeInt)
	if err != nil {
		r.logger.Errorf("get submit list error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "get submit list error"))
		return
	}

	//返回结果
	ctx.Set("data", submits)
}

// SearchSubmitList
//
//	@Summary		Search submit
//	@Description	Search submit
//	@Tags			submit
//	@Accept			json
//	@Produce		json
//	@Param			keyword	query		string					true	"search submit"
//	@Success		200		{object}	middlewares.Response	"ok"
//	@Failure		400		{object}	middlewares.Response	"bad request"
//	@Failure		500		{object}	middlewares.Response	"Internal Server Error"
//	@Router			/user/submit/search_submits [get]
func (r *SubmitHandler) SearchSubmitList(ctx *gin.Context) {
	//获取请求数据
	keyword := ctx.Query("keyword")
	if keyword == "" {
		r.logger.Infof("invalid search keyword")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid search keyword"))
	}

	//搜索提交信息
	submits, err := r.submitService.SearchSubmit(keyword)
	if err != nil {
		r.logger.Errorf("search submit error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "search submit error"))
		return
	}

	//返回结果
	ctx.Set("data", submits)
}

// DeleteSubmit
//
//	@Summary		Delete submit
//	@Description	Delete submit
//	@Tags			submit
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string					true	"delete submit"
//	@Success		200	{object}	middlewares.Response	"ok"
//	@Failure		400	{object}	middlewares.Response	"bad request"
//	@Failure		500	{object}	middlewares.Response	"Internal Server Error"
//	@Router			/admin/submit/delete_submit [get]
func (r *SubmitHandler) DeleteSubmit(ctx *gin.Context) {
	//获取请求数据
	idStr := ctx.Query("id")
	if idStr == "" {
		r.logger.Infof("invalid submit id")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid submit id"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		r.logger.Infof("invalid submit id")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid submit id"))
		return
	}

	//删除提交信息
	err = r.submitService.DeleteSubmit(id)
	if err != nil {
		r.logger.Errorf("delete submit error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "delete submit error"))
		return
	}

	//返回结果
	ctx.Set("data", "delete submit success")
}
