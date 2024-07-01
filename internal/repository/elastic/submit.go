package elastic

import (
	"fmt"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/response"
)

func (es *ESClient) IndexSubmit(submit *response.Submit) error {
	err := indexDoc[response.Submit](es, constant.SUBMIT_INDEX, submit.ID, submit)
	if err != nil {
		return fmt.Errorf("repository layer: elasticsearch, index submit: %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (es *ESClient) GetSubmitById(queryId int) (*response.Submit, error) {
	data := getDocById[response.Submit](es, constant.SUBMIT_INDEX, queryId)
	if data == nil {
		return nil, fmt.Errorf("repository layer: elasticsearch, get submit by id: %w", constant.ErrNotFound)
	}
	return data, nil
}

func (es *ESClient) SearchSubmits(query string) ([]*response.Submit, error) {
	data := searchDocsByQuery[response.Submit](es, constant.SUBMIT_INDEX, query, []string{"title,content"})
	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("repository layer: elasticsearch, search submits: %w", constant.ErrNotFound)
	}
	return data, nil
}

func (es *ESClient) UpdateSubmit(submit *response.Submit) error {
	err := updateDocsById[response.Submit](es, constant.SUBMIT_INDEX, submit.ID, submit)
	if err != nil {
		return fmt.Errorf("repository layer: elasticsearch, update submit: %w", constant.ErrInternal)
	}
	return nil
}

func (es *ESClient) DeleteSubmitById(queryId int) error {
	err := deleteById(es, constant.SUBMIT_INDEX, queryId)
	if err != nil {
		return fmt.Errorf("repository layer: elasticsearch, delete submit by id: %w", constant.ErrInternal)
	}
	return nil
}
