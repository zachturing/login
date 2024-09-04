package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zachturing/login/common/define"
	"github.com/zachturing/login/common/xhttp"
	"github.com/zachturing/util/database/redis"
	"github.com/zachturing/util/log"
	"github.com/zachturing/util/sms"
)

type smsParam struct {
	Phone string `json:"phone" validate:"required,len=11"`
}

func SendSMS(c *gin.Context) {
	var param smsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Debugf("bind error")
	}

	if err := validator.New().Struct(param); err != nil {
		log.Errorf("invalid param:%v, err:%v", param, err)
		xhttp.ParamsError(c, err)
		return
	}

	smsCode, err := sms.Send(param.Phone)
	if err != nil {
		log.Errorf("send sms to phone:%v failed, err:%v", param.Phone, err)
		xhttp.ServerError(c, define.InvalidPhone, define.MapCodeToMsg[define.InvalidPhone])
		return
	}

	cmd := redis.GetGlobalClient().Set(context.TODO(), smsKey(param.Phone), smsCode, define.SMSCodeExpiredTime)
	if cmd.Err() != nil {
		xhttp.ServerError(c, define.SMSCodeSetRedisFailed, define.MapCodeToMsg[define.SMSCodeSetRedisFailed])
		return
	}

	log.Debugf("send sms to phone:%v success", param.Phone)
	xhttp.OK(c)
}

func smsKey(phone string) string {
	return fmt.Sprintf("LOGIN:PHONE:%s", phone)
}
