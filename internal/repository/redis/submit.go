package redis

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (redis *RedisClient) CacheSubmitList(submitList []*entity.Submit) error {
	for i := range submitList {
		err := cacheOrUpdateData[entity.Submit](redis, constant.SUBMIT_TABLE, submitList[i].ID, submitList[i].ID, submitList[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (redis *RedisClient) GetSubmitList(page, pageSize int) (submitList []*entity.Submit) {
	return getDataList[entity.Submit](redis, constant.SUBMIT_TABLE, page, pageSize)
}

func (redis *RedisClient) DeleteSubmitById(submitId int) error {
	return deleteDataById(redis, constant.SUBMIT_TABLE, submitId)
}
