package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zachturing/login/common/define"
	"github.com/zachturing/login/common/xhttp"
	"github.com/zachturing/login/model"
	"github.com/zachturing/login/util"
	"github.com/zachturing/util/database/mysql"
	"github.com/zachturing/util/database/redis"
	"github.com/zachturing/util/log"
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
}

func LoginPhone(c *gin.Context) {
	var param phoneParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		log.Errorf("login phone: invalid param:%v, err:%v", param, err)
		xhttp.ParamsError(c, fmt.Errorf("login phone: invalid param:%v, err:%v", param, err))
		return
	}

	smsCode := redis.GetGlobalClient().Get(context.TODO(), smsKey(param.Phone)).Val()
	if smsCode != param.SMSCode {
		log.Errorf("login phone: %v, origin sms code not match %v->%v", param.Phone, param.SMSCode, smsCode)
		xhttp.ParamsError(c, fmt.Errorf("invalid sms code"))
		return
	}

	user, err := model.QueryUser(param.Phone)
	if err == nil { // 查到，则登录成功，直接返回
		token, _ := util.GenerateToken(int(user.ID))
		log.Debugf("phone login: user:%v success", user.Phone)
		xhttp.Data(c, map[string]string{
			"token": token,
		})
		return
	}

	userID, err := registerUser(param)
	if err != nil {
		log.Errorf("phone login: register user %v failed, err:%v", param.Phone, err)
		xhttp.ServerError(c, define.RegisterFailed, define.MapCodeToMsg[define.RegisterFailed])
		return
	}
	token, _ := util.GenerateToken(userID)
	log.Debugf("phone login: register success, phone:%v, userID:%v", param.Phone, userID)
	xhttp.Data(c, map[string]string{
		"token": token,
	})
	return
}

// registerUser 注册成功返回user_id
func registerUser(param phoneParam) (int, error) {
	var user model.User
	err := mysql.GetGlobalDBIns().Transaction(func(tx *gorm.DB) error {
		user = model.User{
			Phone:            param.Phone,
			RegistrationTime: time.Now(),
			LastLoginTime:    time.Now(),
			Role:             define.ROLE_NORMAL,
			SubDomain:        param.SubDomain, // 用户从哪个二级域名注册的，就绑定在哪个二级域名下，目前会根据二级域名区分代理商，www则为主域名
		}

		return model.CreateUser(&user)
	})
	if err != nil {
		log.Errorf("register user %v failed, err:%v", param.Phone, err)
		return 0, err
	}

	return int(user.ID), nil
}
