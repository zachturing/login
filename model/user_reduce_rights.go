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
