package service

import (
	"context"
	pb "github.com/Cospk/go-mall/api/rpc/gen/go/order"
	"github.com/Cospk/go-mall/internal/api/application/dto"
	"github.com/Cospk/go-mall/internal/api/infrastructure/rpc"
)

type OrderService struct {
	orderClient rpc.OrderServiceClient
}

func NewOrderService(orderClient interface{}) *OrderService {
	return &OrderService{
		orderClient: orderClient.(rpc.OrderServiceClient),
	}
}

func (s *OrderService) OrderCreate(c context.Context, req dto.OrderCreate, userId int64) (*dto.OrderCreateReply, error) {
	order, err := s.orderClient.CreateOrder(c, &pb.OrderCreateRequest{
		UserId:      userId,
		CartItemIds: req.CartItemIdList,
		AddressId:   req.UserAddressId,
	})

	if err != nil {
		return nil, err
	}
	reply := &dto.OrderCreateReply{
		OrderNo: order.OrderNo,
	}
	return reply, nil
}

// UserOrders 用户订单列表
func (s *OrderService) UserOrders(c context.Context, userId int64) ([]*dto.Order, error) {
	orders, err := s.orderClient.GetUserOrders(c, &pb.UserOrdersRequest{
		UserId:   userId,
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		return nil, err
	}
	// todo 这里需要将*dto.Order切片转为pb.UserOrdersRequest切片，这里临时处理一下
	orderList := make([]*dto.Order, 10)

	if orders != nil {

	}
	return orderList, nil
}

// OrderInfo 订单详情
func (s *OrderService) OrderInfo(c context.Context, orderNo string, userId int64) (*dto.Order, error) {
	info, err := s.orderClient.GetOrderInfo(c, &pb.OrderInfoRequest{
		OrderNo: orderNo,
		UserId:  userId,
	})
	if err != nil {
		return nil, err
	}

	orderInfo := &dto.Order{
		OrderNo: info.OrderNo,
		//....
	}
	return orderInfo, err
}

// OrderCancel 用户主动取消订单
func (s *OrderService) OrderCancel(c context.Context, orderNo string, userId int64) error {
	_, err := s.orderClient.CancelOrder(c, &pb.OrderCancelRequest{
		OrderNo: orderNo,
		UserId:  userId,
	})
	if err != nil {
		return err
	}
	return nil
}

// CreateOrderPay 订单发起支付
func (s *OrderService) CreateOrderPay(c context.Context, create *dto.OrderPayCreate, userId int64) error {

	_, err := s.orderClient.OrderCreatePay(c, &pb.OrderPayCreateRequest{
		OrderNo: create.OrderNo,
		UserId:  userId,
		PayType: create.PayType,
	})

	if err != nil {
		return err
	}
	return nil
}
