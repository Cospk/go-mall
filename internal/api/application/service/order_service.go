package service

import "github.com/Cospk/go-mall/internal/api/infrastructure/rpc"

type OrderService struct {
	orderClient rpc.OrderServiceClient
}

func NewOrderService(orderClient interface{}) *OrderService {
	return &OrderService{
		orderClient: orderClient.(rpc.OrderServiceClient),
	}
}
