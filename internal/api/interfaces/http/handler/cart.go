package handler

import (
	"errors"
	"github.com/Cospk/go-mall/internal/api/application/dto"
	"github.com/Cospk/go-mall/internal/api/application/service"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/resp"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"strconv"
)

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler(cartService interface{}) *CartHandler {
	return &CartHandler{
		cartService: service.NewCartService(cartService),
	}
}

// AddToCart 添加商品到购物车
func (h *CartHandler) AddToCart(c *gin.Context) {
	// 从请求中获取用户ID和商品ID
}

// AddCartItem 添加购物车
func (h *CartHandler) AddCartItem(c *gin.Context) {
	request := new(dto.AddCartItem)
	if err := c.ShouldBindJSON(request); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	item, err := h.cartService.AddCartItem(c, request, c.GetInt64("userId"))
	if err != nil {
		// WithCause 记得加, 不然请求的错误日志里记不到错误原因
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))

		return
	}

	resp.NewResponse(c).Success(item)
}

// UpdateCartItem 更改购物项的商品数
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	request := new(dto.CartItemUpdate)
	if err := c.ShouldBindJSON(request); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	err := h.cartService.UpdateCartItem(c, request, c.GetInt64("userId"))

	if err != nil {
		if errors.Is(err, errcode.ErrParams) {
			resp.NewResponse(c).Error(errcode.ErrParams)
		} else {
			// WithCause 记得加, 不然请求的错误日志里记不到错误原因
			resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	resp.NewResponse(c).SuccessOk()
}

// UserCartItems 获取用户购物车中的购物项
func (h *CartHandler) UserCartItems(c *gin.Context) {
	items, err := h.cartService.UserCartItems(c, c.GetInt64("userId"))
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	resp.NewResponse(c).Success(items)
}

// DeleteUserCartItem 删除购物项
func (h *CartHandler) DeleteUserCartItem(c *gin.Context) {
	itemId, _ := strconv.ParseInt(c.Param("item_id"), 10, 64)
	err := h.cartService.DeleteUserCartItem(c, itemId, c.GetInt64("userId"))
	if err != nil {
		if errors.Is(err, errcode.ErrParams) {
			resp.NewResponse(c).Error(errcode.ErrParams)
		} else {
			resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	resp.NewResponse(c).SuccessOk()
}

// CheckCartItemBill 查看购物项账单 -- 确认下单前用来显示商品和支付金额明细
func (h *CartHandler) CheckCartItemBill(c *gin.Context) {
	itemIdList := c.QueryArray("item_id")
	if len(itemIdList) == 0 {
		resp.NewResponse(c).Error(errcode.ErrParams)
	}

	itemIds := lo.Map(itemIdList, func(itemId string, index int) int64 {
		i, _ := strconv.ParseInt(itemId, 10, 64)
		return i
	})

	bill, err := h.cartService.CheckCartItemBill(c, itemIds, c.GetInt64("userId"))
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	resp.NewResponse(c).Success(bill)
}
