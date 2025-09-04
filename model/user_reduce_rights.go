package model

import (
	"github.com/newdee/aipaper-util/database/mysql"
	"gorm.io/gorm"
)

func CreateUserReduceRights(userReduceRights *UserReduceRights, tx *gorm.DB) error {
	if tx == nil {
		tx = mysql.GetGlobalDBIns()
	}
	return tx.Create(userReduceRights).Error
}

func QueryUserReduceRights(userId int64, tx *gorm.DB) (*UserReduceRights, error) {
	if tx == nil {
		tx = mysql.GetGlobalDBIns()
	}
	var userReduceRights UserReduceRights
	tx.Where("user_id = ?", userId).First(&userReduceRights)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &userReduceRights, nil
}

func UpdateUserReduceRightsColumnsByUserId(userId int64, updateColumns map[string]interface{}, tx *gorm.DB) error {
	if len(updateColumns) == 0 {
		return nil
	}
	if tx == nil {
		tx = mysql.GetGlobalDBIns()
	}
	return tx.Model(&UserReduceRights{}).Where("user_id = ?", userId).Updates(updateColumns).Error
}
