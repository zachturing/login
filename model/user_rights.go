package model

import (
	"fmt"
	"github.com/newdee/aipaper-util/database/mysql"
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

func QueryUserRightsByUserId(userID int, tx *gorm.DB) (*UserRights, error) {
	var userRights UserRights
	if tx.Where("user_id = ?", userID).First(&userRights).RowsAffected != 1 {
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
