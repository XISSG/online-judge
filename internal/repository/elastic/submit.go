package elastic

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/response"
)

func (es *ESClient) IndexSubmit(submit *response.Submit) error {
	return indexDoc[response.Submit](es, constant.SUBMIT_INDEX, submit.ID, submit)
}

func (es *ESClient) GetSubmitById(queryId int) *response.Submit {
	return getDocById[response.Submit](es, constant.SUBMIT_INDEX, queryId)
}

func (es *ESClient) SearchSubmits(query string) []*response.Submit {
	return searchDocsByQuery[response.Submit](es, constant.SUBMIT_INDEX, query)
}

func (es *ESClient) UpdateSubmit(submit *response.Submit) error {
	return updateDocsById[response.Submit](es, constant.SUBMIT_INDEX, submit.ID, submit)
}

func (es *ESClient) DeleteSubmitById(queryId int) error {
	return deleteById(es, constant.SUBMIT_INDEX, queryId)
}
