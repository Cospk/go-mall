package service

import (
	"context"
	pb "github.com/Cospk/go-mall/api/rpc/gen/go/cart"
	"github.com/Cospk/go-mall/internal/api/application/dto"
	"github.com/Cospk/go-mall/internal/api/infrastructure/rpc"
)

type CartService struct {
	// 依赖注入
	cartClient rpc.CartServiceClient
}

func NewCartService(cartClient interface{}) *CartService {
	return &CartService{
		cartClient: cartClient.(rpc.CartServiceClient),
	}
}

// AddCartItem 添加购物车
func (s *CartService) AddCartItem(ctx context.Context, req *dto.AddCartItem, userId int64) (*pb.CommonResponse, error) {
	_, err := s.cartClient.AddCartItem(ctx, &pb.AddCartItemRequest{
		UserId:       userId,
		CommodityId:  req.CommodityId,
		CommodityNum: req.CommodityNum,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// UpdateCartItem 更改购物项的商品数
func (s *CartService) UpdateCartItem(ctx context.Context, cartItem *dto.CartItemUpdate, userId int64) error {
	item, err := s.cartClient.UpdateCartItem(ctx, &pb.UpdateCartItemRequest{
		UserId:       userId,
		ItemId:       cartItem.ItemId,
		CommodityNum: cartItem.CommodityNum,
	})

	if err != nil {
		return err
	}
	if item.Success {
		return nil
	}
	return nil
}

// UserCartItems 获取用户购物车中的购物项
func (s *CartService) UserCartItems(ctx context.Context, userId int64) (*pb.CartItemsReply, error) {
	cartItems, err := s.cartClient.UserCartItems(ctx, &pb.UserIdRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}
	return cartItems, nil
}

// DeleteUserCartItem 删除购物项
func (s *CartService) DeleteUserCartItem(ctx context.Context, userId int64, itemId int64) error {
	_, err := s.cartClient.DeleteUserCartItem(ctx, &pb.DeleteCartItemRequest{
		UserId: userId,
		ItemId: itemId,
	})

	if err != nil {
		return err
	}

	return nil
}

// CheckCartItemBill 查看购物项账单 -- 确认下单前用来显示商品和支付金额明细
func (s *CartService) CheckCartItemBill(ctx context.Context, itemIds []int64, userId int64) (*pb.CartItemBillReply, error) {

	bill, err := s.cartClient.CheckCartItemBill(ctx, &pb.CheckCartItemBillRequest{
		UserId:  userId,
		ItemIds: itemIds,
	})
	if err != nil {
		return nil, err
	}
	return bill, nil
}
