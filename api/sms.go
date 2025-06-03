package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/newdee/aipaper-util/log"
	"github.com/newdee/aipaper-util/sms"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	tencentsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/zachturing/login/common/define"
	"github.com/zachturing/login/common/xhttp"
	"github.com/zachturing/login/config"
	"github.com/zachturing/login/redis"
	"github.com/zachturing/login/util"
)

type smsParam struct {
	Phone string `json:"phone" validate:"required,len=11"`
}

func SendTencentSMS(c *gin.Context) {
	var param smsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Debugf("bind error")
	}

	if err := validator.New().Struct(param); err != nil {
		log.Errorf("invalid param:%v, err:%v", param, err)
		xhttp.DiyOkCode(c, define.ErrInvalidParams, define.MapCodeToMsg[define.ErrInvalidParams])
		return
	}

	// 从Apollo获取腾讯云的短信配置
	tencentCfg, err := config.GetTencentSMSConfig()
	if err != nil {
		log.Errorf("get tencent cloud config from apollo failed, err:%v", err)
		xhttp.DiyOkCode(c, define.SendSMSCodeFailed, define.MapCodeToMsg[define.SendSMSCodeFailed])
		return
	}

	// 生成smsCode
	smsCode := util.GenerateSMSCode()

	// 初始化腾讯云client
	credential := common.NewCredential(
		tencentCfg.TencentSecretId,
		tencentCfg.TencentSecretKey,
	)
	client, err := tencentsms.NewClient(credential, regions.Beijing, profile.NewClientProfile())
	if err != nil {
		log.Errorf("new tencent cloud client error, err:%v", err)
		xhttp.DiyOkCode(c, define.SendSMSCodeFailed, define.MapCodeToMsg[define.SendSMSCodeFailed])
		return
	}

	// 发送验证码
	smsReq := tencentsms.NewSendSmsRequest()
	smsReq.PhoneNumberSet = []*string{
		common.StringPtr(fmt.Sprintf("+86%s", param.Phone)),
	}
	smsReq.SignName = &tencentCfg.SignName
	smsReq.SmsSdkAppId = &tencentCfg.SdkAppId
	smsReq.TemplateId = &tencentCfg.TemplateId
	smsReq.TemplateParamSet = []*string{
		common.StringPtr(smsCode),
		common.StringPtr("3"),
	}
	if _, err = client.SendSms(smsReq); err != nil {
		log.Errorf("send sms to phone:%v failed, err:%v", param.Phone, err)
		xhttp.DiyOkCode(c, define.InvalidPhone, define.MapCodeToMsg[define.InvalidPhone])
		return
	}

	// set redis
	cmd := redis.GetGlobalClient().Set(context.TODO(), smsKey(param.Phone), smsCode, define.SMSCodeExpiredTime)
	if cmd.Err() != nil {
		log.Errorf("set redis key:%v failed, err:%v", smsKey(param.Phone), cmd.Err())
		xhttp.DiyOkCode(c, define.SendSMSCodeFailed, define.MapCodeToMsg[define.SendSMSCodeFailed])
		return
	}

	log.Debugf("send sms to phone:%v success", param.Phone)
	xhttp.OK(c)
}

func SendSMS(c *gin.Context) {
	var param smsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Debugf("bind error")
	}

	if err := validator.New().Struct(param); err != nil {
		log.Errorf("invalid param:%v, err:%v", param, err)
		xhttp.DiyOkCode(c, define.ErrInvalidParams, define.MapCodeToMsg[define.ErrInvalidParams])
		return
	}

	smsCode, err := sms.Send(param.Phone)
	if err != nil {
		log.Errorf("send sms to phone:%v failed, err:%v", param.Phone, err)
		xhttp.DiyOkCode(c, define.InvalidPhone, define.MapCodeToMsg[define.InvalidPhone])
		return
	}

	cmd := redis.GetGlobalClient().Set(context.TODO(), smsKey(param.Phone), smsCode, define.SMSCodeExpiredTime)
	if cmd.Err() != nil {
		log.Errorf("set redis key:%v failed, err:%v", smsKey(param.Phone), cmd.Err())
		xhttp.DiyOkCode(c, define.SendSMSCodeFailed, define.MapCodeToMsg[define.SendSMSCodeFailed])
		return
	}

	log.Debugf("send sms to phone:%v success", param.Phone)
	xhttp.OK(c)
}

func smsKey(phone string) string {
	return fmt.Sprintf("LOGIN:PHONE:%s", phone)
}
