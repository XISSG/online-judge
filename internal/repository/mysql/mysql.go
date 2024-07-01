package mysql

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

type RESPONSE interface {
	entity.User | entity.Question | entity.Submit
}

func createData[T RESPONSE](mysql *MysqlClient, table string, t *T) error {
	tx := mysql.client.Begin()
	tx.Table(table).Create(&t)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}

func getDataById[T RESPONSE](mysql *MysqlClient, table string, id int) (*T, error) {
	var t T
	var err error
	tx := mysql.client.Table(table).Where("id = ?", id).First(&t)
	if tx.Error != nil {
		return nil, err
	}
	return &t, nil
}

func getDataList[T RESPONSE](mysql *MysqlClient, table string, page int, pageSize int) ([]*T, error) {
	var t []*T
	var err error
	tx := mysql.client.Table(table).Where("is_delete =?", constant.NOT_DELETED).Order("update_time desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&t)
	if tx.Error != nil {
		return nil, err
	}
	return t, nil
}

func updateDataById[T RESPONSE](mysql *MysqlClient, table string, id int, t *T) error {
	tx := mysql.client.Begin()
	tx.Table(table).Where("id =?", id).Updates(*t)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}

func deleteDataById(mysql *MysqlClient, table string, id int) error {
	tx := mysql.client.Begin()
	tx.Table(table).Where("id = ?", id).Update("is_delete", constant.DELETED)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}

func getRecentData[T RESPONSE](mysql *MysqlClient, table string, lastUpdateTime string) ([]*T, error) {
	var t []*T
	tx := mysql.client.Table(table).Where("is_delete =?", constant.NOT_DELETED).Where("update_time >?", lastUpdateTime).Order("update_time desc").Find(&t)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return t, nil
}
