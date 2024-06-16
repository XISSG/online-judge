package online_judge

import (
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/repository/elastic"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/utils"
	"testing"
)

func Test(t *testing.T) {
	appConfig := config.LoadConfig()

	mysqlClient := mysql.NewMysqlClient(appConfig.Database)

	esClient := elastic.NewElasticSearchClient(appConfig.Elasticsearch)

	question := mysqlClient.GetQuestionList(1, 10)
	for i := range question {
		questionResponse := utils.ConvertQuestionResponse(question[i])
		esClient.IndexQuestion(questionResponse)
	}
}
