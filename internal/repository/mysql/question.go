package mysql

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (mysql *MysqlClient) CreateQuestion(question *entity.Question) error {
	tx := mysql.Begin()
	tx.Model("question").Create(question)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}

func (mysql *MysqlClient) GetQuestionById(questionId int) *entity.Question {
	var question *entity.Question
	mysql.Model("question").First(&question, questionId)
	return question
}

func (mysql *MysqlClient) UpdateQuestion(question *entity.Question) error {
	tx := mysql.Begin()
	tx.Model("question").Where("id = ? ", question.ID).Updates(*question)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	return nil
}

func (mysql *MysqlClient) DeleteQuestion(questionId int) error {
	tx := mysql.Begin()
	tx.Model("question").Where("id = ? ", questionId).Update("is_delete", constant.DELETED)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	return nil
}
