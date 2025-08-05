package httpx

// 导入必要的包
import (
	"fmt"
	"github.com/gin-gonic/gin" // Gin Web 框架
	"net/http"
)

// RespErr 自定义错误响应结构体
// 实现了 error 接口，可以作为 error 类型处理
// 用于构建带有状态码、错误码和错误信息的响应
type RespErr struct {
	StatusCode int    // HTTP 状态码
	Code       int    // 业务错误码
	Message    string // 错误信息
}

func NewRespErr(code, statusCode int, message string) error {
	return &RespErr{Code: code, StatusCode: statusCode, Message: message}
}

// ApiError 实现 error 接口
// 返回格式为 "错误码: 错误信息"
// 例如: "500: internal server error"
func (err *RespErr) Error() string {
	return fmt.Sprintf("%d: %s", err.Code, err.Message)
}

// RespResult 通用响应结构体
// 支持泛型 T，可以用于各种类型的响应数据
// 用于统一 API 的返回格式（支持 JSON 和 XML）
type RespResult[T any] struct {
	Code int    `json:"code" xml:"code"` // 响应码
	Msg  string `json:"msg" xml:"msg"`   // 响应消息
	Data T      `json:"data" xml:"data"` // 响应数据
}

// ApiSuccess 返回成功响应
// 状态码固定为 200，消息固定为 "success"
// 参数 ctx: Gin 上下文
// 参数 data: 响应数据，可以是任何类型
func ApiSuccess[T any](ctx *gin.Context, data T) {
	ctx.JSON(http.StatusOK, RespResult[T]{
		Code: http.StatusOK,
		Msg:  "success",
		Data: data,
	})
}
func ApiNoAuth(ctx *gin.Context, err error) {
	ApiErrWithCode(ctx, err, http.StatusUnauthorized)
}

func ApiNoForbidden(ctx *gin.Context, err error) {
	ApiErrWithCode(ctx, err, http.StatusForbidden)
}

func ApiErrWithCode(ctx *gin.Context, err error, code int) {
	ApiError(ctx, NewRespErr(code, http.StatusOK, err.Error()))
}

// ApiError 返回失败响应
// 根据错误类型构建适当的响应
// 如果是 RespErr 类型，则使用其提供的 StatusCode、Code 和 Message
// 否则默认为 500 错误，msg 为 err.ApiError()
// 参数 ctx: Gin 上下文
// 参数 err: 错误信息
func ApiError(ctx *gin.Context, err error) {
	statusCode := http.StatusOK            // 默认 HTTP 状态码
	msg := err.Error()                     // 默认错误消息
	code := http.StatusInternalServerError // 默认错误码

	// 判断是否为自定义错误
	if respErr, ok := err.(*RespErr); ok {
		if respErr.StatusCode != 0 {
			statusCode = respErr.StatusCode // 使用自定义 HTTP 状态码
		}
		if respErr.Message != "" {
			msg = respErr.Message // 使用自定义错误消息
		}
		if respErr.Code != 0 {
			code = respErr.Code // 使用自定义错误码
		}
	}

	// 返回错误响应
	ctx.JSON(statusCode, RespResult[[]string]{
		Code: code,
		Msg:  msg,
		Data: []string{}, // 错误响应时数据为空数组
	})
}
