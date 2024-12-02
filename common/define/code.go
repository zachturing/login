package define

const (
	OK               = 200
	ErrInvalidParams = 10001
	ErrServer        = 10002

	InvalidPhone      = 40000 // 无效的手机号码
	SendSMSCodeFailed = 40001 // 发送验证码失败
	SMSCodeInvalid    = 40002 // 验证码无效
	RegisterFailed    = 40004 // 注册/登录失败
)

// MapCodeToMsg 返回码对应信息
var MapCodeToMsg = map[int]string{
	OK:               "success",
	ErrInvalidParams: "参数错误",
	ErrServer:        "服务器内部错误",

	// 新增错误码及对应的信息
	InvalidPhone:      "无效的手机号",
	SendSMSCodeFailed: "发送验证码失败",
	SMSCodeInvalid:    "验证码无效",
	RegisterFailed:    "登录出错",
}

const (
	MsgInvalidParams = "参数错误"
)
