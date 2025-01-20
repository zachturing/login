package model

import (
	"fmt"
	"gorm.io/gorm"
)

func QueryAgentBySubDomain(subDomain string, tx *gorm.DB) (*Agent, error) {
	var agent Agent
	if tx.Where("sub_domain = ?", subDomain).First(&agent).RowsAffected != 1 {
		return nil, fmt.Errorf("sub_domain:%v not exist", subDomain)
	}
	return &agent, nil
}
