package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/xissg/online-judge/internal/model/response"
	"strconv"
)

type RESPONSE interface {
	response.Question | response.Submit
}

func indexDoc[T RESPONSE](es *ESClient, index string, indexId int, data *T) error {
	id := strconv.Itoa(indexId)
	_, err := es.client.Index(index).Id(id).Document(*data).Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func searchDocsByQuery[T RESPONSE](es *ESClient, index string, query string) []*T {
	request := &types.Query{
		MatchAll: &types.MatchAllQuery{
			QueryName_: &query,
		},
	}
	res, err := es.client.Search().Index(index).Query(request).Do(context.Background())
	if err != nil {
		return nil
	}

	var results []*T
	for _, hit := range res.Hits.Hits {
		var t T
		if err = json.Unmarshal(hit.Source_, &t); err != nil {
			continue
		}
		results = append(results, &t)
	}
	return results
}

func updateDocsById[T RESPONSE](es *ESClient, index string, queryId int, data *T) error {
	id := strconv.Itoa(queryId)
	_, err := es.client.Update(index, id).Doc(*data).Do(context.Background())
	return err
}
func getDocById[T RESPONSE](es *ESClient, index string, queryId int) *T {
	id := strconv.Itoa(queryId)
	resp, err := es.client.Get(index, id).Do(context.Background())
	if err != nil {
		return nil
	}
	var res T
	err = json.Unmarshal(resp.Source_, &res)
	if err != nil {
		return nil
	}
	return &res
}

func deleteById(es *ESClient, index string, queryId int) error {
	id := strconv.Itoa(queryId)
	_, err := es.client.Delete(index, id).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
