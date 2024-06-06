package utils

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/model/response"
	"time"
)

func ConvertUserEntity(userRequest *request.User) *entity.User {
	var userEntity *entity.User
	userEntity.ID = Snowflake()
	userEntity.UserName = userRequest.UserName
	userEntity.UserPassword = MD5Crypt(userRequest.UserPassword)
	userEntity.AvatarURL = userRequest.AvatarURL
	userEntity.CreateTime = time.Now().UTC()
	userEntity.UpdateTime = time.Now().UTC()
	userEntity.IsDelete = constant.NOT_DELETED
	userEntity.UserRole = constant.USER
	return userEntity
}

func ConvertUseResponse(userEntity *entity.User) *response.User {
	var userResponse *response.User
	userResponse.ID = userEntity.ID
	userResponse.UserName = userEntity.UserName
	userResponse.AvatarURL = userEntity.AvatarURL
	userResponse.CreateTime = userEntity.CreateTime
	userResponse.UpdateTime = userEntity.UpdateTime
	return userResponse
}
