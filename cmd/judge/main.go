package main

import (
	"github.com/xissg/online-judge/internal/config"
	logger2 "github.com/xissg/online-judge/internal/logger"
	"github.com/xissg/online-judge/internal/repository/ai"
	"github.com/xissg/online-judge/internal/repository/docker"
	"github.com/xissg/online-judge/internal/repository/elastic"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/repository/rabbitmq"
	"github.com/xissg/online-judge/internal/repository/redis"
	"github.com/xissg/online-judge/internal/service"
)

func main() {
	appConfig := config.LoadConfig()

	logger, _ := logger2.NewLogger(appConfig.Log)

	dockerClient := docker.NewDockerClient()
	mysqlClient := mysql.NewMysqlClient(appConfig.Database)
	redisClient := redis.NewRedisClient(appConfig.Redis)
	esClient := elastic.NewElasticSearchClient(appConfig.Elasticsearch)

	aiClient := ai.NewAIClient(appConfig.AI)
	rabbitMqClient := rabbitmq.NewRabbitMQClient(appConfig.RabbitMQ)

	questionSvc := service.NewQuestionService(mysqlClient, esClient, redisClient)
	submitSvc := service.NewSubmitService(mysqlClient, esClient, redisClient)
	aiSvc := service.NewAIService(aiClient)
	rabbitMqSvc := service.NewRabbitMqService(rabbitMqClient)

	judgeSvc := service.NewJudgeService(dockerClient, questionSvc, submitSvc, aiSvc, logger)

	rabbitMqSvc.Consume(judgeSvc.Run)
}
