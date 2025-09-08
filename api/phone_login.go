package api

import (
	"context"
	"github.com/newdee/aipaper-util/config"
	"github.com/zachturing/login/redis"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/newdee/aipaper-util/database/mysql"
	"github.com/newdee/aipaper-util/log"
	"github.com/zachturing/login/common/define"
	"github.com/zachturing/login/common/xhttp"
	apolloConfig "github.com/zachturing/login/config"
	"github.com/zachturing/login/model"
	"github.com/zachturing/login/util"
	"gorm.io/gorm"
)

// 定义phoneParam结构体
type phoneParam struct {
	// 用户手机号
	Phone string `json:"phone" validate:"required,len=11"`
	// 短信验证码
	SMSCode string `json:"sms_code" validate:"required,len=6"`
	// 用户的代理商，使用域名代替
	//SubDomain string `json:"sub_domain"` // 弃用
	// 邀请码，可以为空
	InvCode string `json:"inv_code"`
	// 百度的数据回传接口需要的code
	BdVid string `json:"bd_vid"`
	// 密码
	Password string `json:"password"`
}

type loginResponse struct {
	Token            string `json:"token"`
	ExpiredTimestamp int64  `json:"expired_timestamp"`
}

func LoginPhone(c *gin.Context) {
	var param phoneParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		log.Errorf("login phone: invalid param:%v, err:%v", param, err)
		xhttp.DiyOkCode(c, define.ErrInvalidParams, define.MapCodeToMsg[define.ErrInvalidParams])
		return
	}

	if param.Password != "" {
		// TODO：校验用户密码，后续加，默认一个密码内部能登录所有用户的账户
		if param.Password != "mixpaer@wandou" {
			log.Errorf("login phone: %s, password: %s not match", param.Phone, param.Password)
			xhttp.DiyOkCode(c, define.PasswordInvalid, define.MapCodeToMsg[define.PasswordInvalid])
			return
		}
	} else {
		smsCode := redis.GetGlobalClient().Get(context.TODO(), smsKey(param.Phone)).Val()
		if smsCode != param.SMSCode {
			log.Errorf("login phone: %v, origin sms code not match %v->%v", param.Phone, param.SMSCode, smsCode)
			xhttp.DiyOkCode(c, define.SMSCodeInvalid, define.MapCodeToMsg[define.SMSCodeInvalid])
			return
		}
	}

	// 用户已存在，登录成功，直接返回
	user, err := model.QueryUser(param.Phone)
	if err == nil {
		token, expiredTimeStamp, _ := util.GenerateToken(int(user.ID))
		log.Debugf("phone login: user:%v success", user.Phone)

		// 更新用户的登录时间
		updatedColumns := map[string]interface{}{
			"last_login_time": time.Now(),
		}
		if err = model.UpdateUserColumns(user.ID, updatedColumns, nil); err != nil {
			log.Errorf("phone login: update user %v failed, err:%v", user.Phone, err)
		}
		xhttp.Data(c, loginResponse{
			Token:            token,
			ExpiredTimestamp: expiredTimeStamp,
		})
		return
	}

	// 非生产环境禁止注册用户
	if define.Env != config.ProdEnv {
		phoneArray, err := apolloConfig.GetAllowRegistryPhone()
		if err != nil {
			xhttp.DiyOkCode(c, define.RegisterFailed, "非生产环境禁止新用户注册！")
			return
		}
		if !util.ContainsStr(phoneArray, param.Phone) {
			xhttp.DiyOkCode(c, define.RegisterFailed, "非生产环境禁止新用户注册！")
			return
		}
	}

	// 用户注册
	userID, err := registerUser(param)
	if err != nil {
		log.Errorf("phone login: register user %v failed, err:%v", param.Phone, err)
		xhttp.DiyOkCode(c, define.RegisterFailed, define.MapCodeToMsg[define.RegisterFailed])
		return
	}

	// 注册完成调用百度推广的数据回传接口，调用失败记录日志，不要抛错
	err = CallBaiduUploadConvertData(param.BdVid)
	if err != nil {
		log.Errorf("user register: call baidu upload convert data %v failed, err:%v", param.BdVid, err)
	}

	// 生成token
	token, expiredTimeStamp, _ := util.GenerateToken(userID)
	log.Debugf("phone login: register success, phone:%v, userID:%v", param.Phone, userID)
	xhttp.Data(c, loginResponse{
		Token:            token,
		ExpiredTimestamp: expiredTimeStamp,
	})
	return
}

