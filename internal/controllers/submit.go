package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/service"
	"strconv"
)

type SubmitHandler struct {
	submitService service.SubmitService
}

func NewSubmitHandler(submitService service.SubmitService) *SubmitHandler {
	return &SubmitHandler{
		submitService: submitService,
	}
}

func (r *SubmitHandler) RegisterRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	{
		admin.GET("/delete_submit", r.deleteSubmit)
	}
	submit := router.Group("/submit")
	{
		submit.POST("/create_submit", r.createSubmit)
		submit.GET("/get_submits", r.getSubmitList)
		submit.GET("/search_submits", r.searchSubmitList)
	}
}

func (r *SubmitHandler) createSubmit(ctx *gin.Context) {
	request := request.Submit{}
	userId := 1 //TODO: 从session中获取
	ctx.ShouldBindJSON(&request)
	r.submitService.CreateSubmit(&request, userId)
}

func (r *SubmitHandler) getSubmitList(ctx *gin.Context) {
	page := ctx.Query("page")
	pageSize := ctx.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	r.submitService.GetSubmitList(pageInt, pageSizeInt)
}

func (r *SubmitHandler) searchSubmitList(ctx *gin.Context) {
	keryword := ctx.Query("keyword")
	r.submitService.SearchSubmit(keryword)
}
func (r *SubmitHandler) deleteSubmit(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, _ := strconv.Atoi(idStr)
	r.submitService.DeleteSubmit(id)
}
