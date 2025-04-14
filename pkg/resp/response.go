package resp

import (
	errcode2 "github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/gin-gonic/gin"
)

type response struct {
	ctx       *gin.Context
	Code      int         `json:"Code"`
	Msg       string      `json:"msg"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data,omitempty"`
	PageInfo  *PageInfo   `json:"page_info,omitempty"`
}

// NewResponse 构造一个响应,根据需要调用其他方法
func NewResponse(ctx *gin.Context) *response {
	return &response{ctx: ctx}
}

func (r *response) SetPageInfo(pageInfo *PageInfo) *response {
	r.PageInfo = pageInfo
	return r
}

// Success 成功并给出数据的响应
func (r *response) Success(data interface{}) {
	r.Code = errcode2.Success.Code()
	r.Msg = errcode2.Success.Msg()
	requestId := ""
	if _, exists := r.ctx.Get("Trace-Id"); exists {
		requestId = r.ctx.GetString("Trace-Id")
	}
	r.RequestId = requestId
	r.Data = data
	r.ctx.JSON(errcode2.Success.HttpStatusCode(), r)
}

func (r *response) SuccessOk() {
	r.Success("")
}

// Error 失败并给出错误的响应
func (r *response) Error(err *errcode2.AppError) {
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
