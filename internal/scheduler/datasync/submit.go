package datasync

import (
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/repository/elastic"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/utils"
)

// 同步用户提交
type Submit struct {
	es             *elastic.ESClient
	mysql          *mysql.MysqlClient
	lastUpdateTime string
}

func NewSubmitSync(esConfig config.ElasticsearchConfig, mysqlConfig config.DatabaseConfig) DataSync {
	es := elastic.NewElasticSearchClient(esConfig)
	mysqlClient := mysql.NewMysqlClient(mysqlConfig)
	return &Submit{
		es:             es,
		mysql:          mysqlClient,
		lastUpdateTime: "",
	}
}

func (s *Submit) SyncData() {
	submits, err := s.mysql.GetRecentSubmit(s.lastUpdateTime)
	if err != nil {
		return
	}
	for _, submit := range submits {
		submitResponse := utils.ConvertSubmitResponse(submit)
		s.es.UpdateSubmit(submitResponse)
	}
	s.lastUpdateTime = submits[0].UpdateTime
}
