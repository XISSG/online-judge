package mysql

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (mysql *MysqlClient) CreateQuestion(question *entity.Question) error {
	err := createData[entity.Question](mysql, constant.QUESTION_TABLE, question)
	if err != nil {
		return err
	}
	return nil
}

func (mysql *MysqlClient) GetQuestionById(questionId int) *entity.Question {
	question := getDataById[entity.Question](mysql, constant.QUESTION_TABLE, questionId)
	return question
}

func (mysql *MysqlClient) GetQuestionList(page, pageSize int) (questionList []*entity.Question) {
	return getDataList[entity.Question](mysql, constant.QUESTION_TABLE, page, pageSize)
}
func (mysql *MysqlClient) UpdateQuestion(question *entity.Question) error {
	err := updateDataById[entity.Question](mysql, constant.QUESTION_TABLE, question.ID, question)
	if err != nil {
		return err
	}
	return nil
}

func (mysql *MysqlClient) DeleteQuestion(questionId int) error {
	err := deleteDataById(mysql, constant.QUESTION_TABLE, questionId)
	if err != nil {
		return err
	}
	return nil
}

func (mysql *MysqlClient) GetRecentQuestion(lastUpdateTime string) []*entity.Question {
	questionList, err := getRecentData[entity.Question](mysql, constant.QUESTION_TABLE, lastUpdateTime)
	if err != nil {
		return nil
	}
	return questionList
}
