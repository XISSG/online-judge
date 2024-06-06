package mysql

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (mysql *MysqlClient) CreateSubmit(submit *entity.Submit) error {
	tx := mysql.Begin()
	tx.Model("submit").Create(submit)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}

func (mysql *MysqlClient) GetSubmitById(submitId int) *entity.Submit {
	var submit *entity.Submit
	mysql.Model("submit").First(&submit, submitId)
	return submit
}

func (mysql *MysqlClient) DeleteSubmit(submitId int) error {
	tx := mysql.Begin()
	tx.Model("submit").Where("id = ? ", submitId).Update("is_delete", constant.DELETED)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	return nil
}