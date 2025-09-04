package model

import (
	"github.com/newdee/aipaper-util/database/mysql"
	"gorm.io/gorm"
)

func CreateInvitationLogs(invitationLogs *InvitationLogs, tx *gorm.DB) error {
	if tx == nil {
		tx = mysql.GetGlobalDBIns()
	}
	return tx.Create(invitationLogs).Error
}

func CountInvitedUsers(inviterId int64, tx *gorm.DB) (int64, error) {
	if tx == nil {
		tx = mysql.GetGlobalDBIns()
	}
	var count int64
	err := tx.Model(&InvitationLogs{}).Where("inviter_id = ?", inviterId).Count(&count).Error
	return count, err
}
