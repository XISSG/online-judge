package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/online-judge/internal/controllers"
	"github.com/xissg/online-judge/internal/middlewares"
	"go.uber.org/zap"
)

func Router(
	router *gin.Engine,
	userHandler *controllers.UserHandler,
	questionHandler *controllers.QuestionHandler,
	submitHandler *controllers.SubmitHandler,
	pictureHandler *controllers.PictureHandler,
	invokeHandler *controllers.InvokeHandler,
	logger *zap.SugaredLogger,
) {

	router.POST("/user/register", userHandler.Register)
	router.POST("/user/login", userHandler.Login)

	user := router.Group("/user")
	user.Use(middlewares.AuthLogin(logger))
	{
		user.GET("/question/get_questions", questionHandler.GetQuestionList)
		user.GET("/question/search_questions", questionHandler.SearchQuestionList)
		user.POST("/submit/create_submit", submitHandler.CreateSubmit)
		user.GET("/submit/get_submits", submitHandler.GetSubmitList)
		user.GET("/submit/search_submits", submitHandler.SearchSubmitList)
		user.GET("/picture/avatar", pictureHandler.GetAvatar)
	}

	admin := router.Group("/admin")
	admin.Use(middlewares.AuthAdmin(logger))
	{
		admin.GET("/user/get_users", userHandler.GetUserList)
		admin.POST("/user/update_user", userHandler.UpdateUser)
		admin.GET("/user/delete_user", userHandler.DeleteUser)
		admin.GET("/user/ban_user", userHandler.BanUser)
		admin.POST("/question/create_question", questionHandler.CreateQuestion)
		admin.POST("/question/update_question", questionHandler.UpdateQuestion)
		admin.GET("/question/delete_question", questionHandler.DeleteQuestion)
		admin.GET("/submit/delete_submit", submitHandler.DeleteSubmit)
		admin.POST("/invoke/count", invokeHandler.GetInvokeCount)
	}
}
