package model

import (
	"github.com/newdee/aipaper-util/database/mysql"
	"gorm.io/gorm"
	"time"
)

type UserReduceRights struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID        int64     `gorm:"column:user_id;not null" json:"user_id"`
	RemainingNum  int       `gorm:"column:remaining_num;not null" json:"remaining_num"`
	UsedReduceNum int       `gorm:"column:used_reduce_num" json:"used_reduce_num"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
}

func (u *UserReduceRights) TableName() string {
	return "user_reduce_rights"
}

func GiftUserRights(userId int64, tx *gorm.DB) error {
	if tx == nil {
		tx = mysql.GetGlobalDBIns()
	}
	return tx.Create(&UserReduceRights{
		UserID:        userId,
		RemainingNum:  3,
		UsedReduceNum: 0,
	}).Error
}
