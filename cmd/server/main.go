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
	"github.com/xissg/online-judge/internal/service"
	"go.uber.org/ratelimit"
	"log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/xissg/online-judge/internal/docs"
)

func main() {

	appConfig := config.LoadConfig()
	limiter := ratelimit.New(5)
	logger, err := logger2.NewLogger(appConfig.Log)
	if err != nil {
		panic(err)
	}

	r := gin.New()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(middlewares.CORS())
	r.Use(middlewares.InvokeCount())
	r.Use(middlewares.InvokeLimit(limiter))
	r.Use(middlewares.ResponseMiddleware(logger))
	r.Use(middlewares.RecoveryMiddleware(logger))

	mysqlClient := mysql.NewMysqlClient(appConfig.Database)
	redisClient := redis.NewRedisClient(appConfig.Redis)
	esClient := elastic.NewElasticSearchClient(appConfig.Elasticsearch)
	rabbitMqClient := rabbitmq.NewRabbitMQClient(appConfig.RabbitMQ)

	userService := service.NewUserService(mysqlClient)
	questionService := service.NewQuestionService(mysqlClient, esClient, redisClient)
	submitService := service.NewSubmitService(mysqlClient, esClient, redisClient)
	rabbitMqService := service.NewRabbitMqService(rabbitMqClient)

	userHandler := controllers.NewUserHandler(userService, logger)
	userHandler.RegisterRoutes(r)

	questionHandler := controllers.NewQuestionHandler(questionService, logger)
	questionHandler.RegisterRoutes(r)

	submitHandler := controllers.NewSubmitHandler(submitService, rabbitMqService, logger)
	submitHandler.RegisterRoutes(r)

	invokeHandler := controllers.NewInvokeHandler(redisClient, logger)
	invokeHandler.RegisterRoutes(r)

	if err := r.Run(fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
