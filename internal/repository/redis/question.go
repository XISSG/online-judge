package redis

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (redis *RedisClient) CacheQuestionList(questionList []*entity.Question) error {
	for i := range questionList {
		err := cacheOrUpdateData[entity.Question](redis, constant.QUESTION_TABLE, questionList[i].ID, questionList[i].ID, questionList[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (redis *RedisClient) GetQuestionList(page, pageSize int) (questionList []*entity.Question) {
	return getDataList[entity.Question](redis, constant.QUESTION_TABLE, page, pageSize)
}

func (redis *RedisClient) DeleteQuestionById(questionId int) error {
	return deleteDataById(redis, constant.QUESTION_TABLE, questionId)
}
