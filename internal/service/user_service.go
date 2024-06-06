package service

import (
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/model/response"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/utils"
)

type UserService interface {
	CreateUser(userRequest *request.User) error
	GetUserById(userId string) *response.User
	GetUserByUsername(username string) *response.User
	GetUserListByUsername(username string) []*response.User
	UpdateUserPassword(userId, password string) error
	UpdateUserAvatar(userId, avatar string) error
	DeleteUserById(userId string) error
	BanUserById(userId string) error
	CheckUser(userName, password string) bool
}

type userService struct {
	mysql *mysql.MysqlClient
}

func NewUserService(mysql *mysql.MysqlClient) UserService {
	return &userService{
		mysql: mysql,
	}
}

func (s *userService) CreateUser(userRequest *request.User) error {
	user := utils.ConvertUserEntity(userRequest)
	err := s.mysql.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) GetUserById(userId string) *response.User {
	user := s.mysql.GetUserById(userId)
	if user == nil {
		return nil
	}
	userResponse := utils.ConvertUseResponse(user)
	return userResponse
}
func (s *userService) GetUserByUsername(username string) *response.User {
	user := s.mysql.GetUserByName(username)
	if user == nil {
		return nil
	}
	return utils.ConvertUseResponse(user)
}

func (s *userService) GetUserListByUsername(username string) []*response.User {
	var userResponses []*response.User
	users := s.mysql.GetUserListByName(username)
	if len(users) == 0 {
		return nil
	}

	for i := range users {
		userResponse := utils.ConvertUseResponse(users[i])
		userResponses = append(userResponses, userResponse)
	}

	return userResponses
}

func (s *userService) UpdateUserPassword(userId, password string) error {
	password = utils.MD5Crypt(password)
	err := s.mysql.UpdateUserPassword(userId, password)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUserAvatar(userId, avatar string) error {
	err := s.mysql.UpdateUserAvatar(userId, avatar)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) DeleteUserById(userId string) error {
	err := s.mysql.DeleteUser(userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) BanUserById(userId string) error {
	err := s.mysql.BanUser(userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) CheckUser(userName, password string) bool {
	user := s.mysql.GetUserByName(userName)
	if user == nil {
		return false
	}
	if utils.MD5Crypt(password) != user.UserPassword {
		return false
	}
	return true
}
