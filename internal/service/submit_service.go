package service

import (
	"errors"
	"github.com/xissg/online-judge/internal/model/entity"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/model/response"
	"github.com/xissg/online-judge/internal/repository/elasticsearch"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/repository/redis"
	"github.com/xissg/online-judge/internal/utils"
)

type SubmitService interface {
	CreateSubmit(submit *request.Submit, userId int) error
	SearchSubmit(query string) ([]*response.Submit, error)
	GetSubmitList(page, pageSize int) ([]*response.Submit, error)
	DeleteSubmit(id int) error
}

type submitService struct {
	mysql *mysql.MysqlClient
	es    *elasticsearch.ESClient
	redis *redis.RedisClient
}

func NewSubmitService(mysql *mysql.MysqlClient, es *elasticsearch.ESClient, redis *redis.RedisClient) SubmitService {
	return &submitService{
		mysql: mysql,
		es:    es,
		redis: redis,
	}
}

func (q *submitService) CreateSubmit(submit *request.Submit, userId int) error {
	data := utils.ConvertSubmitEntity(submit, userId)
	if data == nil {
		return errors.New("data marshalling error")
	}
	return q.mysql.CreateSubmit(data)
}

func (q *submitService) SearchSubmit(query string) ([]*response.Submit, error) {
	submits := q.es.SearchSubmits(query)
	if submits == nil {
		return nil, errors.New("not found query submit")
	}
	return submits, nil
}

func (q *submitService) GetSubmitList(page, pageSize int) ([]*response.Submit, error) {
	var submits []*entity.Submit

	submits = q.redis.GetSubmitList(page, pageSize)
	if submits == nil {
		submits = q.mysql.GetSubmitList(page, pageSize)
		if submits == nil {
			return nil, errors.New("not found query submit")
		}
		_ = q.redis.CacheSubmitList(submits)
	}

	var result []*response.Submit

	for i := range submits {
		submit := utils.ConvertSubmitResponse(submits[i])
		result = append(result, submit)
	}
	return result, nil
}

func (q *submitService) DeleteSubmit(id int) error {
	err := q.mysql.DeleteSubmit(id)
	err = q.es.DeleteSubmitById(id)
	err = q.redis.DeleteSubmitById(id)
	return err
}
