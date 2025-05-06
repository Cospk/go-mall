package payment

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// PaymentService 支付服务
type PaymentService struct {
	orderRepo model.OrderRepository
}

// NewPaymentService 创建支付服务
func NewPaymentService(orderRepo model.OrderRepository) *PaymentService {
	return &PaymentService{
		orderRepo: orderRepo,
	}
}

// CreatePayment 创建支付
func (s *PaymentService) CreatePayment(ctx context.Context, orderID, paymentType, tradeType string) (map[string]interface{}, error) {
	// 获取订单信息
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	if order.Status != model.OrderStatusPending {
		return nil, errors.New("订单状态不允许支付")
	}

	// 获取支付实例
	p, err := payment.Get(paymentType)
	if err != nil {
		return nil, fmt.Errorf("不支持的支付方式: %w", err)
	}

	// 构建支付请求
	req := &payment.PaymentRequest{
		OrderID:     orderID,
		TotalAmount: order.TotalAmount,
		Subject:     fmt.Sprintf("订单 %s", orderID),
		Body:        fmt.Sprintf("商品购买 - 订单号: %s", orderID),
		ClientIP:    "127.0.0.1", // 应从请求中获取
		TradeType:   tradeType,
		TimeExpire:  time.Now().Add(2 * time.Hour), // 订单2小时后过期
	}

	// 创建支付
	resp, err := p.Pay(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("创建支付失败: %w", err)
	}

	// 更新订单支付信息
	order.PaymentType = paymentType
	order.PaymentTradeNo = resp.TradeNo
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("更新订单支付信息失败: %w", err)
	}

	// 返回支付数据
	result := map[string]interface{}{
		"order_id":     orderID,
		"payment_type": paymentType,
	}

	// 根据不同支付方式和交易类型，返回不同的支付数据
	switch paymentType {
	case payment.PaymentTypeWechat:
		switch tradeType {
		case "JSAPI":
			result["prepay_id"] = resp.PrepayID
			result["pay_sign"] = resp.PaySign
		case "NATIVE":
			result["code_url"] = resp.CodeURL
		case "APP":
			result["prepay_id"] = resp.PrepayID
		case "H5":
			result["payment_url"] = resp.PaymentURL
		}
	case payment.PaymentTypeAlipay:
		switch tradeType {
		case "PAGE":
			result["payment_form"] = resp.PaymentForm
		case "WAP", "APP":
			result["payment_data"] = resp.PaymentData
		}
	}

	return result, nil
}

// HandlePaymentNotify 处理支付通知
func (s *PaymentService) HandlePaymentNotify(ctx context.Context, paymentType string, data []byte) error {
	// 获取支付实例
	p, err := payment.Get(paymentType)
	if err != nil {
		return fmt.Errorf("不支持的支付方式: %w", err)
	}

	// 解析支付通知
	notify, err := p.ParseNotify(ctx, data)
	if err != nil {
		return fmt.Errorf("解析支付通知失败: %w", err)
	}

	// 获取订单信息
	order, err := s.orderRepo.GetByID(ctx, notify.OutTradeNo)
	if err != nil {
		return fmt.Errorf("获取订单失败: %w", err)
	}

	// 验证订单金额
	if order.TotalAmount != notify.TotalAmount {
		return errors.New("订单金额不匹配")
	}

	// 更新订单状态
	if notify.TradeStatus == "SUCCESS" || notify.TradeStatus == "TRADE_SUCCESS" {
		order.Status = model.OrderStatusPaid
		order.PaymentTime = &notify.PayTime
		order.PaymentTradeNo = notify.TradeNo

		if err := s.orderRepo.Update(ctx, order); err != nil {
			return fmt.Errorf("更新订单状态失败: %w", err)
		}

		// 触发订单支付成功事件
		// eventBus.Publish("order.paid", order)
	}

	return nil
}

