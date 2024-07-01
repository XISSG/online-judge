package redis

import (
	"fmt"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (redis *RedisClient) CacheQuestionList(questionList []*entity.Question) error {
	var ids []int
	for _, q := range questionList {
		ids = append(ids, q.ID)
	}
	err := cacheOrUpdateData[entity.Question](redis, constant.QUESTION_TABLE, ids, ids, questionList)
	if err != nil {
		return fmt.Errorf("repository layer: redis, cache question list: %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (redis *RedisClient) GetQuestionList(page, pageSize int) (questionList []*entity.Question, err error) {
	data := getDataList[entity.Question](redis, constant.QUESTION_TABLE, page, pageSize)
	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("repository layer: redis, get question list: %w", constant.ErrNotFound)
	}
	return data, nil
}

func (redis *RedisClient) DeleteQuestionById(questionId int) error {
	err := deleteDataById(redis, constant.QUESTION_TABLE, questionId)
	if err != nil {
		return fmt.Errorf("repository layer: redis, delete question by id: %w %+v", constant.ErrInternal, err)
	}
	return nil
}
