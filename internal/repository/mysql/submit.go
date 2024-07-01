package mysql

import (
	"fmt"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (mysql *MysqlClient) CreateSubmit(submit *entity.Submit) error {
	err := createData[entity.Submit](mysql, constant.SUBMIT_TABLE, submit)
	if err != nil {
		return fmt.Errorf("repository layer: mysql, create submit, %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (mysql *MysqlClient) GetSubmitById(submitId int) (*entity.Submit, error) {
	data, err := getDataById[entity.Submit](mysql, constant.SUBMIT_TABLE, submitId)
	if err != nil {
		return nil, fmt.Errorf("repository layer: mysql, get submit by id, %w %+v", constant.ErrInternal, err)
	}
	if data == nil {
		return nil, fmt.Errorf("repository layer: mysql, get submit by id, %w", constant.ErrNotFound)
	}
	return data, nil
}

func (mysql *MysqlClient) GetSubmitList(page, pageSize int) ([]*entity.Submit, error) {
	data, err := getDataList[entity.Submit](mysql, constant.SUBMIT_TABLE, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("repository layer: mysql, get submit list, %w %+v", constant.ErrInternal, err)
	}

	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("repository layer: mysql, get submit list, %w", constant.ErrNotFound)
	}
	return data, nil
}

func (mysql *MysqlClient) UpdateSubmit(submit *entity.Submit) error {
	err := updateDataById[entity.Submit](mysql, constant.SUBMIT_TABLE, submit.ID, submit)
	if err != nil {
		return fmt.Errorf("repository layer: mysql, update submit, %w %+v", constant.ErrInternal, err)
	}
	return nil
}
func (mysql *MysqlClient) DeleteSubmit(submitId int) error {
	err := deleteDataById(mysql, constant.SUBMIT_TABLE, submitId)
	if err != nil {
		return fmt.Errorf("repository layer: mysql, delete submit, %w %+v", constant.ErrInternal, err)
	}
	return nil
}
func (mysql *MysqlClient) GetRecentSubmit(lastUpdateTime string) ([]*entity.Submit, error) {
	submitList, err := getRecentData[entity.Submit](mysql, constant.SUBMIT_TABLE, lastUpdateTime)
	if err != nil {
		return nil, fmt.Errorf("repository layer: mysql, get recent submit, %w %+v", constant.ErrInternal, err)
	}
	if submitList == nil || len(submitList) == 0 {
		return nil, fmt.Errorf("repository layer: mysql, get recent submit, %w", constant.ErrNotFound)
	}
	return submitList, nil
}
