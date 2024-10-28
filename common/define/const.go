package define

import (
	"os"
	"time"

	"github.com/zachturing/util/config"
	"github.com/zachturing/util/log"
)

var (
	Env = config.DevEnv
)

func init() {
	env := os.Getenv("ENV")
	if len(env) > 0 {
		Env = config.EnvType(env)
	}
	log.Infof("ENV is %v", Env)
}

const (
	// SMSCodeExpiredTime 短信验证码过期时间
	SMSCodeExpiredTime = time.Minute * 3

	// TokenExpireTime token有效时间
	TokenExpireTime = time.Hour * 24 * 7
)

const (
	ROLE_NORMAL      = "LEVEL_NORMAL"      // 普通用户，注册默认是这个角色
	ROLE_PROXY       = "LEVEL_PROXY"       // 代理商，申请之后会有审核
	ROLE_SUPER_ADMIN = "LEVEL_SUPER_ADMIN" // 超级管理员
)
