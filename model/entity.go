package model

import "time"

// User 代表数据库中的一条用户记录
type User struct {
	ID               int64     `json:"id"`                // 自增ID，作为主键
	Phone            string    `json:"phone"`             // 手机号码
	RegistrationTime time.Time `json:"registration_time"` // 注册时间
	LastLoginTime    time.Time `json:"last_login_time"`   // 上一次登录时间
	Role             string    `json:"role"`              // 用户角色，可以为空
	SubDomain        string    `json:"sub_domain"`        // 代理商二级域名前缀，如zhangsan.mixpaper.cn，前缀为zhangsan
}

func (u *User) TableName() string {
	return "user"
}
