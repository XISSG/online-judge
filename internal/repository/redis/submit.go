package redis

import (
	"fmt"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (redis *RedisClient) CacheSubmitList(submitList []*entity.Submit) error {
	var ids []int
	for _, q := range submitList {
		ids = append(ids, q.ID)
	}
	err := cacheOrUpdateData[entity.Submit](redis, constant.QUESTION_TABLE, ids, ids, submitList)
	if err != nil {
		return fmt.Errorf("repository layer: redis, cache submit list: %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (redis *RedisClient) GetSubmitList(page, pageSize int) (submitList []*entity.Submit, err error) {
	data := getDataList[entity.Submit](redis, constant.SUBMIT_TABLE, page, pageSize)
	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("repository layer: redis, get submit list: %w", constant.ErrNotFound)
	}
	return data, nil
}

func (redis *RedisClient) DeleteSubmitById(submitId int) error {
	err := deleteDataById(redis, constant.SUBMIT_TABLE, submitId)
	if err != nil {
		return fmt.Errorf("repository layer: redis, delete submit by id: %w %+v", constant.ErrInternal, err)
	}
	return nil
}
