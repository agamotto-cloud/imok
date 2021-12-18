package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RespEntity struct {
	Success bool `json:"success"`
	//0是正確,其他需要參考接口錯誤碼對應
	Code int `json:"statusCode"`
	//响应消息
	Msg string `json:"msg"`
	//接口数据
	Data interface{} `json:"data"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, RespEntity{
		code == 0,
		code,
		msg,
		data,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "ok", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, nil, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "ok", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "ok", c)
}

func FailWithCodeMessage(code int, message string, c *gin.Context) {
	Result(code, nil, message, c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}
