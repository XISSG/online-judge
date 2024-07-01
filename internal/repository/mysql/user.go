package mysql

import (
	"fmt"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (mysql *MysqlClient) CreateUser(user *entity.User) error {
	err := createData[entity.User](mysql, constant.USER_TABLE, user)
	if err != nil {
		return fmt.Errorf("repository layer: mysql, create user, %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (mysql *MysqlClient) GetUserById(userId int) (*entity.User, error) {
	user, err := getDataById[entity.User](mysql, constant.USER_TABLE, userId)
	if err != nil {
		return nil, fmt.Errorf("repository layer: mysql, get user by id, %w %+v", constant.ErrInternal, err)
	}

	if user == nil {
		return nil, fmt.Errorf("repository layer: mysql, get user by id, %w", constant.ErrNotFound)
	}
	return user, nil
}

func (mysql *MysqlClient) GetUserByName(name string) (*entity.User, error) {
	var user *entity.User
	tx := mysql.client.Table(constant.USER_TABLE).Where("user_name = ?", name).First(&user)
	if tx.Error != nil {
		return nil, fmt.Errorf("repository layer: mysql, get user by name, %w %+v", constant.ErrInternal, tx.Error)
	}
	if user == nil {
		return nil, fmt.Errorf("repository layer mysql: get user by name %w", constant.ErrNotFound)
	}
	return user, nil
}

func (mysql *MysqlClient) GetUserListByName(userName string) ([]*entity.User, error) {
	var users []*entity.User
	tx := mysql.client.Table(constant.USER_TABLE).Where("user_name like ?", "%"+userName+"%").Find(&users)
	if tx.Error != nil {
		return nil, fmt.Errorf("repository layer: mysql, get user, %w %+v", constant.ErrInternal, tx.Error)
	}

	if users == nil || len(users) == 0 {
		return nil, fmt.Errorf("repository layer: mysql, get user, %w", constant.ErrNotFound)
	}

	return users, nil
}

func (mysql *MysqlClient) UpdateUser(user *entity.User) error {
	err := updateDataById[entity.User](mysql, constant.USER_TABLE, user.ID, user)
	if err != nil {
		return fmt.Errorf("repository layer: mysql, update user, %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (mysql *MysqlClient) BanUser(userId int) error {
	user := &entity.User{
		ID:       userId,
		UserRole: constant.BAN,
	}
	err := updateDataById[entity.User](mysql, constant.USER_TABLE, userId, user)
	if err != nil {
		return fmt.Errorf("repository layer: mysql, ban user, %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (mysql *MysqlClient) DeleteUser(userId int) error {
	err := deleteDataById(mysql, constant.USER_TABLE, userId)
	if err != nil {
		return fmt.Errorf("repository layer: mysql, delete user, %w %+v", constant.ErrInternal, err)
	}
	return nil
}
