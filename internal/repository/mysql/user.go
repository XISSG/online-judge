package mysql

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (mysql *MysqlClient) CreateUser(user *entity.User) error {
	err := createData[entity.User](mysql, constant.USER_TABLE, user)
	if err != nil {
		return err
	}
	return nil
}

func (mysql *MysqlClient) GetUserById(userId int) *entity.User {
	user := getDataById[entity.User](mysql, constant.USER_TABLE, userId)
	return user
}

func (mysql *MysqlClient) GetUserByName(name string) *entity.User {
	var user *entity.User
	tx := mysql.client.Model("user").Where("user_name = ?", name).First(&user)
	if tx.Error != nil {
		return nil
	}
	return user
}

func (mysql *MysqlClient) GetUserListByName(userName string) []*entity.User {
	var users []*entity.User
	tx := mysql.client.Model("user").Where("user_name like ?", "%"+userName+"%").Find(&users)
	if tx.Error != nil {
		return nil
	}

	return users
}

func (mysql *MysqlClient) UpdateUser(user *entity.User) error {
	return updateDataById[entity.User](mysql, constant.USER_TABLE, user.ID, user)
}

func (mysql *MysqlClient) BanUser(userId int) error {
	tx := mysql.client.Begin()
	tx.Model("user").Where("id = ?", userId).Update("user_role", constant.BAN)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}

func (mysql *MysqlClient) DeleteUser(userId int) error {
	err := deleteDataById(mysql, constant.USER_TABLE, userId)
	if err != nil {
		return err
	}
	return nil
}
