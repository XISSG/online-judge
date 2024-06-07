package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/service"
	"strconv"
)

type QuestionHandler struct {
	questionService service.QuestionService
}

func NewQuestionHandler(questionService service.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
	}
}

func (r *QuestionHandler) RegisterRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	{
		admin.POST("/create_question", r.createQuestion)
		admin.POST("/update_question", r.updateQuestion)
		admin.GET("/delete_question", r.deleteQuestion)
	}
	question := router.Group("/question")
	{
		question.GET("/get_questions", r.getQuestionList)
		question.GET("/search_questions", r.searchQuestionList)
	}
}

func (r *QuestionHandler) createQuestion(ctx *gin.Context) {
	request := request.Question{}
	ctx.ShouldBindJSON(&request)
	r.questionService.CreateQuestion(&request)
}

func (r *QuestionHandler) updateQuestion(ctx *gin.Context) {
	request := request.UpdateQuestion{}
	ctx.ShouldBindJSON(&request)
	r.questionService.UpdateQuestion(&request)
}

func (r *QuestionHandler) getQuestionList(ctx *gin.Context) {
	page := ctx.Query("page")
	pageSize := ctx.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	r.questionService.GetQuestionList(pageInt, pageSizeInt)
}

func (r *QuestionHandler) searchQuestionList(ctx *gin.Context) {
	keryword := ctx.Query("keyword")
	r.questionService.SearchQuestion(keryword)
}
func (r *QuestionHandler) deleteQuestion(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, _ := strconv.Atoi(idStr)
	r.questionService.DeleteQuestion(id)
}
