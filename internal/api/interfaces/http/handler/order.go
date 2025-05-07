package handler

import (
	"github.com/Cospk/go-mall/internal/api/application/dto"
	"github.com/Cospk/go-mall/internal/api/application/service"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/resp"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	// 注入订单服务
	orderService *service.OrderService
}

func NewOrderHandler(orderService interface{}) *OrderHandler {
	return &OrderHandler{
		orderService: service.NewOrderService(orderService),
	}
}

// OrderCreate 创建订单
func (h *OrderHandler) OrderCreate(c *gin.Context) {
	request := new(dto.OrderCreate)
	if err := c.ShouldBindJSON(request); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	reply, err := h.orderService.OrderCreate(c, *request, c.GetInt64("userId"))
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	resp.NewResponse(c).Success(reply)
}

// UserOrders 用户订单列表
func (h *OrderHandler) UserOrders(c *gin.Context) {
	orders, err := h.orderService.UserOrders(c, c.GetInt64("userId"))
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
	}

	// todo
	// 返回分页订单列表,这里直接返回了
	resp.NewResponse(c).Success(orders)
}

// OrderInfo 订单详情
func (h *OrderHandler) OrderInfo(c *gin.Context) {
	info, err := h.orderService.OrderInfo(c, c.Param("order_no"), c.GetInt64("userId"))
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	resp.NewResponse(c).Success(info)
}

// OrderCancel 用户主动取消订单
func (h *OrderHandler) OrderCancel(c *gin.Context) {
	err := h.orderService.OrderCancel(c, c.Param("order_no"), c.GetInt64("userId"))
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	resp.NewResponse(c).SuccessOk()
}

// CreateOrderPay 订单发起支付
func (h *OrderHandler) CreateOrderPay(c *gin.Context) {
	request := new(dto.OrderPayCreate)
	if err := c.ShouldBindJSON(request); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	err := h.orderService.CreateOrderPay(c, request, c.GetInt64("userId"))
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))

		return
	}
	resp.NewResponse(c).SuccessOk()
}
