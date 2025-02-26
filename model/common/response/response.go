package response

import (
	"github.com/Cospk/go-mall/common/errcode"
	logger "github.com/Cospk/go-mall/global"
	"github.com/gin-gonic/gin"
)

type response struct {
	ctx       *gin.Context
	Code      int         `json:"Code"`
	Msg       string      `json:"msg"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data,omitempty"`
}

// NewResponse 构造一个响应,根据需要调用其他方法
func NewResponse(ctx *gin.Context) *response {
	return &response{ctx: ctx}
}

// Success 成功并给出数据的响应
func (r *response) Success(data interface{}) {
	r.Code = errcode.Success.Code()
	r.Msg = errcode.Success.Msg()
	requestId := ""
	if _, exists := r.ctx.Get("Trace-Id"); exists {
		requestId = r.ctx.GetString("Trace-Id")
	}
	r.RequestId = requestId
	r.Data = data
	r.ctx.JSON(errcode.Success.HttpStatusCode(), r)
}

func (r *response) SuccessOk() {
	r.Success("")
}

// Fail 失败并给出错误的响应
func (r *response) Fail(err *errcode.AppError) {
	r.Code = err.Code()
	r.Msg = err.Msg()
	requestId := ""
	if _, exists := r.ctx.Get("Trace-Id"); exists {
		requestId = r.ctx.GetString("Trace-Id")
	}
	r.RequestId = requestId
	// Error记录到日志
	logger.NewLogger(r.ctx).Error("api_response_error", "err", err)
	r.ctx.JSON(err.HttpStatusCode(), r)
}
