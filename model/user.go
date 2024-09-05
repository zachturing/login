package model

import (
	"fmt"

	"github.com/zachturing/util/database/mysql"
)

// CreateUser 创建用户
func CreateUser(user *User) error {
	if user == nil {
		return fmt.Errorf("avatar nil")
	}
	return mysql.GetGlobalDBIns().Create(user).Error
}

// QueryUser 根据手机号查询用户
func QueryUser(phone string) (*User, error) {
	user := new(User)
	tx := mysql.GetGlobalDBIns().Where("phone = ?", phone).First(&user)
	if tx.RowsAffected != 1 {
		return nil, fmt.Errorf("phone:%v not exist", phone)
	}
	return user, nil
}
