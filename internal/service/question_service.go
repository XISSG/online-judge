package service

import (
	"errors"
	"github.com/xissg/online-judge/internal/model/entity"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/model/response"
	"github.com/xissg/online-judge/internal/repository/elastic"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/repository/redis"
	"github.com/xissg/online-judge/internal/utils"
)

type QuestionService interface {
	CreateQuestion(question *request.Question) error
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

func (q *questionService) CreateQuestion(question *request.Question) error {
	data := utils.ConvertQuestionEntity(question)
	if data == nil {
		return errors.New("data marshalling error")
	}
	return q.mysql.CreateQuestion(data)
}

func (q *questionService) SearchQuestion(query string) ([]*response.Question, error) {
	questions := q.es.SearchQuestions(query)
	if questions == nil {
		return nil, errors.New("not found query question")
	}
	return questions, nil
}

func (q *questionService) GetQuestionList(page, pageSize int) ([]*response.Question, error) {
	var questions []*entity.Question

	questions = q.redis.GetQuestionList(page, pageSize)
	if questions == nil {
		questions = q.mysql.GetQuestionList(page, pageSize)
		if questions == nil {
			return nil, errors.New("not found query question")
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
		return nil
	}
	questionResponse := utils.ConvertQuestionResponse(questionEntity)
	err = q.es.UpdateQuestion(questionResponse)
	err = q.redis.DeleteQuestionById(question.ID)
	if err != nil {
		return err
	}
	return nil
}

func (q *questionService) DeleteQuestion(id int) error {
	err := q.mysql.DeleteQuestion(id)
	err = q.es.DeleteQuestionById(id)
	err = q.redis.DeleteQuestionById(id)
	return err
}
func (q *questionService) UpdateQuestionAcceptNum(id int) error {
	questionEntity := q.mysql.GetQuestionById(id)
	updateRequest := &request.UpdateQuestion{
		ID:        id,
		AcceptNum: questionEntity.AcceptNum + 1,
		SubmitNum: questionEntity.SubmitNum + 1,
	}
	return q.UpdateQuestion(updateRequest)
}

func (q *questionService) UpdateQuestionSubmitNum(id int) error {
	questionEntity := q.mysql.GetQuestionById(id)
	updateRequest := &request.UpdateQuestion{
		ID:        id,
		SubmitNum: questionEntity.SubmitNum + 1,
	}
	return q.UpdateQuestion(updateRequest)
}

func (q *questionService) GetQuestionById(id int) (*entity.Question, error) {
	questionEntity := q.mysql.GetQuestionById(id)
	if questionEntity == nil {
		return nil, errors.New("not found question")
	}
	return questionEntity, nil
}
