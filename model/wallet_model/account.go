package wallet_model

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID        int             `json:"id" gorm:"primaryKey;autoIncrement;column:id"`  // 自增主键
	UserID    int             `json:"user_id" gorm:"not null;column:user_id"`        // 用户ID
	Currency  string          `json:"currency" gorm:"not null;column:currency"`      // 币种，CNY-人民币、USD-美元，默认：CNY
	Balance   decimal.Decimal `json:"balance" gorm:"not null;column:balance"`        // 资金余额
	Status    string          `json:"status" gorm:"not null;column:status"`          // 账户状态，ACTIVE-活跃、FROZEN-冻结、CLOSED-关闭，默认：ACTIVE
	CreatedAt time.Time       `json:"created_at" gorm:"not null"`                    // 创建时间
	UpdatedAt time.Time       `json:"updated_at" gorm:"not null"`                    // 更新时间
	Remarks   string          `json:"remarks" gorm:"type:varchar(255);default null"` // 备注字段
}

func (Account) TableName() string {
	return "account"
}

func CreateAccount(account *Account, tx *gorm.DB) error {
	if account == nil {
		return fmt.Errorf("avatar nil")
	}
	return tx.Create(account).Error
}
