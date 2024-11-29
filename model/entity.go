package model

import "time"

// User 代表数据库中的一条用户记录
type User struct {
	ID               int64     `json:"id"`                // 自增ID，作为主键
	Phone            string    `json:"phone"`             // 手机号码
	RegistrationTime time.Time `json:"registration_time"` // 注册时间
	LastLoginTime    time.Time `json:"last_login_time"`   // 上一次登录时间
	Role             string    `json:"role"`              // 用户级别
	Permission       string    `json:"permission"`        // 用户权限
	AgentId          int       `json:"agent_id"`          // 用户所属的代理商ID
	InvCode          string    `json:"inv_code"`          // 用户邀请码，根据userId生成，每个用户唯一
}

type Agent struct {
	ID         int64     `json:"id"`          // 代理商ID，自增
	Phone      string    `json:"phone"`       // 代理商手机号
	Name       string    `json:"name"`        // 代理商名称
	CreatedAt  time.Time `json:"created_at"`  // 代理商注册时间
	UpdatedAt  time.Time `json:"updated_at"`  // 更新时间
	SubDomain  string    `json:"sub_domain"`  // 代理商二级域名
	Verified   bool      `json:"verified"`    // 是否通过审核
	DomainFlag bool      `json:"domain_flag"` // 是否已绑定二级域名
	AgentLevel int       `json:"agent_level"` // 代理商等级，1-5级
	ParentId   int       `json:"parent_id"`   // 上级代理商ID
}

type UserRights struct {
	ID                 int64     `json:"id"`                   // 权益表的Id，自增
	UserId             int64     `json:"user_id"`              // 用户ID
	InvUsers           int       `json:"inv_users"`            // 用户已邀请的人数
	DuplicateCheckNums int       `json:"duplicate_check_nums"` // PaperYY免费查重次数，每次限制1W字
	UsedCheckNums      int       `json:"used_check_nums"`      // 已使用的查重次数
	CreatedAt          time.Time `json:"created_at"`           // 创建时间
	UpdatedAt          time.Time `json:"updated_at"`           // 更新时间
}

func (u *User) TableName() string {
	return "user"
}
