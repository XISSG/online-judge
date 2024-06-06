package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/xissg/online-judge/internal/model/response"
	"log"
	"strconv"
)

func (es *ElasticsearchClient) IndexQuestions(question *response.Question) error {
	res, err := es.Client.Index("question").
		Id(strconv.Itoa(question.ID)).
		Document(*question).
		Do(context.Background())

	if err != nil {
		return err
	}

	log.Println(res.Result)
	return nil
}

func (es *ElasticsearchClient) GetQuestionById(queryId string) *response.Question {
	resp, err := es.Client.Get("question", queryId).Do(context.Background())
	if err != nil {
		return nil
	}
	var res response.Question
	err = json.Unmarshal(resp.Source_, &res)
	if err != nil {
		return nil
	}
	return &res
}

func (es *ElasticsearchClient) SearchQuestions(query string) []*response.Question {
	request := &types.Query{
		MatchAll: &types.MatchAllQuery{
			QueryName_: &query,
		},
	}
	res, err := es.Client.Search().Index("question").Query(request).Do(context.Background())
	if err != nil {
		return nil
	}

	var questions []*response.Question
	for _, hit := range res.Hits.Hits {
		var question response.Question
		if err = json.Unmarshal(hit.Source_, &question); err != nil {
			log.Printf("Error unmarshalling hit: %s", err)
			continue
		}
		questions = append(questions, &question)
	}
	return questions
}
func (es *ElasticsearchClient) DeleteQuestionById(id string) error {
	_, err := es.Client.Delete("question", id).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
