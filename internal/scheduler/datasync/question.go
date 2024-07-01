package datasync

import (
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/repository/elastic"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/utils"
)

// 同步题目
type Question struct {
	es             *elastic.ESClient
	mysql          *mysql.MysqlClient
	lastUpdateTime string
}

func NewQuestionSync(esConfig config.ElasticsearchConfig, mysqlConfig config.DatabaseConfig) DataSync {
	es := elastic.NewElasticSearchClient(esConfig)
	mysqlClient := mysql.NewMysqlClient(mysqlConfig)
	return &Question{
		es:             es,
		mysql:          mysqlClient,
		lastUpdateTime: "",
	}
}

func (q *Question) SyncData() {

	questions, err := q.mysql.GetRecentQuestion(q.lastUpdateTime)
	if err != nil {
		return
	}
	for _, question := range questions {
		questionResponse := utils.ConvertQuestionResponse(question)
		q.es.UpdateQuestion(questionResponse)
	}
	q.lastUpdateTime = questions[0].UpdateTime
}
