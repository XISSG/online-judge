package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xissg/online-judge/internal/middlewares"
	"github.com/xissg/online-judge/internal/model/common"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/service"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type QuestionHandler struct {
	questionService service.QuestionService
	logger          *zap.SugaredLogger
	validator       *validator.Validate
}

func NewQuestionHandler(questionService service.QuestionService, logger *zap.SugaredLogger) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
		logger:          logger,
		validator:       validator.New(),
	}
}

// CreateQuestion
//
//	@Summary		Create question
//	@Description	Create question
//	@Tags			question
//	@Accept			json
//	@Produce		json
//	@Param			createRequest	body		request.Question		true	"create question"
//	@Success		200				{object}	middlewares.Response	"ok"
//	@Failure		400				{object}	middlewares.Response	"bad request"
//	@Failure		500				{object}	middlewares.Response	"Internal Server Error"
//	@Router			/admin/question/create_question [post]
func (r *QuestionHandler) CreateQuestion(ctx *gin.Context) {
	//获取请求数据
	req := request.Question{
		JudgeConfig: &common.Config{},
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		r.logger.Infof("question request error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "question request error"))
		return
	}

	//校验数据
	err = r.validator.Struct(req)
	if err != nil {
		r.logger.Infof("question validate error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "question validate error"))
		return
	}

	//创建题目
	userId, _ := ctx.Get("user")
	err = r.questionService.CreateQuestion(&req, userId.(int))
	if err != nil {
		r.logger.Errorf("create question error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "create question error"))
		return
	}

	//返回结果
	ctx.Set("data", "create question success")
}

// UpdateQuestion
//
//	@Summary		Update question
//	@Description	Update question
//	@Tags			question
//	@Accept			json
//	@Produce		json
//	@Param			updateRequest	body		request.UpdateQuestion	true	"update question"
//	@Success		200				{object}	middlewares.Response	"ok"
//	@Failure		400				{object}	middlewares.Response	"bad request"
//	@Failure		500				{object}	middlewares.Response	"Internal Server Error"
//	@Router			/admin/question/update_question [post]
func (r *QuestionHandler) UpdateQuestion(ctx *gin.Context) {
	//获取请求数据
	req := request.UpdateQuestion{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		r.logger.Infof("question request error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "question request error"))
		return
	}

	//校验数据
	err = r.validator.Struct(req)
	if err != nil {
		r.logger.Infof("question validate error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "question validate error"))
		return
	}

	//更新题目
	err = r.questionService.UpdateQuestion(&req)
	if err != nil {
		r.logger.Errorf("update question error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "update question error"))
		return
	}

	//返回结果
	ctx.Set("data", "update question success")
}

// GetQuestionList
//
//	@Summary		get question list
//	@Description	get question list
//	@Tags			question
//	@Accept			json
//	@Produce		json
//	@Param			page		query		string					true	"page number"
//	@Param			page_size	query		string					true	"page size"
//	@Success		200			{object}	middlewares.Response	"ok"
//	@Failure		400			{object}	middlewares.Response	"bad request"
//	@Failure		500			{object}	middlewares.Response	"Internal Server Error"
//	@Router			/user/question/get_questions [get]
func (r *QuestionHandler) GetQuestionList(ctx *gin.Context) {
	//获取请求数据
	page := ctx.Query("page")
	pageSize := ctx.Query("page_size")

	//校验数据
	if page == "" || pageSize == "" {
		r.logger.Infof("invalid question page or page size")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid question page or page size"))
		return
	}

	pageInt, err := strconv.Atoi(page)
	pageSizeInt, err := strconv.Atoi(pageSize)

	if err != nil {
		r.logger.Infof("invalid question page or page size")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid question page or page size"))
		return
	}

	//获取题目列表
	questions, err := r.questionService.GetQuestionList(pageInt, pageSizeInt)
	if err != nil {
		r.logger.Errorf("get question list error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "get question list error"))
		return
	}

	//返回结果
	ctx.Set("data", questions)
}

// SearchQuestionList
//
//	@Summary		Search question
//	@Description	Search question
//	@Tags			question
//	@Accept			json
//	@Produce		json
//	@Param			keyword	query		string					true	"search question"
//	@Success		200		{object}	middlewares.Response	"ok"
//	@Failure		400		{object}	middlewares.Response	"bad request"
//	@Failure		500		{object}	middlewares.Response	"Internal Server Error"
//	@Router			/user/question/search_questions [get]
func (r *QuestionHandler) SearchQuestionList(ctx *gin.Context) {
	//获取请求数据
	keyword := ctx.Query("keyword")

	//校验数据
	if keyword == "" {
		r.logger.Infof("invalid search keyword")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid search keyword"))
	}

	//获取题目列表
	question, err := r.questionService.SearchQuestion(keyword)
	if err != nil {
		r.logger.Errorf("search question error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "search question error"))
		return
	}

	//返回结果
	ctx.Set("data", question)
}

// DeleteQuestion
//
//	@Summary		Delete question
//	@Description	Delete question
//	@Tags			question
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string					true	"delete question"
//	@Success		200	{object}	middlewares.Response	"ok"
//	@Failure		400	{object}	middlewares.Response	"bad request"
//	@Failure		500	{object}	middlewares.Response	"Internal Server Error"
//	@Router			/admin/question/delete_question [get]
func (r *QuestionHandler) DeleteQuestion(ctx *gin.Context) {
	//获取请求数据
	idStr := ctx.Query("id")

	//校验数据
	if idStr == "" {
		r.logger.Infof("invalid question id")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid question id"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		r.logger.Infof("invalid question id")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid question id"))
		return
	}

	//删除题目
	err = r.questionService.DeleteQuestion(id)
	if err != nil {
		r.logger.Errorf("delete question error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "delete question error"))
		return
	}

	//返回结果
	ctx.Set("data", "delete question success")
}
