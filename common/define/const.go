package define

import (
	"github.com/zachturing/util/config"
	"time"
)

const (
	Env = config.DevEnv
)

const (
	// SMSCodeExpiredTime 短信验证码过期时间
	SMSCodeExpiredTime = time.Minute * 3

	// TokenExpireTime token有效时间
	TokenExpireTime = time.Hour * 24 * 7
)
