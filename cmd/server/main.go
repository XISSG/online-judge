package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/controllers"
	logger2 "github.com/xissg/online-judge/internal/logger"
	"github.com/xissg/online-judge/internal/middlewares"
	"github.com/xissg/online-judge/internal/repository/elastic"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/repository/rabbitmq"
	"github.com/xissg/online-judge/internal/repository/redis"
	"github.com/xissg/online-judge/internal/router"
	"github.com/xissg/online-judge/internal/service"
	"go.uber.org/ratelimit"
	"log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/xissg/online-judge/internal/docs"
)

func main() {

	//初始化配置
	appConfig := config.LoadConfig()
	limiter := ratelimit.New(5)
	logger, err := logger2.NewLogger(appConfig.Log)
	if err != nil {
		panic(err)
	}

	r := gin.New()

	//注册swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//注册middleware
	r.Use(middlewares.CORS())
	r.Use(middlewares.InvokeCount())
	r.Use(middlewares.InvokeLimit(limiter))
	r.Use(middlewares.ResponseMiddleware(logger))
	r.Use(middlewares.RecoveryMiddleware(logger))

	//初始化repository
	mysqlClient := mysql.NewMysqlClient(appConfig.Database)
	redisClient := redis.NewRedisClient(appConfig.Redis)
	esClient := elastic.NewElasticSearchClient(appConfig.Elasticsearch)
	rabbitMqClient := rabbitmq.NewRabbitMQClient(appConfig.RabbitMQ)

	//初始化service
	userService := service.NewUserService(mysqlClient)
	questionService := service.NewQuestionService(mysqlClient, esClient, redisClient)
	submitService := service.NewSubmitService(mysqlClient, esClient, redisClient)
	rabbitMqService := service.NewRabbitMqService(rabbitMqClient)

	//初始化handler
	userHandler := controllers.NewUserHandler(userService, logger)
	questionHandler := controllers.NewQuestionHandler(questionService, logger)
	submitHandler := controllers.NewSubmitHandler(submitService, rabbitMqService, logger)
	pictureHandler := controllers.NewPictureHandler(appConfig.Picture, logger)
	invokeHandler := controllers.NewInvokeHandler(redisClient, logger)

	//注册路由
	router.Router(r, userHandler, questionHandler, submitHandler, pictureHandler, invokeHandler, logger)
	//启动服务器
	if err := r.Run(fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
