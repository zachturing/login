package model

import (
	"fmt"
	"github.com/zachturing/util/database/mysql"
)

func QueryAgentBySubDomain(subDomain string) (*Agent, error) {
	agent := new(Agent)
	tx := mysql.GetGlobalDBIns().Where("sub_domain = ?", subDomain).First(agent)
	if tx.RowsAffected != 1 {
		return nil, fmt.Errorf("sub_domain:%v not exist", subDomain)
	}
	return agent, nil
}