// registerUser 注册成功返回user_id
func registerUser(param phoneParam) (int, error) {
	var userId int
	err := mysql.GetGlobalDBIns().Transaction(func(tx *gorm.DB) error {
		// 新用户注册
		user, err := registerNewUser(tx, param)
		if err != nil {
			return err
		}

		// 为新用户初始化分销商账户
		if err = initDistributionAccount(tx, user); err != nil {
			return err
		}

		// 为新用户初始化降AIGC权益记录
		if err = initUserReduceRights(tx, user); err != nil {
			return err
		}

		// 处理邀请双方权益
		if user.ParentUserId != 0 {
			if err = processInviteRights(tx, user); err != nil {
				return err
			}
		}

		// 返回参数
		userId = int(user.ID)
		return nil
	})

	if err != nil {
		log.Errorf("register user %v failed, err:%v", param.Phone, err)
		return 0, err
	}

	return userId, nil
}

// registerNewUser 注册新用户
func registerNewUser(tx *gorm.DB, param phoneParam) (*model.User, error) {
	user := &model.User{
		Phone:            param.Phone,
		UserName:         util.GenerateUserName(param.Phone),
		RegistrationTime: time.Now(),
		LastLoginTime:    time.Now(),
		Role:             define.LevelNormal,
		Permission:       define.PermissionNormal,
	}
	// 根据入参中的邀请码解码获取邀请人ID
	if inviteUserId, decodeError := util.DecodeInvCodeToUID(param.InvCode); decodeError == nil {
		user.ParentUserId = int64(inviteUserId)
	}
	// 创建用户
	if err := model.CreateUser(user, tx); err != nil {
		return nil, err
	}
	// 生成当前用户的邀请码
	user.InvCode = util.GenerateInvCodeByUserId(uint64(user.ID))
	if err := model.UpdateUserColumns(user.ID,
		map[string]interface{}{
			"inv_code": user.InvCode,
		},
		tx); err != nil {
		return nil, err
	}
	return user, nil
}

