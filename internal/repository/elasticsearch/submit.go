package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/xissg/online-judge/internal/model/response"
	"log"
	"strconv"
)

func (es *ElasticsearchClient) IndexSubmit(submit *response.Submit) error {
	res, err := es.Client.Index("submit").
		Id(strconv.Itoa(submit.ID)).
		Document(*submit).
		Do(context.Background())

	if err != nil {
		return err
	}

	log.Println(res.Result)
	return nil
}

func (es *ElasticsearchClient) GetSubmitById(queryId string) *response.Submit {
	resp, err := es.Client.Get("submit", queryId).Do(context.Background())
	if err != nil {
		return nil
	}
	var res response.Submit
	err = json.Unmarshal(resp.Source_, &res)
	if err != nil {
		return nil
	}
	return &res
}

func (es *ElasticsearchClient) SearchSubmits(query string) []*response.Submit {
	request := &types.Query{
		MatchAll: &types.MatchAllQuery{
			QueryName_: &query,
		},
	}
	res, err := es.Client.Search().Index("submit").Query(request).Do(context.Background())
	if err != nil {
		return nil
	}

	var submits []*response.Submit
	for _, hit := range res.Hits.Hits {
		var submit response.Submit
		if err = json.Unmarshal(hit.Source_, &submit); err != nil {
			log.Printf("Error unmarshalling hit: %s", err)
			continue
		}
		submits = append(submits, &submit)
	}
	return submits
}

func (es *ElasticsearchClient) DeleteSubmitById(id string) error {
	_, err := es.Client.Delete("submit", id).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
