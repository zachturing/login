package define

const (
	OK               = 200
	ErrInvalidParams = 10001
	ErrServer        = 10002
)

// MapCodeToMsg 返回码对应信息
var MapCodeToMsg = map[int]string{
	OK:               "success",
	ErrInvalidParams: "invalid params",
	ErrServer:        "server error",
}

const (
	MsgInvalidParams = "参数错误"
)

