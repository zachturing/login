package define

const (
	OK               = 200
	ErrInvalidParams = 10001
	ErrServer        = 10002

	// InvalidPhone 无效的手机号码
	InvalidPhone = 40000
	// SendSMSCodeFailed 发送验证码失败
	SendSMSCodeFailed = 40001
	// SMSCodeInvalid 验证码无效
	SMSCodeInvalid = 40002

	SMSCodeSetRedisFailed = 40003
)

// MapCodeToMsg 返回码对应信息
var MapCodeToMsg = map[int]string{
	OK:               "success",
	ErrInvalidParams: "invalid params",
	ErrServer:        "server error",

	// 新增错误码及对应的信息
	InvalidPhone:          "invalid phone number",
	SendSMSCodeFailed:     "failed to send SMS code",
	SMSCodeInvalid:        "invalid SMS code",
	SMSCodeSetRedisFailed: "failed to set SMS code to redis",
}

const (
	MsgInvalidParams = "参数错误"
)
