package model

import "time"

// User 代表数据库中的一条用户记录
type User struct {
	ID               int64     `json:"id"`                // 自增ID，作为主键
	Phone            string    `json:"phone"`             // 手机号码
	UserName         string    `json:"user_name"`         // 用户名
	RegistrationTime time.Time `json:"registration_time"` // 注册时间
	LastLoginTime    time.Time `json:"last_login_time"`   // 上一次登录时间
	Role             string    `json:"role"`              // 用户级别
	Permission       string    `json:"permission"`        // 用户权限
	AgentId          int       `json:"agent_id"`          // 用户所属的代理商Id
	ParentUserId     int64     `json:"parent_user_id"`    // 用户所属的邀请人Id
	InvCode          string    `json:"inv_code"`          // 用户邀请码，根据userId生成，每个用户唯一
}

// InvitationLogs 代表数据库中的一条邀请记录
type InvitationLogs struct {
	ID                 int64     `json:"id"`
	InviteeId          int64     `json:"invitee_id"`
	InviteeName        string    `json:"invitee_name"`
	InviterId          int64     `json:"inviter_id"`
	InviteeRewardsType string    `json:"invitee_rewards_type"`
	InviterRewardsType string    `json:"inviter_rewards_type"`
	Remarks            string    `json:"remarks"`
	CreatedAt          time.Time `json:"created_at"`
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

// DistributionAccount 分销商账户表
type DistributionAccount struct {
	ID                 int64     `json:"id"`                   // 自增主键
	UserID             int64     `json:"user_id"`              // user表id
	Currency           string    `json:"currency"`             // 币种，CNY-人民币、USD-美元，默认：CNY
	Balance            float64   `json:"balance"`              // 账户余额，两位小数
	FrozenAmount       float64   `json:"frozen_amount"`        // 待结算金额，两位小数
	WithdrawnAmount    float64   `json:"withdrawn_amount"`     // 已提现金额，两位小数
	TotalIncome        float64   `json:"total_income"`         // 总收益，两位小数
	DirectPercent      float64   `json:"direct_percent"`       // 直推分成比例：默认20%
	IndirectPercent    float64   `json:"indirect_percent"`     // 间推分成比例：默认10%
	UserUpgradePercent float64   `json:"user_upgrade_percent"` // 直推的用户升级成为代理时，代理费分成：默认80%
	Status             string    `json:"status"`               // 账户状态，ACTIVE-活跃、FROZEN-冻结、CLOSED-关闭，默认：ACTIVE
	CreatedAt          time.Time `json:"created_at"`           // 账户创建时间
	UpdatedAt          time.Time `json:"updated_at"`           // 账户更新时间
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

type UserReduceRights struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID        int64     `gorm:"column:user_id;not null" json:"user_id"`
	RemainingNum  int       `gorm:"column:remaining_num;not null" json:"remaining_num"`
	UsedReduceNum int       `gorm:"column:used_reduce_num" json:"used_reduce_num"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
}

func (u *User) TableName() string {
	return "user"
}

func (a *Agent) TableName() string {
	return "agents"
}

func (a *DistributionAccount) TableName() string {
	return "distribution_account"
}

func (u *UserRights) TableName() string {
	return "user_rights"
}

func (u *UserReduceRights) TableName() string {
	return "user_reduce_rights"
}

func (i *InvitationLogs) TableName() string {
	return "invitation_logs"
}
