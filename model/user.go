package model

import (
	"fmt"
	"gorm.io/gorm"

	"github.com/newdee/aipaper-util/database/mysql"
)

// CreateUser 创建用户
func CreateUser(user *User, tx *gorm.DB) error {
	if user == nil {
		return fmt.Errorf("avatar nil")
	}
	return tx.Create(user).Error
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

func UpdateUserColumns(userId int64, updateColumns map[string]interface{}, tx *gorm.DB) error {
	if len(updateColumns) == 0 {
		return nil
	}
	if tx == nil {
		tx = mysql.GetGlobalDBIns()
	}
	return tx.Model(&User{}).Where("id = ?", userId).Updates(updateColumns).Error
}
