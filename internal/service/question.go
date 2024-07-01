package service

import (
	"errors"
	"fmt"
	"github.com/xissg/online-judge/internal/model/entity"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/model/response"
	"github.com/xissg/online-judge/internal/repository/elastic"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/repository/redis"
	"github.com/xissg/online-judge/internal/utils"
)

type QuestionService interface {
	CreateQuestion(question *request.Question, userId int) error
	SearchQuestion(query string) ([]*response.Question, error)
	GetQuestionList(page, pageSize int) ([]*response.Question, error)
	UpdateQuestion(question *request.UpdateQuestion) error
	DeleteQuestion(id int) error
	UpdateQuestionAcceptNum(id int) error
	UpdateQuestionSubmitNum(id int) error
	GetQuestionById(id int) (*entity.Question, error)
}

type questionService struct {
	mysql *mysql.MysqlClient
	es    *elastic.ESClient
	redis *redis.RedisClient
}

func NewQuestionService(mysql *mysql.MysqlClient, es *elastic.ESClient, redis *redis.RedisClient) QuestionService {
	return &questionService{
		mysql: mysql,
		es:    es,
		redis: redis,
	}
}

func (q *questionService) CreateQuestion(question *request.Question, userId int) error {
	data := utils.ConvertQuestionEntity(question, userId)
	if data == nil {
		return errors.New("service layer: question, data marshalling error")
	}
	err := q.mysql.CreateQuestion(data)
	if err != nil {
		return err
	}
	resp := utils.ConvertQuestionResponse(data)
	err = q.es.IndexQuestion(resp)
	if err != nil {
		return fmt.Errorf("service layer: question -> %w", err)
	}
	return nil
}

func (q *questionService) SearchQuestion(query string) ([]*response.Question, error) {
	questions, err := q.es.SearchQuestions(query)
	if err != nil {
		return nil, fmt.Errorf("service layer: question -> %w", err)
	}
	return questions, nil
}

func (q *questionService) GetQuestionList(page, pageSize int) ([]*response.Question, error) {
	var questions []*entity.Question
	var err error

	questions, err = q.redis.GetQuestionList(page, pageSize)
	if questions == nil {
		questions, err = q.mysql.GetQuestionList(page, pageSize)
		if questions == nil {
			return nil, fmt.Errorf("service layer: question -> %w", err)
		}
		_ = q.redis.CacheQuestionList(questions)
	}

	var result []*response.Question

	for i := range questions {
		question := utils.ConvertQuestionResponse(questions[i])
		result = append(result, question)
	}
	return result, nil
}

func (q *questionService) UpdateQuestion(question *request.UpdateQuestion) error {
	questionEntity := utils.UpdateQuestionToQuestionEntity(question)
	err := q.mysql.UpdateQuestion(questionEntity)
	if err != nil {
		return fmt.Errorf("service layer: question -> %w", err)
	}
	questionResponse := utils.ConvertQuestionResponse(questionEntity)
	err = q.es.UpdateQuestion(questionResponse)
	err = q.redis.DeleteQuestionById(question.ID)
	if err != nil {
		return fmt.Errorf("service layer: question -> %w", err)
	}
	return nil
}

func (q *questionService) DeleteQuestion(id int) error {
	err := q.mysql.DeleteQuestion(id)
	err = q.es.DeleteQuestionById(id)
	err = q.redis.DeleteQuestionById(id)
	if err != nil {
		return fmt.Errorf("service layer: question -> %w", err)
	}
	return nil
}
func (q *questionService) UpdateQuestionAcceptNum(id int) error {
	questionEntity, err := q.mysql.GetQuestionById(id)
	if err != nil {
		return fmt.Errorf("service layer: question -> %w", err)
	}
	updateRequest := &request.UpdateQuestion{
		ID:        id,
		AcceptNum: questionEntity.AcceptNum + 1,
		SubmitNum: questionEntity.SubmitNum + 1,
	}
	err = q.UpdateQuestion(updateRequest)
	if err != nil {
		return fmt.Errorf("service layer: question -> %w", err)
	}
	return nil
}

func (q *questionService) UpdateQuestionSubmitNum(id int) error {
	questionEntity, err := q.mysql.GetQuestionById(id)
	if err != nil {
		return fmt.Errorf("service layer: question -> %w", err)
	}
	updateRequest := &request.UpdateQuestion{
		ID:        id,
		SubmitNum: questionEntity.SubmitNum + 1,
	}
	err = q.UpdateQuestion(updateRequest)
	if err != nil {
		return fmt.Errorf("service layer: question -> %w", err)
	}
	return nil
}

func (q *questionService) GetQuestionById(id int) (*entity.Question, error) {
	questionEntity, err := q.mysql.GetQuestionById(id)
	if err != nil {
		return nil, fmt.Errorf("service layer: question -> %w", err)
	}
	return questionEntity, nil
}
