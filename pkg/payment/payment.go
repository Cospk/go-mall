package payment

import (
	"context"
	"time"
)

// PaymentRequest 支付请求参数
type PaymentRequest struct {
	OrderID     string    `json:"order_id"`     // 订单ID
	TotalAmount float64   `json:"total_amount"` // 支付金额
	Subject     string    `json:"subject"`      // 订单标题
	Body        string    `json:"body"`         // 订单描述
	ClientIP    string    `json:"client_ip"`    // 客户端IP
	NotifyURL   string    `json:"notify_url"`   // 支付结果通知URL
	ReturnURL   string    `json:"return_url"`   // 支付完成后跳转URL
	TradeType   string    `json:"trade_type"`   // 交易类型(如: APP, JSAPI, NATIVE等)
	OpenID      string    `json:"open_id"`      // 用户标识(微信JSAPI支付必填)
	TimeExpire  time.Time `json:"time_expire"`  // 订单失效时间
	ProductID   string    `json:"product_id"`   // 商品ID
	Attach      string    `json:"attach"`       // 附加数据
}

// PaymentResponse 支付响应参数
type PaymentResponse struct {
	TradeNo     string                 `json:"trade_no"`     // 支付平台交易号
	OutTradeNo  string                 `json:"out_trade_no"` // 商户订单号
	PaymentURL  string                 `json:"payment_url"`  // 支付链接(如果有)
	CodeURL     string                 `json:"code_url"`     // 二维码链接(如果有)
	PrepayID    string                 `json:"prepay_id"`    // 预支付ID(微信支付)
	PaySign     string                 `json:"pay_sign"`     // 支付签名(微信支付)
	NonceStr    string                 `json:"nonce_str"`    // 随机字符串
	PaymentForm string                 `json:"payment_form"` // 支付表单(支付宝网页支付)
	PaymentData map[string]interface{} `json:"payment_data"` // 其他支付数据
}

// PaymentNotify 支付通知参数
type PaymentNotify struct {
	TradeNo     string                 `json:"trade_no"`     // 支付平台交易号
	OutTradeNo  string                 `json:"out_trade_no"` // 商户订单号
	TotalAmount float64                `json:"total_amount"` // 支付金额
	PayTime     time.Time              `json:"pay_time"`     // 支付时间
	TradeStatus string                 `json:"trade_status"` // 交易状态
	BuyerID     string                 `json:"buyer_id"`     // 买家ID
	RawData     map[string]interface{} `json:"raw_data"`     // 原始通知数据
}

// QueryResponse 查询响应参数
type QueryResponse struct {
	TradeNo     string                 `json:"trade_no"`     // 支付平台交易号
	OutTradeNo  string                 `json:"out_trade_no"` // 商户订单号
	TotalAmount float64                `json:"total_amount"` // 支付金额
	TradeStatus string                 `json:"trade_status"` // 交易状态
	PayTime     time.Time              `json:"pay_time"`     // 支付时间
	BuyerID     string                 `json:"buyer_id"`     // 买家ID
	RawData     map[string]interface{} `json:"raw_data"`     // 原始查询数据
}

// RefundRequest 退款请求参数
type RefundRequest struct {
	OutTradeNo   string  `json:"out_trade_no"`  // 商户订单号
	TradeNo      string  `json:"trade_no"`      // 支付平台交易号
	OutRefundNo  string  `json:"out_refund_no"` // 商户退款单号
	TotalAmount  float64 `json:"total_amount"`  // 订单总金额
	RefundAmount float64 `json:"refund_amount"` // 退款金额
	RefundReason string  `json:"refund_reason"` // 退款原因
	NotifyURL    string  `json:"notify_url"`    // 退款结果通知URL
}

// RefundResponse 退款响应参数
type RefundResponse struct {
	OutTradeNo   string                 `json:"out_trade_no"`  // 商户订单号
	TradeNo      string                 `json:"trade_no"`      // 支付平台交易号
	OutRefundNo  string                 `json:"out_refund_no"` // 商户退款单号
	RefundID     string                 `json:"refund_id"`     // 支付平台退款单号
	RefundAmount float64                `json:"refund_amount"` // 退款金额
	RawData      map[string]interface{} `json:"raw_data"`      // 原始退款数据
}

// Payment 支付接口
type Payment interface {
	// Pay 创建支付订单
	Pay(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error)

	// ParseNotify 解析支付结果通知
	ParseNotify(ctx context.Context, data []byte) (*PaymentNotify, error)

	// Query 查询支付订单
	Query(ctx context.Context, outTradeNo, tradeNo string) (*QueryResponse, error)

	// Refund 申请退款
	Refund(ctx context.Context, req *RefundRequest) (*RefundResponse, error)

	// QueryRefund 查询退款
	QueryRefund(ctx context.Context, outRefundNo, refundID string) (*RefundResponse, error)

	// CloseOrder 关闭订单
	CloseOrder(ctx context.Context, outTradeNo string) error
}
