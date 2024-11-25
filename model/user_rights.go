package model

import (
	"fmt"
	"github.com/zachturing/util/database/mysql"
	"gorm.io/gorm"
)

func CreateUserRights(userRights *UserRights, tx *gorm.DB) error {
	if userRights == nil {
		return fmt.Errorf("avatar nil")
	}
	if tx != nil {
		return tx.Create(userRights).Error
	}
	return mysql.GetGlobalDBIns().Create(userRights).Error
}

func QueryUserRightsByUserId(userID int) (*UserRights, error) {
	var userRights UserRights
	tx := mysql.GetGlobalDBIns().Where("user_id = ?", userID).First(&userRights)
	if tx.RowsAffected != 1 {
		return nil, fmt.Errorf("userId:%v not exist", userID)
	}
	return &userRights, nil
}

func SaveUserRights(userRights *UserRights, tx *gorm.DB) error {
	if userRights == nil {
		return fmt.Errorf("avatar nil")
	}
	if tx != nil {
		return tx.Save(userRights).Error
	}
	return mysql.GetGlobalDBIns().Save(userRights).Error
}