// QueryPayment 查询支付状态
func (s *PaymentService) QueryPayment(ctx context.Context, orderID string) (map[string]interface{}, error) {
	// 获取订单信息
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	if order.PaymentType == "" {
		return nil, errors.New("订单未创建支付")
	}

	// 获取支付实例
	p, err := payment.Get(order.PaymentType)
	if err != nil {
		return nil, fmt.Errorf("不支持的支付方式: %w", err)
	}

	// 查询支付状态
	resp, err := p.Query(ctx, orderID, order.PaymentTradeNo)
	if err != nil {
		return nil, fmt.Errorf("查询支付状态失败: %w", err)
	}

	// 如果订单已支付但本地状态未更新，则更新订单状态
	if (resp.TradeStatus == "SUCCESS" || resp.TradeStatus == "TRADE_SUCCESS") && order.Status == model.OrderStatusPending {
		order.Status = model.OrderStatusPaid
		order.PaymentTime = &resp.PayTime

		if err := s.orderRepo.Update(ctx, order); err != nil {
			return nil, fmt.Errorf("更新订单状态失败: %w", err)
		}

		// 触发订单支付成功事件
		// eventBus.Publish("order.paid", order)
	}

	// 返回查询结果
	return map[string]interface{}{
		"order_id":     orderID,
		"payment_type": order.PaymentType,
		"trade_no":     resp.TradeNo,
		"trade_status": resp.TradeStatus,
		"total_amount": resp.TotalAmount,
		"pay_time":     resp.PayTime,
	}, nil
}

// RefundPayment 申请退款
func (s *PaymentService) RefundPayment(ctx context.Context, orderID, refundReason string, refundAmount float64) (map[string]interface{}, error) {
	// 获取订单信息
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	if order.Status != model.OrderStatusPaid && order.Status != model.OrderStatusShipped {
		return nil, errors.New("订单状态不允许退款")
	}

	if order.PaymentType == "" || order.PaymentTradeNo == "" {
		return nil, errors.New("订单未支付，无法退款")
	}

	// 获取支付实例
	p, err := payment.Get(order.PaymentType)
	if err != nil {
		return nil, fmt.Errorf("不支持的支付方式: %w", err)
	}

	// 生成退款单号
	outRefundNo := fmt.Sprintf("refund_%s_%d", orderID, time.Now().Unix())

	// 构建退款请求
	req := &payment.RefundRequest{
		OutTradeNo:   orderID,
		TradeNo:      order.PaymentTradeNo,
		OutRefundNo:  outRefundNo,
		TotalAmount:  order.TotalAmount,
		RefundAmount: refundAmount,
		RefundReason: refundReason,
	}

	// 申请退款
	resp, err := p.Refund(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("申请退款失败: %w", err)
	}

	// 更新订单状态
	order.Status = model.OrderStatusRefunding
	order.RefundAmount = refundAmount
	order.RefundReason = refundReason
	order.RefundNo = outRefundNo

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

	// 触发订单退款事件
	// eventBus.Publish("order.refunding", order)

	// 返回退款结果
	return map[string]interface{}{
		"order_id":      orderID,
		"refund_id":     resp.RefundID,
		"out_refund_no": resp.OutRefundNo,
		"refund_amount": resp.RefundAmount,
	}, nil
}

// QueryRefund 查询退款状态
func (s *PaymentService) QueryRefund(ctx context.Context, orderID string) (map[string]interface{}, error) {
	// 获取订单信息
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	if order.RefundNo == "" {
		return nil, errors.New("订单未申请退款")
	}

	// 获取支付实例
	p, err := payment.Get(order.PaymentType)
	if err != nil {
		return nil, fmt.Errorf("不支持的支付方式: %w", err)
	}

	// 查询退款状态
	resp, err := p.QueryRefund(ctx, order.RefundNo, "")
	if err != nil {
		return nil, fmt.Errorf("查询退款状态失败: %w", err)
	}

	// 返回查询结果
	return map[string]interface{}{
		"order_id":      orderID,
		"refund_id":     resp.RefundID,
		"out_refund_no": resp.OutRefundNo,
		"refund_amount": resp.RefundAmount,
		"raw_data":      resp.RawData,
	}, nil
}

// ClosePayment 关闭支付
func (s *PaymentService) ClosePayment(ctx context.Context, orderID string) error {
	// 获取订单信息
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("获取订单失败: %w", err)
	}

	if order.Status != model.OrderStatusPending {
		return errors.New("只有未支付的订单才能关闭支付")
	}

	if order.PaymentType == "" {
		return errors.New("订单未创建支付")
	}

	// 获取支付实例
	p, err := payment.Get(order.PaymentType)
	if err != nil {
		return fmt.Errorf("不支持的支付方式: %w", err)
	}

	// 关闭支付订单
	if err := p.CloseOrder(ctx, orderID); err != nil {
		return fmt.Errorf("关闭支付订单失败: %w", err)
	}

	// 更新订单状态
	order.Status = model.OrderStatusCancelled

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("更新订单状态失败: %w", err)
	}

	// 触发订单取消事件
	// eventBus.Publish("order.cancelled", order)

	return nil
}
