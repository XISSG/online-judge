package mysql

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/entity"
)

func (mysql *MysqlClient) CreateUser(user *entity.User) error {
	tx := mysql.Begin()
	tx.Model("user").Create(user)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}

func (mysql *MysqlClient) GetUserById(userId string) *entity.User {
	var user *entity.User
	tx := mysql.Model("user").Where("id = ?", userId).First(&user)
	if tx.Error != nil {
		return nil
	}
	return user
}

func (mysql *MysqlClient) GetUserByName(name string) *entity.User {
	var user *entity.User
	tx := mysql.Model("user").Where("user_name = ?", name).First(&user)
	if tx.Error != nil {
		return nil
	}
	return user
}

func (mysql *MysqlClient) GetUserListByName(userName string) []*entity.User {
	var users []*entity.User
	tx := mysql.Model("user").Where("user_name like ?", "%"+userName+"%").Find(&users)
	if tx.Error != nil {
		return nil
	}

	return users
}

func (mysql *MysqlClient) UpdateUserPassword(userId string, password string) error {
	tx := mysql.Begin()
	tx.Model("user").Where("id = ?", userId).Update("password", password)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}

func (mysql *MysqlClient) UpdateUserAvatar(userId string, avatar string) error {

	tx := mysql.Begin()
	tx.Model("user").Where("id = ?", userId).Update("avatar_url", avatar)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}

func (mysql *MysqlClient) BanUser(userId string) error {
	tx := mysql.Begin()
	tx.Model("user").Where("id = ?", userId).Update("user_role", constant.BAN)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}

func (mysql *MysqlClient) DeleteUser(userId string) error {
	tx := mysql.Begin()
	tx.Model("user").Where("id = ?", userId).Update("is_delete", 1)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}
