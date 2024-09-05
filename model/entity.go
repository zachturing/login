package model

import "time"

// User 代表数据库中的一条用户记录
type User struct {
	ID               int64     `json:"id"`                // 自增ID，作为主键
	Phone            string    `json:"phone"`             // 手机号码
	RegistrationTime time.Time `json:"registration_time"` // 注册时间
	LastLoginTime    time.Time `json:"last_login_time"`   // 上一次登录时间
}

func (u *User) TableName() string {
	return "user"
}
