package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是返回给前端的标准化 JSON 结构体。
// 使用 `json:"..."` 标签来定义在序列化为 JSON 时每个字段的键名。
type Response struct {
	Code    int         `json:"code"`    // 业务状态码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

// 定义常用业务状态吗
const (
	CodeSuccess      = 200 // 成功
	CodeError        = 500 // 通用错误
	CodeUnauthorized = 401 // 认证失败
)

// result 是一个内部辅助函数，用于构造并发送 JSON 响应。
// 它接受业务状态码、消息、数据以及 Gin 上下文作为参数。
func result(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// Success 函数用于返回一个表示成功的响应。
// 它封装了成功的状态码和默认消息。
// data 参数是需要返回给客户端的具体数据。
func Success(data interface{}, c *gin.Context) {
	result(CodeSuccess, "success", data, c)
}

// Error 函数用于返回一个表示失败的响应。
// 它封装了失败的状态码和错误消息。
func Error(msg string, c *gin.Context) {
	result(CodeError, msg, nil, c)
}

// ErrorWithData 函数用于返回一个表示失败的响应，但可以携带一些额外数据。
// 例如，在表单验证失败时，可以用 data 字段返回具体的字段错误信息。
func ErrorWithData(msg string, data interface{}, c *gin.Context) {
	result(CodeError, msg, data, c)
}

// Unauthorized 函数用于返回一个表示认证失败的响应。
// 我们遵循项目现有的约定，HTTP 状态码依然返回 200，由前端根据 JSON 中的 code 字段来处理逻辑。
func Unauthorized(msg string, c *gin.Context) {
	result(CodeUnauthorized, msg, nil, c)
}
