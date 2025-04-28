package model

import (
	"github.com/newdee/aipaper-util/database/mysql"
	"gorm.io/gorm"
	"time"
)

type UserReduceLogs struct {
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID           int64     `gorm:"column:user_id;not null" json:"user_id"`
	PreReduceNum     int       `gorm:"column:pre_reduce_num;not null" json:"pre_reduce_num"`
	ChangeNum        int       `gorm:"column:change_num;not null" json:"change_num"`
	PostReduceNum    int       `gorm:"column:post_reduce_num;not null" json:"post_reduce_num"`
	ChangeReason     string    `gorm:"column:change_reason;type:varchar(10);not null" json:"change_reason"`
	OriginalContents string    `gorm:"column:original_contents" json:"original_contents"`
	PostContents     string    `gorm:"column:post_contents" json:"post_contents"`
	CreatedAt        time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
}

func (u *UserReduceLogs) TableName() string {
	return "user_reduce_logs"
}

func InsertUserReduceLogs(userId int64, preReduceNum, changeNum, postReduceNum int, changeReason, originalContents, postContents string, tx *gorm.DB) error {
	if tx == nil {
		tx = mysql.GetGlobalDBIns()
	}
	userReduceLogs := &UserReduceLogs{
		UserID:           userId,
		PreReduceNum:     preReduceNum,
		ChangeNum:        changeNum,
		PostReduceNum:    postReduceNum,
		ChangeReason:     changeReason,
		OriginalContents: originalContents,
		PostContents:     postContents,
	}
	return tx.Create(userReduceLogs).Error
}
