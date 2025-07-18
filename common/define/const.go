package define

import (
	"os"
	"time"

	"github.com/newdee/aipaper-util/config"
	"github.com/newdee/aipaper-util/log"
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

// 用户邀请码
const (
	PRIME1               = 3         // 与56互质
	PRIME2               = 5         // 与邀请码长度 6 互质
	SALT                 = 123456789 // 随意一个数值
	DefaultInvCodeLength = 6         // 邀请码长度，默认为6，请勿改动，邀请码的实现算法考虑到了互质的关系
)

const (
	BaiduAPIURL   = "https://ocpc.baidu.com/ocpcapi/api/uploadConvertData"
	BaiduApiToken = "0FKgl3lVEe0FCnVO7E8J4fG19jO1kH1t@jMcYhBnhQcwjxInk2Acvqhp1B9fF8ewz"
)

// UserReduceLogs表的change_reason
const (
	ChangeReasonUserOperation = "用户执行降AIGC" // 用户执行降AIGC
	ChangeReasonGift          = "后台赠送权益次数"  // 后台赠送，用于活动推广等
	ChangeReasonUserRecharge  = "用户付费购买次数"  // 用于用户购买降AIGC套餐的场景
)
