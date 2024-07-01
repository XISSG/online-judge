package mysql

import (
	"fmt"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (mysql *MysqlClient) CreateQuestion(question *entity.Question) error {
	err := createData[entity.Question](mysql, constant.QUESTION_TABLE, question)
	if err != nil {
		return fmt.Errorf("repository layer: mysql, create question, %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (mysql *MysqlClient) GetQuestionById(questionId int) (*entity.Question, error) {
	data, err := getDataById[entity.Question](mysql, constant.QUESTION_TABLE, questionId)
	if err != nil {
		return nil, fmt.Errorf("repository layer: mysql, get question by id, %w %+v", constant.ErrInternal, err)
	}
	if data == nil {
		return nil, fmt.Errorf("repository layer: mysql, get question by id, %w", constant.ErrNotFound)
	}
	return data, nil
}

func (mysql *MysqlClient) GetQuestionList(page, pageSize int) ([]*entity.Question, error) {
	data, err := getDataList[entity.Question](mysql, constant.QUESTION_TABLE, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("repository layer: mysql, get question list, %w %+v", constant.ErrInternal, err)
	}

	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("repository layer: mysql, get question list, %w", constant.ErrNotFound)
	}
	return data, nil
}
func (mysql *MysqlClient) UpdateQuestion(question *entity.Question) error {
	err := updateDataById[entity.Question](mysql, constant.QUESTION_TABLE, question.ID, question)
	if err != nil {
		return fmt.Errorf("repository layer: mysql, update question, %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (mysql *MysqlClient) DeleteQuestion(questionId int) error {
	err := deleteDataById(mysql, constant.QUESTION_TABLE, questionId)
	if err != nil {
		return fmt.Errorf("repository layer: mysql, delete question, %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (mysql *MysqlClient) GetRecentQuestion(lastUpdateTime string) ([]*entity.Question, error) {
	questionList, err := getRecentData[entity.Question](mysql, constant.QUESTION_TABLE, lastUpdateTime)
	if err != nil {
		return nil, fmt.Errorf("repository layer: mysql, get recent question, %w %+v", constant.ErrInternal, err)
	}
	if questionList == nil || len(questionList) == 0 {
		return nil, fmt.Errorf("repository layer: mysql, get recent question, %w", constant.ErrNotFound)
	}
	return questionList, nil
}
