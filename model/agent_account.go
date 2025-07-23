package model

import (
	"github.com/newdee/aipaper-util/database/mysql"
	"gorm.io/gorm"
)

func CreateAgentAccount(agentAccount *AgentAccount, tx *gorm.DB) error {
	if tx == nil {
		tx = mysql.GetGlobalDBIns()
	}
	return tx.Create(agentAccount).Error
}
