package redis

import (
	"fmt"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/model/entity"
	"testing"
)

func TestRedisClient(t *testing.T) {
	appConfig := config.LoadConfig()
	client := NewRedisClient(appConfig.Redis)
	dataList := []*entity.Question{
		{
			ID:    1,
			Title: "test1",
		},
		{
			ID:      2,
			Title:   "test2",
			Content: "test",
		},
	}
	err := client.CacheQuestionList(dataList)
	if err != nil {
		panic(err)
	}
	questions := client.GetQuestionList(1, 10)
	fmt.Println(questions)
	err = client.DeleteQuestionById(1)
	if err != nil {
		panic(err)
	}
}
