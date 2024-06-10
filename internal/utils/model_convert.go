package utils

import (
	"encoding/json"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/common"
	"github.com/xissg/online-judge/internal/model/entity"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/model/response"
	"time"
)

func ConvertUserEntity(userRequest *request.User) *entity.User {
	var userEntity *entity.User

	userEntity.ID = Snowflake()
	userEntity.UserRole = constant.USER
	userEntity.UserName = userRequest.UserName
	userEntity.IsDelete = constant.NOT_DELETED
	userEntity.AvatarURL = userRequest.AvatarURL
	userEntity.CreateTime = time.Now().Format(time.RFC3339Nano)
	userEntity.UpdateTime = time.Now().Format(time.RFC3339Nano)
	userEntity.UserPassword = MD5Crypt(userRequest.UserPassword)

	return userEntity
}

func ConvertUserResponse(userEntity *entity.User) *response.User {
	var userResponse *response.User

	userResponse.ID = userEntity.ID
	userResponse.UserName = userEntity.UserName
	userResponse.AvatarURL = userEntity.AvatarURL
	userResponse.CreateTime = userEntity.CreateTime
	userResponse.UpdateTime = userEntity.UpdateTime

	return userResponse
}

func ConvertQuestionEntity(questionRequest *request.Question) *entity.Question {
	var questionEntity *entity.Question

	tag, err := json.Marshal(questionRequest.Tag)
	answer, err := json.Marshal(questionRequest.Answer)
	judgeCase, err := json.Marshal(questionRequest.JudgeCase)
	judgeConfig, err := json.Marshal(questionRequest.JudgeConfig)

	if err != nil {
		return nil
	}

	questionEntity.AcceptNum = 0
	questionEntity.SubmitNum = 0
	questionEntity.ID = Snowflake()
	questionEntity.Tag = string(tag)
	questionEntity.Answer = string(answer)
	questionEntity.JudgeCase = string(judgeCase)
	questionEntity.Title = questionRequest.Title
	questionEntity.IsDelete = constant.NOT_DELETED
	questionEntity.JudgeConfig = string(judgeConfig)
	questionEntity.Content = questionRequest.Content
	questionEntity.CreateTime = time.Now().Format(time.RFC3339Nano)
	questionEntity.UpdateTime = time.Now().Format(time.RFC3339Nano)

	return questionEntity
}

func ConvertQuestionResponse(questionEntity *entity.Question) *response.Question {
	var tag []string
	var answer []string
	var judgeCase []string
	var judgeConfig common.Config
	var questionResponse *response.Question

	err := json.Unmarshal([]byte(questionEntity.Tag), &tag)
	err = json.Unmarshal([]byte(questionEntity.Answer), &answer)
	err = json.Unmarshal([]byte(questionEntity.JudgeCase), &judgeCase)
	err = json.Unmarshal([]byte(questionEntity.JudgeConfig), &judgeConfig)

	if err != nil {
		return nil
	}

	questionResponse.Tag = tag
	questionResponse.Answer = answer
	questionResponse.JudgeCase = judgeCase
	questionResponse.ID = questionEntity.ID
	questionResponse.JudgeConfig = judgeConfig
	questionResponse.Title = questionEntity.Title
	questionResponse.UserId = questionEntity.UserId
	questionResponse.Content = questionEntity.Content
	questionResponse.AcceptNum = questionEntity.AcceptNum
	questionResponse.SubmitNum = questionEntity.SubmitNum
	questionResponse.CreateTime = questionEntity.CreateTime
	questionResponse.UpdateTime = questionEntity.UpdateTime

	return questionResponse
}

func UpdateQuestionToQuestionEntity(updateRequest *request.UpdateQuestion) *entity.Question {
	var questionEntity *entity.Question

	if updateRequest.Tag != nil {
		tag, _ := json.Marshal(updateRequest.Tag)
		questionEntity.Tag = string(tag)
	}

	if updateRequest.Answer != nil {
		answer, _ := json.Marshal(updateRequest.Answer)
		questionEntity.Answer = string(answer)
	}

	if updateRequest.JudgeCase != nil {
		judgeCase, _ := json.Marshal(updateRequest.JudgeCase)
		questionEntity.JudgeCase = string(judgeCase)
	}

	if updateRequest.JudgeConfig != nil {
		judgeConfig, _ := json.Marshal(updateRequest.JudgeConfig)
		questionEntity.JudgeConfig = string(judgeConfig)
	}

	questionEntity.ID = updateRequest.ID
	questionEntity.Title = updateRequest.Title
	questionEntity.Content = updateRequest.Content
	questionEntity.UpdateTime = time.Now().Format(time.RFC3339Nano)

	return questionEntity
}

func ConvertSubmitEntity(submitRequest *request.Submit, userId int) *entity.Submit {
	var submitEntity *entity.Submit

	submitEntity.Status = constant.WATING_STATUS
	submitEntity.UserId = userId
	submitEntity.ID = Snowflake()
	submitEntity.JudgeResult = ""
	submitEntity.Code = submitRequest.Code
	submitEntity.IsDelete = constant.NOT_DELETED
	submitEntity.Language = submitRequest.Language
	submitEntity.QuestionId = submitRequest.QuestionId
	submitEntity.CreateTime = time.Now().Format(time.RFC3339Nano)
	submitEntity.UpdateTime = time.Now().Format(time.RFC3339Nano)

	return submitEntity
}
func ConvertSubmitResponse(submitEntity *entity.Submit) *response.Submit {
	var submitResponse *response.Submit

	submitResponse.ID = submitEntity.ID
	submitResponse.UserId = submitEntity.UserId
	submitResponse.QuestionId = submitEntity.QuestionId
	submitResponse.Status = submitEntity.Status
	submitResponse.JudgeResult = submitEntity.JudgeResult
	submitResponse.CreateTime = submitEntity.CreateTime
	submitResponse.UpdateTime = submitEntity.UpdateTime

	return submitResponse
}

func UpdateSubmitToSubmitEntity(updateRequest *request.UpdateSubmit) *entity.Submit {
	var submitEntity *entity.Submit

	submitEntity.ID = updateRequest.ID
	submitEntity.Status = updateRequest.Status
	submitEntity.JudgeResult = updateRequest.JudgeResult
	submitEntity.UpdateTime = time.Now().Format(time.RFC3339Nano)
	return submitEntity
}
