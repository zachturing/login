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
	PermissionNormal       = "NORMAL"       // 普通用户，注册默认是这个角色
	PermissionAgent        = "AGENT"        // 域名代理商，弃用
	PermissionDistribution = "DISTRIBUTION" // 分销商
	PermissionFranchisee   = "FRANCHISEE"   // 加盟商
	PermissionSuperAdmin   = "SUPER_ADMIN"  // 超级管理员
)

// 用户级别，用于区分用户是普通用户、VIP用户，会在生成论文上走不通的算法模型
const (
	LevelAnonymous = "LEVEL_ANONYMOUS" // 目前已弃用，系统已经不支持未登录的用户生成大纲了
	LevelNormal    = "LEVEL_NORMAL"
	LevelVip       = "LEVEL_VIP"
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

// 分销商账户相关常量
const (
	CurrencyCny = "CNY"

	DistributionAccountStatusPendingActivation = "PENDING_ACTIVATION" // 分销商账户状态：待激活
	DistributionAccountStatusNormal            = "NORMAL"             // 分销商账户状态：正常
	DistributionAccountStatusFrozen            = "FROZEN"             // 分销商账户状态：冻结
	DistributionAccountStatusClosed            = "CLOSED"             // 分销商账户状态：关闭
)

// UserReduceLogs表的change_reason
const (
	ChangeReasonUserOperation = "用户执行降AIGC" // 用户执行降AIGC
	ChangeReasonGift          = "后台赠送权益次数"  // 后台赠送，用于活动推广等
	ChangeReasonRegistry      = "注册赠送权益次数"  // 通过邀请码注册赠送
	ChangeReasonUserRecharge  = "用户付费购买次数"  // 用于用户购买降AIGC套餐的场景
	ChangeReasonUserCDK       = "用户兑换卡密"    // 用于用户购买降AIGC套餐的场景
	ChangeReasonRefund        = "因退款退回权益"   // 用户用户退款场景
)

// 优惠券常量
const (
	CouponTypeDiscount      = 1 // 折扣券
	CouponTypeFullReduction = 2 // 满减券
	CouponTypeAIGC          = 3 // 降AIGC次数券\

	CouponChannelTB       = "tb"       // 淘宝
	CouponChannelBilibili = "bili"     // B站
	CouponChannelDy       = "dy"       // 抖音
	CouponChannelXhs      = "xhs"      // 小红书
	CouponChannelWx1      = "wx1"      // 微商
	CouponChannelWx2      = "wx2"      // 微信视频号
	CouponChannelUserBuy  = "buy"      // 用户购入
	CouponChannelRegistry = "registry" // 注册赠送
	CouponChannelInvite   = "invite"   // 邀请赠送
	CouponChannelOther    = "other"    // 其他

	CouponStatusUnused    = 1 // 未使用
	CouponStatusUsed      = 2 // 已使用
	CouponStatusExpired   = 3 // 已过期
	CouponStatusFrozen    = 4 // 已冻结
	CouponStatusInvalid   = 5 // 已作废
	CouponStatusExchanged = 6 // 已兑换
)
