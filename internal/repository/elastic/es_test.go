package elastic

import (
	"encoding/json"
	"fmt"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/model/response"
	"testing"
	"time"
	"unsafe"
)

func TestQuestionSearch(t *testing.T) {
	appConfig := config.LoadConfig()
	client := NewElasticSearchClient(appConfig.Elasticsearch)
	var index *response.Question
	index = &response.Question{
		ID:        1,
		Title:     "test",
		Content:   "test content",
		Tag:       []string{"test", "test"},
		Answer:    []string{"test", "test"},
		SubmitNum: 10,
		AcceptNum: 1,
		JudgeCase: []string{"test", "test"},
		JudgeConfig: []response.Config{{
			TimeLimit:   time.Hour,
			MemoryLimit: 100,
		},
		},
		UserId:     1,
		CreateTime: time.Now().Format(time.RFC3339Nano),
		UpdateTime: time.Now().Format(time.RFC3339Nano),
	}

	fmt.Println(index.CreateTime)
	err := client.IndexQuestion(index)
	if err != nil {
		panic(err)
	}
	questions := client.SearchQuestions("test")
	fmt.Println(questions)
	res := client.GetQuestionById(1)
	fmt.Println(*res)
}

func TestSubmitSearch(t *testing.T) {
	appConfig := config.LoadConfig()
	client := NewElasticSearchClient(appConfig.Elasticsearch)
	var index *response.Submit
	index = &response.Submit{
		ID:          1,
		JudgeResult: []string{"test"},
		QuestionId:  1,
		UserId:      1,
		Status:      "ACCEPTED",
		CreateTime:  time.Now().Format(time.RFC3339Nano),
		UpdateTime:  time.Now().Format(time.RFC3339Nano),
	}
	err := client.IndexSubmit(index)
	if err != nil {
		panic(err)
	}
	submits := client.SearchSubmits("test")
	fmt.Println(submits)
	res := client.GetSubmitById(1)
	fmt.Println(*res)
}

func TestUpdateDocument(t *testing.T) {
	appConfig := config.LoadConfig()
	client := NewElasticSearchClient(appConfig.Elasticsearch)
	var index *response.Question
	index = &response.Question{
		ID:      1,
		Title:   "test",
		Content: "update content",
	}

	err := client.UpdateQuestion(index)
	if err != nil {
		panic(err)
	}
}

func TestModel(t *testing.T) {
	data1 := response.Question{
		ID: 1,
	}
	fmt.Println(unsafe.Sizeof(data1))
	serializedata1, _ := json.Marshal(data1)
	fmt.Println(len(serializedata1))
	data2 := response.Question{
		ID:         2,
		Title:      "hello world",
		Content:    "hello world",
		Tag:        []string{"test", "test"},
		Answer:     []string{"test", "test"},
		SubmitNum:  10,
		AcceptNum:  1,
		UserId:     1,
		CreateTime: time.Now().Format(time.RFC3339Nano),
		UpdateTime: time.Now().Format(time.RFC3339Nano),
		JudgeCase:  []string{"test", "test"},
	}
	fmt.Println(unsafe.Sizeof(data2))
	serializedata2, _ := json.Marshal(data2)
	fmt.Println(len(serializedata2))
}
