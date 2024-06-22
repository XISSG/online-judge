package service

import (
	"github.com/xissg/online-judge/internal/model/entity"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/model/response"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/utils"
	"time"
)

type UserService interface {
	CreateUser(userRequest *request.User) error
	GetUserById(userId int) *response.User
	GetUserByUsername(username string) *response.User
	GetUserListByUsername(username string) []*response.User
	UpdateUserPassword(userId int, password string) error
	UpdateUserAvatar(userId int, avatar string) error
	DeleteUserById(userId int) error
	BanUserById(userId int) error
	CheckUserNameAndPasswd(userName, password string) (bool, int, string)
	CheckUserExists(userName string) bool
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

func (s *userService) GetUserById(userId int) *response.User {
	user := s.mysql.GetUserById(userId)
	if user == nil {
		return nil
	}
	userResponse := utils.ConvertUserResponse(user)
	return userResponse
}
func (s *userService) GetUserByUsername(username string) *response.User {
	user := s.mysql.GetUserByName(username)
	if user == nil {
		return nil
	}
	return utils.ConvertUserResponse(user)
}

func (s *userService) GetUserListByUsername(username string) []*response.User {
	var userResponses []*response.User
	users := s.mysql.GetUserListByName(username)
	if len(users) == 0 {
		return nil
	}

	for i := range users {
		userResponse := utils.ConvertUserResponse(users[i])
		userResponses = append(userResponses, userResponse)
	}

	return userResponses
}

func (s *userService) UpdateUserPassword(userId int, password string) error {
	password = utils.MD5Crypt(password)
	user := &entity.User{
		ID:           userId,
		UserPassword: password,
		UpdateTime:   time.Now().Format(time.RFC3339Nano),
	}
	err := s.mysql.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUserAvatar(userId int, avatar string) error {
	user := &entity.User{
		ID:         userId,
		AvatarURL:  avatar,
		UpdateTime: time.Now().Format(time.RFC3339Nano),
	}
	err := s.mysql.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) DeleteUserById(userId int) error {
	err := s.mysql.DeleteUser(userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) BanUserById(userId int) error {
	err := s.mysql.BanUser(userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) CheckUserNameAndPasswd(userName, password string) (bool, int, string) {
	user := s.mysql.GetUserByName(userName)
	if user == nil {
		return false, 0, ""
	}
	if utils.MD5Crypt(password) != user.UserPassword {
		return false, 0, ""
	}
	return true, user.ID, user.UserRole
}

func (s *userService) CheckUserExists(userName string) bool {
	user := s.mysql.GetUserByName(userName)
	if user == nil {
		return false
	}
	return true
}
