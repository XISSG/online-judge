package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/middlewares"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/service"
	"github.com/xissg/online-judge/internal/utils"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService service.UserService
	logger      *zap.SugaredLogger
	validator   *validator.Validate
}

func NewUserHandler(userService service.UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
		validator:   validator.New(),
	}
}

func (r *UserHandler) RegisterRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	admin.Use(middlewares.AuthAdmin())
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

// register
//
//	@Summary		User registration
//	@Description	User registration
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			createRequest	body		request.User			true	"create user"
//	@Success		200				{object}	middlewares.Response	"ok"
//	@Failure		400				{object}	middlewares.Response	"bad request"
//	@Failure		500				{object}	middlewares.Response	"Internal Server Error"
//	@Router			/user/register [post]
func (r *UserHandler) register(ctx *gin.Context) {
	//获取请求数据
	user := request.User{}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		r.logger.Infof("register error: %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "register error"))
		return
	}

	//校验数据
	err = r.validator.Struct(user)
	if user.AvatarURL != "" {
		err := r.validator.Var(user.AvatarURL, "url")
		if err != nil {
			r.logger.Infof("illegal avatar URL: %v", err)
			ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "illegal avatar url"))
		}
	}

	//判断用户名是否已存在
	exist := r.userService.CheckUserExists(user.UserName)
	if exist {
		r.logger.Errorf("Register user name is already exist")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "register user name is already exist"))
		return
	}

	//创建用户
	err = r.userService.CreateUser(&user)
	if err != nil {
		r.logger.Errorf("create user error: %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "create user error"))
		return
	}

	//返回结果
	ctx.Set("data", "register success")
}

// login
//
//	@Summary		User login
//	@Description	User login
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			loginRequest	body		request.Login			true	"user login"
//	@Success		200				{object}	middlewares.Response	"ok"
//	@Failure		400				{object}	middlewares.Response	"bad request"
//	@Failure		500				{object}	middlewares.Response	"Internal Server Error"
//	@Router			/user/login [post]
func (r *UserHandler) login(ctx *gin.Context) {
	//获取请求数据
	loginRequest := request.Login{}
	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		r.logger.Infof("login error: %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "login error"))
		return
	}

	//校验数据
	err = r.validator.Struct(loginRequest)
	if err != nil {
		r.logger.Infof("login error: %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "login error"))
		return
	}

	//判断用户名和密码是否匹配
	ok, userId, userRole := r.userService.CheckUserNameAndPasswd(loginRequest.UserName, loginRequest.UserPassword)
	if !ok {
		ctx.JSON(400, "login failed")
		return
	}

	//生成jwt唯一标识
	jwtConfig := config.LoadConfig().Jwt
	tokenString, err := utils.Generate(userId, userRole, jwtConfig)
	if err != nil {
		r.logger.Errorf("generate token error: %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "login error"))
		return
	}

	//返回结果
	ctx.Set("data", tokenString)
}

// getUserList
//
//	@Summary		Get user list
//	@Description	Get user list
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id		query		string				true	"get user list by id"
//	@Param			name	query		string				true	"get user list by name"
//	@Success		200		{object}	middlewares.Response	"ok"
//	@Failure		400		{object}	middlewares.Response	"bad request"
//	@Failure		500		{object}	middlewares.Response	"Internal Server Error"
//	@Router			/admin/get_users [get]
func (r *UserHandler) getUserList(ctx *gin.Context) {
	//获取请求数据
	idStr := ctx.Query("id")
	name := ctx.Query("name")

	//判断请求类型
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			r.logger.Infof("Invalid id: %v", err)
			ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid id"))
			return
		}
		//获取数据并返回
		user := r.userService.GetUserById(id)
		ctx.Set("data", user)
		return
	}

	if name != "" {
		//获取数据并返回
		users := r.userService.GetUserListByUsername(name)
		ctx.Set("data", users)
		return
	}
}

// updateUser
//
//	@Summary		Update user
//	@Description	Update user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			updateRequest	body		request.UpdateUser		true	"update user"
//	@Success		200				{object}	middlewares.Response	"ok"
//	@Failure		400				{object}	middlewares.Response	"bad request"
//	@Failure		500				{object}	middlewares.Response	"Internal Server Error"
//	@Router			/admin/update_user [post]
func (r *UserHandler) updateUser(ctx *gin.Context) {
	//获取请求数据
	updateRequest := request.UpdateUser{}
	err := ctx.ShouldBindJSON(&updateRequest)
	if err != nil {
		r.logger.Infof("Error updating user request error %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Error updating user request"))
		return
	}

	//判断数据类型
	if updateRequest.Type == "password" {
		//更新结果并返回
		err := r.userService.UpdateUserPassword(updateRequest.Body.ID, updateRequest.Body.Data)
		if err != nil {
			r.logger.Errorf("Error updating user password error %v", err)
			ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Error updating user password"))
			return
		}
		ctx.Set("data", "success")
		return
	}

	if updateRequest.Type == "avatar" {
		//更新结果并返回
		err := r.userService.UpdateUserAvatar(updateRequest.Body.ID, updateRequest.Body.Data)
		if err != nil {
			r.logger.Errorf("Error updating user avatar error %v", err)
			ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "Error updating user avatar"))
			return
		}
		ctx.Set("data", "success")
		return
	}
	//数据类型不合法
	r.logger.Infof("invalid update request")
	ctx.JSON(http.StatusBadRequest, "update data type must be password or avatar")
}

// deleteUser
//
//	@Summary		User registration
//	@Description	User registration
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string					true	"delete user by id"
//	@Success		200	{object}	middlewares.Response	"ok"
//	@Failure		400	{object}	middlewares.Response	"bad request"
//	@Failure		500	{object}	middlewares.Response	"Internal Server Error"
//	@Router			/admin/delete_user [get]
func (r *UserHandler) deleteUser(ctx *gin.Context) {
	//获取请求数据
	idStr := ctx.Query("id")
	if idStr == "" {
		r.logger.Infof("invalid id")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid id"))
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		r.logger.Infof("invalid id")
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "invalid id"))
		return
	}

	//删除用户
	err = r.userService.DeleteUserById(id)
	if err != nil {
		r.logger.Errorf("delete user error: %v", err)
		ctx.JSON(http.StatusBadRequest, middlewares.ErrorResponse(http.StatusBadRequest, "delete user error"))
		return
	}

	//返回结果
	ctx.Set("data", "success")
}
