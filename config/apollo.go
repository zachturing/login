package config

import "github.com/newdee/aipaper-util/config"

type TencentSMSConfig struct {
	Enable           bool   `json:"enable"`
	TencentSecretId  string `json:"tencent_secret_id"`
	TencentSecretKey string `json:"tencent_secret_key"`
	SdkAppId         string `json:"sdk_app_id"`
	TemplateId       string `json:"template_id"`
	SignName         string `json:"sign_name"`
}

func GetTencentSMSConfig() (*TencentSMSConfig, error) {
	cfg, err := config.Get(config.Common)
	if err != nil {
		return nil, err
	}
	var param TencentSMSConfig
	err = cfg.GetWithUnmarshal("tencent_sms", &param, &config.JSONUnmarshaler{})
	if err != nil {
		return nil, err
	}
	return &param, nil
}
