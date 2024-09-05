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
