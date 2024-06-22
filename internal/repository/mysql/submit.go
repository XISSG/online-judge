package mysql

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (mysql *MysqlClient) CreateSubmit(submit *entity.Submit) error {
	err := createData[entity.Submit](mysql, constant.SUBMIT_TABLE, submit)
	if err != nil {
		return err
	}
	return nil
}

func (mysql *MysqlClient) GetSubmitById(submitId int) *entity.Submit {
	return getDataById[entity.Submit](mysql, constant.SUBMIT_TABLE, submitId)
}

func (mysql *MysqlClient) GetSubmitList(page, pageSize int) (submitList []*entity.Submit) {
	return getDataList[entity.Submit](mysql, constant.SUBMIT_TABLE, page, pageSize)
}

func (mysql *MysqlClient) UpdateSubmit(submit *entity.Submit) error {
	return updateDataById[entity.Submit](mysql, constant.SUBMIT_TABLE, submit.ID, submit)
}
func (mysql *MysqlClient) DeleteSubmit(submitId int) error {
	err := deleteDataById(mysql, constant.SUBMIT_TABLE, submitId)
	if err != nil {
		return err
	}
	return nil
}
func (mysql *MysqlClient) GetRecentSubmit(lastUpdateTime string) []*entity.Submit {
	submitList, err := getRecentData[entity.Submit](mysql, constant.SUBMIT_TABLE, lastUpdateTime)
	if err != nil {
		return nil
	}
	return submitList
}
