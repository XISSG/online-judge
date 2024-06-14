package redis

import (
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
		return err
	}
	return nil
}

func (redis *RedisClient) GetSubmitList(page, pageSize int) (submitList []*entity.Submit) {
	return getDataList[entity.Submit](redis, constant.SUBMIT_TABLE, page, pageSize)
}

func (redis *RedisClient) DeleteSubmitById(submitId int) error {
	return deleteDataById(redis, constant.SUBMIT_TABLE, submitId)
}
