package elasticsearch

import (
	"fmt"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/model/response"
	"strconv"
	"testing"
	"time"
)

func TestElasticSearch(t *testing.T) {
	appConfig := config.LoadConfig()
	client := NewElasticSearchClient(appConfig.Elasticsearch)
	var index *response.Question
	index = &response.Question{
		ID:          1,
		Title:       "test",
		Content:     "test content",
		//Tag:         ["test", "test"],
		//Answer:      ["{"test","test""],
		SubmitNum:   10,
		AcceptNum:   1,
		//JudgeCase:   ["{"test":"test"}"],
		//JudgeConfig: ["{"test":"test"}"],
		ThumNum:     10,
		UserId:      1,
		CreateTime:  time.Now().Format(time.RFC3339Nano),
		UpdateTime:  time.Now().Format(time.RFC3339Nano),
	}

	fmt.Println(index.CreateTime)
	err := client.IndexQuestions(index)
	if err != nil {
		panic(err)
	}
	questions := client.SearchQuestions("test")
	fmt.Println(questions)
	res := client.GetQuestionById(strconv.Itoa(1))
	fmt.Println(*res)
	err = client.DeleteQuestionById(strconv.Itoa(1))
	if err != nil {
		panic(err)
	}
}
