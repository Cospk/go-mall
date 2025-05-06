package service

import (
	"context"
	pb "github.com/Cospk/go-mall/api/rpc/gen/go/commodity"
	"github.com/Cospk/go-mall/internal/api/infrastructure/rpc"
	"github.com/Cospk/go-mall/pkg/errcode"
)

type CommodityService struct {
	commodityClient rpc.CommodityServiceClient
}

func NewCommodityService(commodityClient interface{}) *CommodityService {
	return &CommodityService{
		commodityClient: commodityClient.(rpc.CommodityServiceClient),
	}
}

// GetCategoryHierarchy 获取按层级划分后的所有分类
func (s *CommodityService) GetCategoryHierarchy(ctx context.Context) (*pb.CategoryHierarchyReply, error) {
	hierarchy, err := s.commodityClient.GetCategoryHierarchy(ctx, &pb.EmptyRequest{})

	if err != nil {
		return nil, errcode.ErrServer.WithCause(err)
	}
	return hierarchy, nil
}

// GetCategoriesWithParentId 按parentId查询分类列表
func (s *CommodityService) GetCategoriesWithParentId(ctx context.Context, parentId int64) (*pb.CategoriesReply, error) {
	subCategories, err := s.commodityClient.GetCategoriesWithParentId(ctx, &pb.ParentIdRequest{
		ParentId: parentId,
	})
	if err != nil {
		return nil, errcode.ErrServer.WithCause(err)
	}
	return subCategories, nil
}

// CommoditiesInCategory 分类商品列表
func (s *CommodityService) CommoditiesInCategory(ctx context.Context, categoryId int64) (*pb.CategoriesReply, error) {
	list, err := s.commodityClient.GetCategoriesWithParentId(ctx, &pb.ParentIdRequest{
		ParentId: categoryId,
	})
	if err != nil {
		return nil, errcode.ErrServer.WithCause(err)
	}
	return list, nil

}

// CommoditySearch 搜索商品
func (s *CommodityService) CommoditySearch(ctx context.Context, keyword string) (*pb.CommoditiesReply, error) {
	commodities, err := s.commodityClient.CommoditySearch(ctx, &pb.SearchRequest{
		Keyword:  keyword,
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		return nil, errcode.ErrServer.WithCause(err)
	}
	return commodities, nil
}

func (s *CommodityService) CommodityInfo(ctx context.Context, commodityId int64) (*pb.CommodityDetailReply, error) {
	info, err := s.commodityClient.CommodityInfo(ctx, &pb.CommodityIdRequest{
		CommodityId: commodityId,
	})
	if err != nil {
		return nil, errcode.ErrServer.WithCause(err)
	}
	return info, nil
}
