package api

import (
	"context"
	"github.com/newdee/aipaper-util/config"
	"github.com/zachturing/login/redis"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/newdee/aipaper-util/database/mysql"
	"github.com/newdee/aipaper-util/log"
	"github.com/zachturing/login/common/define"
	"github.com/zachturing/login/common/xhttp"
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
	SubDomain string `json:"sub_domain"`
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
		xhttp.DiyOkCode(c, define.RegisterFailed, "非生产环境禁止新用户注册！")
		return
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
	var user model.User
	err := mysql.GetGlobalDBIns().Transaction(func(tx *gorm.DB) error {
		// 根据域名查询对应的代理商
		agent, err := model.QueryAgentBySubDomain(param.SubDomain, tx)
		if err != nil {
			return err
		}

		// 新用户注册
		user = model.User{
			Phone:            param.Phone,
			UserName:         util.GenerateUserName(param.Phone),
			RegistrationTime: time.Now(),
			LastLoginTime:    time.Now(),
			Role:             define.LEVEL_NORMAL,
			Permission:       define.PERMISSON_NORMAL,
			AgentId:          int(agent.ID),
		}
		// 根据入参中的邀请码解码获取邀请人ID
		if inviteUserId, decodeError := util.DecodeInvCodeToUID(param.InvCode); decodeError == nil {
			user.ParentUserId = int64(inviteUserId)
		}
		if err = model.CreateUser(&user, tx); err != nil {
			return err
		}

		// 创建分销商账户
		var agentAccount = model.DistributionAccount{
			UserID:             user.ID,
			Currency:           define.CurrencyCny,
			Status:             define.AccountStatusNormal,
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
		if err = model.CreateDistributionAccount(&agentAccount, tx); err != nil {
			return err
		}

		// 生成当前用户的邀请码
		updatedColumns := map[string]interface{}{
			"inv_code": util.GenerateInvCodeByUserId(uint64(user.ID)),
		}
		if err = model.UpdateUserColumns(user.ID, updatedColumns, tx); err != nil {
			return err
		}

		// 如果用户是被邀请的，则生成邀请记录
		if user.ParentUserId != 0 {
			invitationLogs := model.InvitationLogs{
				InviterId:          user.ParentUserId,
				InviteeId:          user.ID,
				InviterRewardsType: "待定义",
				InviteeRewardsType: "待定义",
				Remarks:            "",
			}
			if err = model.CreateInvitationLogs(&invitationLogs, tx); err != nil {
				return err
			}
		}

		// TODO：基于配置项赠送权益

		// 新用户注册生成空的降AIGC权益记录
		var userReduceRights = model.UserReduceRights{
			UserID:        user.ID,
			RemainingNum:  0,
			UsedReduceNum: 0,
		}
		if err = model.CreateUserReduceRights(&userReduceRights, tx); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Errorf("register user %v failed, err:%v", param.Phone, err)
		return 0, err
	}

	return int(user.ID), nil
}
