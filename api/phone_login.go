package api

import (
	"context"
	"github.com/newdee/aipaper-util/config"
	"github.com/shopspring/decimal"
	"github.com/zachturing/login/model/wallet_model"
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

		// 用户注册
		user = model.User{
			Phone:            param.Phone,
			RegistrationTime: time.Now(),
			LastLoginTime:    time.Now(),
			Role:             define.LEVEL_NORMAL,
			Permission:       define.PERMISSON_NORMAL,
			AgentId:          int(agent.ID),
		}
		if err = model.CreateUser(&user, tx); err != nil {
			return err
		}

		// 创建积分账户
		account := wallet_model.Account{
			UserID:    int(user.ID),
			Currency:  define.CurrencyCNY,
			Balance:   decimal.NewFromInt(0),
			Status:    define.AccountStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err = wallet_model.CreateAccount(&account, tx); err != nil {
			return err
		}

		// 注册成功之后为该用户绑定一个唯一的邀请码
		user.InvCode = util.GenerateInvCodeByUserId(uint64(user.ID))
		if err = model.UpdateUserInvCode(&user, tx); err != nil {
			return err
		}

		// 赠送一次免费PaperYY查重
		userRights := model.UserRights{
			UserId:             user.ID,
			InvUsers:           0,
			DuplicateCheckNums: 1,
			UsedCheckNums:      0,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		if err = model.CreateUserRights(&userRights, tx); err != nil {
			return err
		}

		// 如果通过别人的邀请码链接进入，则为邀请人赠送一次查重权益
		if param.InvCode != "" {
			invUserId, err := util.DecodeInvCodeToUID(param.InvCode)
			if err != nil {
				return err
			}

			// 查询邀请人的权益
			invUserRights, err := model.QueryUserRightsByUserId(int(invUserId), tx)
			if err == nil {
				invUserRights.InvUsers += 1
				invUserRights.DuplicateCheckNums += 1
			}

			// 如果没有查到邀请人的权益记录，则插入一条新的权益记录
			if invUserRights == nil {
				invUserRights = &model.UserRights{
					UserId:             int64(invUserId),
					InvUsers:           1,
					DuplicateCheckNums: 1,
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
				}
			}
			return model.SaveUserRights(invUserRights, tx)
		}

		// 新用户注册赠送10次降AIGC次数
		if err = model.GiftUserRights(user.ID, tx); err != nil {
			return err
		}

		// 写入UserReduceLogs
		if err = model.InsertUserReduceLogs(user.ID,
			0,
			3,
			3,
			define.ChangeReasonGift,
			"",
			"",
			tx); err != nil {
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
