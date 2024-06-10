package elastic

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/response"
)

func (es *ESClient) IndexQuestion(question *response.Question) error {
	return indexDoc[response.Question](es, constant.QUESTION_INDEX, question.ID, question)
}

func (es *ESClient) GetQuestionById(queryId int) *response.Question {
	return getDocById[response.Question](es, constant.QUESTION_INDEX, queryId)
}

func (es *ESClient) SearchQuestions(query string) []*response.Question {
	return searchDocsByQuery[response.Question](es, constant.QUESTION_INDEX, query)
}

func (es *ESClient) DeleteQuestionById(queryId int) error {
	return deleteById(es, constant.QUESTION_INDEX, queryId)
}

func (es *ESClient) UpdateQuestion(question *response.Question) error {
	return updateDocsById[response.Question](es, constant.QUESTION_INDEX, question.ID, question)
}

func (es *ESClient) UpdateQuestionAcceptNum(questionId int) error {
	question := es.GetQuestionById(questionId)
	acceptNum := question.AcceptNum + 1
	data := &response.Question{
		AcceptNum: acceptNum,
	}
	return updateDocsById[response.Question](es, constant.QUESTION_INDEX, questionId, data)
}

func (es *ESClient) UpdateQuestionSubmitNum(questionId int) error {
	question := es.GetQuestionById(questionId)
	submitNum := question.SubmitNum + 1
	data := &response.Question{
		AcceptNum: submitNum,
	}
	question.SubmitNum++
	return updateDocsById[response.Question](es, constant.QUESTION_INDEX, questionId, data)
}
