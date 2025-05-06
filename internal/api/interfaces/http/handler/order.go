package handler

import "github.com/Cospk/go-mall/internal/api/application/service"

type OrderHandler struct {
	// 注入订单服务
	orderService *service.OrderService
}

func NewOrderHandler(orderService interface{}) *OrderHandler {
	return &OrderHandler{
		orderService: service.NewOrderService(orderService),
	}
}
