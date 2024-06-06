package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (r *UserHandler) RegisterRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	{
		admin.GET("/get_users", r.getUserList)
		admin.POST("/update_user", r.updateUser)
		admin.POST("/delete_user", r.deleteUser)
	}
	user := router.Group("/user")
	{
		user.POST("/register", r.register)
		user.POST("/login", r.login)

	}
}

func (r *UserHandler) register(ctx *gin.Context) {
	user := request.User{}
	ctx.ShouldBindJSON(&user)
	if user.AvatarURL == "" {
		//TODO:用户头像生成,提供一个本网站开放的头像网站,本网站提供一个图像接口，该字段为空时调用该接口，将其返回值的url填充到该字段
	}
	r.userService.CreateUser(&user)
	ctx.JSON(200, user)
}

func (r *UserHandler) login(ctx *gin.Context) {
	loginRequest := request.Login{}
	ctx.ShouldBindJSON(&loginRequest)
	ok := r.userService.CheckUser(loginRequest.UserName, loginRequest.UserPassword)
	if !ok {
		ctx.JSON(400, "Login failed")
		return
	}
	//TODO:生成jwt token
	ctx.JSON(200, "success")
}

func (r *UserHandler) getUserList(ctx *gin.Context) {
	id := ctx.Query("id")
	name := ctx.Query("name")
	if id != "" {
		user := r.userService.GetUserById(id)
		ctx.JSON(200, user)
		return
	}
	if name != "" {
		users := r.userService.GetUserByUsername(name)
		ctx.JSON(200, users)
		return
	}
	ctx.JSON(200, nil)
}

func (r *UserHandler) updateUser(ctx *gin.Context) {
	updateRequest := request.UpdateUser{}
	err := ctx.ShouldBindJSON(&updateRequest)
	if err != nil {
		return
	}
	if updateRequest.Type == "password" {
		r.userService.UpdateUserPassword(updateRequest.Body.ID, updateRequest.Body.Data)
		ctx.JSON(200, "success")
		return
	}

	if updateRequest.Type == "avatar" {
		//TODO:更新头像一般为用户自己上传图片，服务器保存，需优化逻辑,本网站提供一个文件上传接口，本机调用然后更新该字段，rpc调用？
		//TODO:完成调用接口
		r.userService.UpdateUserAvatar(updateRequest.Body.ID, updateRequest.Body.Data)
		ctx.JSON(200, "success")
		return
	}
	ctx.JSON(405, "error data")
}

func (r *UserHandler) deleteUser(ctx *gin.Context) {
	id := ctx.Query("id")
	r.userService.DeleteUserById(id)
	ctx.JSON(200, "success")
}
