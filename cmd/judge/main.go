package main

import (
	"github.com/xissg/online-judge/internal/config"
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

	dockerClient := docker.NewDockerClient()
	mysqlClient := mysql.NewMysqlClient(appConfig.Database)
	redisClient := redis.NewRedisClient(appConfig.Redis)
	esClient := elastic.NewElasticSearchClient(appConfig.Elasticsearch)
	aiClient := ai.NewAIClient(appConfig.AI)
	rabbitMqClient := rabbitmq.NewRabbitMQClient(appConfig.RabbitMQ)

	questionSvc := service.NewQuestionService(mysqlClient, esClient, redisClient)
	submitSvc := service.NewSubmitService(mysqlClient, esClient, redisClient)
	aiSvc := service.NewAIService(aiClient, appConfig.AI)
	rabbitMqSvc := service.NewRabbitMqService(rabbitMqClient)

	judgeSvc := service.NewJudgeService(dockerClient, questionSvc, submitSvc, aiSvc, rabbitMqSvc)

	rabbitMqSvc.Consume(judgeSvc.Run)
}
