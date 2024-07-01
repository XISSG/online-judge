package service

import (
	"errors"
	"fmt"
	"github.com/xissg/online-judge/internal/model/entity"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/model/response"
	"github.com/xissg/online-judge/internal/repository/elastic"
	"github.com/xissg/online-judge/internal/repository/mysql"
	rdb "github.com/xissg/online-judge/internal/repository/redis"
	"github.com/xissg/online-judge/internal/utils"
)

type SubmitService interface {
	CreateSubmit(submit *request.Submit, userId int) (int, error)
	SearchSubmit(query string) ([]*response.Submit, error)
	GetSubmitList(page, pageSize int) ([]*response.Submit, error)
	UpdateSubmit(submit *request.UpdateSubmit) error
	DeleteSubmit(id int) error
	GetSubmitById(id int) (*entity.Submit, error)
}

type submitService struct {
	mysql *mysql.MysqlClient
	es    *elastic.ESClient
	redis *rdb.RedisClient
}

func NewSubmitService(mysql *mysql.MysqlClient, es *elastic.ESClient, redis *rdb.RedisClient) SubmitService {
	return &submitService{
		mysql: mysql,
		es:    es,
		redis: redis,
	}
}

func (q *submitService) CreateSubmit(submit *request.Submit, userId int) (int, error) {
	data := utils.ConvertSubmitEntity(submit, userId)
	if data == nil {
		return 0, errors.New("service layer: submit, data marshalling error")
	}

	err := q.mysql.CreateSubmit(data)
	resp := utils.ConvertSubmitResponse(data)
	err = q.es.IndexSubmit(resp)
	if err != nil {
		return 0, fmt.Errorf("service layer: submit -> %w", err)
	}

	return data.ID, nil
}

func (q *submitService) SearchSubmit(query string) ([]*response.Submit, error) {
	submits, err := q.es.SearchSubmits(query)
	if err != nil {
		return nil, fmt.Errorf("service layer: submit -> %w", err)
	}
	return submits, nil
}

func (q *submitService) GetSubmitList(page, pageSize int) ([]*response.Submit, error) {
	var submits []*entity.Submit
	var err error

	submits, err = q.redis.GetSubmitList(page, pageSize)
	if submits == nil {
		submits, err = q.mysql.GetSubmitList(page, pageSize)
		if err != nil {
			return nil, fmt.Errorf("service layer: submit -> %w", err)
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

func (q *submitService) UpdateSubmit(submit *request.UpdateSubmit) error {
	submitEntity := utils.UpdateSubmitToSubmitEntity(submit)
	err := q.mysql.UpdateSubmit(submitEntity)
	if err != nil {
		return fmt.Errorf("service layer: submit -> %w", err)
	}
	submitResponse := utils.ConvertSubmitResponse(submitEntity)
	err = q.es.UpdateSubmit(submitResponse)
	err = q.redis.DeleteSubmitById(submit.ID)
	if err != nil {
		return fmt.Errorf("service layer: submit -> %w", err)
	}
	return nil
}

func (q *submitService) DeleteSubmit(id int) error {
	err := q.mysql.DeleteSubmit(id)
	err = q.es.DeleteSubmitById(id)
	err = q.redis.DeleteSubmitById(id)
	if err != nil {
		return fmt.Errorf("service layer: submit -> %w", err)
	}
	return nil
}

func (q *submitService) GetSubmitById(id int) (*entity.Submit, error) {
	res, err := q.mysql.GetSubmitById(id)
	if err != nil {
		return nil, fmt.Errorf("service layer: submit -> %w", err)
	}
	return res, nil
}