// 为新用户初始化分销商账户
func initDistributionAccount(tx *gorm.DB, user *model.User) error {
	var agentAccount = model.DistributionAccount{
		UserID:             user.ID,
		Currency:           define.CurrencyCny,
		Status:             define.DistributionAccountStatusNormal,
		Balance:            0.0,
		FrozenAmount:       0.0,
		WithdrawnAmount:    0.0,
		TotalIncome:        0.0,
		DirectPercent:      0.2,
		IndirectPercent:    0.0,
		UserUpgradePercent: 0.0,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if err := model.CreateDistributionAccount(&agentAccount, tx); err != nil {
		return err
	}
	return nil
}

// 为新用户初始化降AIGC权益记录
func initUserReduceRights(tx *gorm.DB, user *model.User) error {
	var userReduceRights = model.UserReduceRights{
		UserID:        user.ID,
		RemainingNum:  0,
		UsedReduceNum: 0,
	}
	if err := model.CreateUserReduceRights(&userReduceRights, tx); err != nil {
		return err
	}
	return nil
}

// 处理邀请权益
func processInviteRights(tx *gorm.DB, user *model.User) error {
	// 1.1 赠送被邀请人 10次降AIGC次数
	if err := processUserReduceRights(user.ID, 10, tx); err != nil {
		return err
	}
	// 1.2 赠送邀请人 10次降AIGC次数
	if err := processUserReduceRights(user.ParentUserId, 10, tx); err != nil {
		return err
	}
	// 2、被邀请人赠送9折优惠券
	couponCode := &model.Coupon{
		Type:           define.CouponTypeDiscount,
		RuleId:         -1,
		CreateUserId:   int(user.ID),
		ExchangeUserId: user.ID,
		CouponCode:     util.BuildCouponCode(define.CouponChannelRegistry, strconv.Itoa(define.CouponTypeDiscount)),
		DiscountRate:   0.9,
		RightsNum:      0,
		CreateTime:     time.Now(),
		UsedTime:       nil,
		ExpireTime:     time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC),
		Channel:        define.CouponChannelRegistry,
		Status:         define.CouponStatusExchanged,
	}
	if err := tx.Create(couponCode).Error; err != nil {
		return err
	}

	// 3、邀请人基于梯度获取优惠券奖励：3人-8折，10人-7折，20人-6折，30人-5折
	invitedCounts, err := model.CountInvitedUsers(user.ParentUserId, tx)
	if err != nil {
		return err
	}
	discountRate := 0.0
	inviterRewards := "降AIGC10次" // 邀请人奖励
	if invitedCounts == 1 {
		discountRate = 0.9
		inviterRewards = "降AIGC10次、9折正文优惠券"
	} else if invitedCounts == 3 {
		discountRate = 0.8
		inviterRewards = "降AIGC10次、8折正文优惠券"
	} else if invitedCounts == 10 {
		discountRate = 0.7
		inviterRewards = "降AIGC10次、7折正文优惠券"
	} else if invitedCounts == 20 {
		discountRate = 0.6
		inviterRewards = "降AIGC10次、6折正文优惠券"
	} else if invitedCounts == 30 {
		discountRate = 0.5
		inviterRewards = "降AIGC10次、5折正文优惠券"
	}
	if discountRate > 0.0 {
		inviterCoupon := &model.Coupon{
			Type:           define.CouponTypeDiscount,
			RuleId:         -1,
			CreateUserId:   int(user.ID),
			ExchangeUserId: user.ID,
			CouponCode:     util.BuildCouponCode(define.CouponChannelInvite, strconv.Itoa(define.CouponTypeDiscount)),
			DiscountRate:   discountRate,
			RightsNum:      0,
			CreateTime:     time.Now(),
			UsedTime:       nil,
			ExpireTime:     time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC),
		}
		if err = tx.Create(inviterCoupon).Error; err != nil {
			return err
		}
	}

	// 4、 生成邀请记录
	invitationLogs := model.InvitationLogs{
		InviterId:      user.ParentUserId,
		InviteeId:      user.ID,
		InviteeName:    user.UserName,
		InviterRewards: inviterRewards,
		InviteeRewards: "降AIGC10次、9折正文优惠券",
		Remarks:        "",
	}
	if err = model.CreateInvitationLogs(&invitationLogs, tx); err != nil {
		return err
	}
	return nil
}

func processUserReduceRights(userId int64, rightsNums int, tx *gorm.DB) error {
	// 查询用户权益
	userReduceRights, err := model.QueryUserReduceRights(userId, tx)
	if err != nil {
		return err
	}
	// 更新用户权益
	preRightsNum := userReduceRights.RemainingNum
	postRightsNum := preRightsNum + rightsNums
	if err = model.UpdateUserReduceRightsColumnsByUserId(userId, map[string]interface{}{
		"remaining_num": postRightsNum,
	}, tx); err != nil {
		return err
	}
	// 记录权益日志
	userReduceRightsLog := &model.UserReduceLogs{
		UserID:           userId,
		PreReduceNum:     preRightsNum,
		ChangeNum:        rightsNums,
		PostReduceNum:    postRightsNum,
		ChangeReason:     define.ChangeReasonRegistry,
		OriginalContents: "",
		PostContents:     "",
	}
	if err = tx.Create(userReduceRightsLog).Error; err != nil {
		return err
	}
	return nil
}
