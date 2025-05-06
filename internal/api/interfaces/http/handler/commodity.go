package handler

import (
	"github.com/Cospk/go-mall/internal/api/application/dto"
	"github.com/Cospk/go-mall/internal/api/application/service"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/resp"
	"github.com/gin-gonic/gin"
	"strconv"
)

type CommodityHandler struct {
	// 注入CommodityService
	CommodityService *service.CommodityService
}

func NewCommodityHandler(commodityService interface{}) *CommodityHandler {
	return &CommodityHandler{
		CommodityService: service.NewCommodityService(commodityService),
	}
}

// GetCategoryHierarchy 获取按层级划分后的所有分类
func (h *CommodityHandler) GetCategoryHierarchy(c *gin.Context) {
	replyData, err := h.CommodityService.GetCategoryHierarchy(c)
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	resp.NewResponse(c).Success(replyData)
}

// GetCategoriesWithParentId 按parentId查询分类列表
func (h *CommodityHandler) GetCategoriesWithParentId(c *gin.Context) {
	parentId, _ := strconv.ParseInt(c.Query("parent_id"), 10, 64)
	categoryList, err := h.CommodityService.GetCategoriesWithParentId(c, parentId)
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	resp.NewResponse(c).Success(categoryList)
}

// CommoditiesInCategory 分类商品列表
func (h *CommodityHandler) CommoditiesInCategory(c *gin.Context) {
	categoryId, _ := strconv.ParseInt(c.Query("category_id"), 10, 64)
	categoryList, err2 := h.CommodityService.CommoditiesInCategory(c, categoryId)
	if err2 != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err2))
		return
	}
	resp.NewResponse(c).Success(categoryList)
}

// CommoditySearch 搜索商品
func (h *CommodityHandler) CommoditySearch(c *gin.Context) {
	searchQuery := new(dto.CommoditySearch)
	if err := c.ShouldBindQuery(searchQuery); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	search, err := h.CommodityService.CommoditySearch(c, searchQuery.Keyword)
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	resp.NewResponse(c).Success(search)
}

func (h *CommodityHandler) CommodityInfo(c *gin.Context) {
	commodityId, _ := strconv.ParseInt(c.Param("commodity_id"), 10, 64)
	if commodityId <= 0 {
		resp.NewResponse(c).Error(errcode.ErrParams)
		return
	}

	info, err := h.CommodityService.CommodityInfo(c, commodityId)

	if err == nil {
		resp.NewResponse(c).Error(errcode.ErrParams)
		return
	}

	resp.NewResponse(c).Success(info)
}
