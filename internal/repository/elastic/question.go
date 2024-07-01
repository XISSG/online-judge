package elastic

import (
	"fmt"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/response"
)

func (es *ESClient) IndexQuestion(question *response.Question) error {
	err := indexDoc[response.Question](es, constant.QUESTION_INDEX, question.ID, question)
	if err != nil {
		return fmt.Errorf("repository layer: elasticsearch, index question: %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (es *ESClient) GetQuestionById(queryId int) (*response.Question, error) {
	data := getDocById[response.Question](es, constant.QUESTION_INDEX, queryId)
	if data == nil {
		return nil, fmt.Errorf("repository layer: elasticsearch, get question by id: %w", constant.ErrNotFound)
	}
	return data, nil
}

func (es *ESClient) SearchQuestions(query string) ([]*response.Question, error) {
	data := searchDocsByQuery[response.Question](es, constant.QUESTION_INDEX, query, []string{"title", "content", "tags"})
	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("repository layer: elasticsearch, search questions: %w", constant.ErrNotFound)
	}
	return data, nil
}

func (es *ESClient) UpdateQuestion(question *response.Question) error {
	err := updateDocsById[response.Question](es, constant.QUESTION_INDEX, question.ID, question)
	if err != nil {
		return fmt.Errorf("repository layer: elasticsearch, update question: %w", constant.ErrInternal)
	}
	return nil
}

func (es *ESClient) DeleteQuestionById(queryId int) error {
	err := deleteById(es, constant.QUESTION_INDEX, queryId)
	if err != nil {
		return fmt.Errorf("repository layer: elasticsearch, delete question by id: %w", constant.ErrInternal)
	}
	return nil
}

func (es *ESClient) UpdateQuestionAcceptNum(questionId int) error {
	question, err := es.GetQuestionById(questionId)
	if err != nil {
		return err
	}
	acceptNum := question.AcceptNum + 1
	data := &response.Question{
		AcceptNum: acceptNum,
	}
	err = updateDocsById[response.Question](es, constant.QUESTION_INDEX, questionId, data)
	if err != nil {
		return fmt.Errorf("repository layer: elasticsearch, update question accept num: %w", constant.ErrInternal)
	}
	return nil
}

func (es *ESClient) UpdateQuestionSubmitNum(questionId int) error {
	question, err := es.GetQuestionById(questionId)
	if err != nil {
		return err
	}
	submitNum := question.SubmitNum + 1
	data := &response.Question{
		AcceptNum: submitNum,
	}
	question.SubmitNum++
	err = updateDocsById[response.Question](es, constant.QUESTION_INDEX, questionId, data)
	if err != nil {
		return fmt.Errorf("repository layer: elasticsearch, update question submit num: %w", constant.ErrInternal)
	}
	return nil
}
