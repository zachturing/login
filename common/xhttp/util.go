package xhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zachturing/login/common/define"
)


// BaseResponse 基础response
type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

// OK response success
func OK(c *gin.Context) {
	c.JSON(http.StatusOK, &BaseResponse{
		Code:    define.OK,
		Message: define.MapCodeToMsg[define.OK],
		Result:  nil,
	})
}

// Data 返回指定数据response
func Data(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &BaseResponse{
		Code:    define.OK,
		Message: define.MapCodeToMsg[define.OK],
		Result:  data,
	})
}

//
// 一些常见的错误封装
//

// ServerError 服务端错误
func ServerError(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusInternalServerError, &BaseResponse{
		Code:    code,
		Message: msg,
		Result:  nil,
	})
}

// ClientError 客户端错误
func ClientError(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusBadRequest, &BaseResponse{
		Code:    code,
		Message: msg,
		Result:  nil,
	})
}

// ParamsError 解析参数错误
func ParamsError(c *gin.Context, err error) {
	ClientError(c, define.ErrInvalidParams, define.MsgInvalidParams+":"+ err.Error())
}
