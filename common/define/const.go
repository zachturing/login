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

// 用户权限：普通用户、代理商、超管用户
const (
	PERMISSON_NORMAL      = "NORMAL"      // 普通用户，注册默认是这个角色
	PERMISSON_AGENT       = "AGENT"       // 代理商，申请之后会有审核
	PERMISSON_SUPER_ADMIN = "SUPER_ADMIN" // 超级管理员
)

// 用户级别，用于区分用户是普通用户、VIP用户，会在生成论文上走不通的算法模型
const (
	LEVEL_ANONYMOUS = "LEVEL_ANONYMOUS" // 目前已弃用，系统已经不支持未登录的用户生成大纲了
	LEVEL_NORMAL    = "LEVEL_NORMAL"
	LEVEL_VIP       = "LEVEL_VIP"
)
