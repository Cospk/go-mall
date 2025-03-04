package service

import (
	"context"
	"github.com/Cospk/go-mall/api/request"
	"github.com/Cospk/go-mall/api/response"
	"github.com/Cospk/go-mall/internal/dal/cache"
	"github.com/Cospk/go-mall/internal/logic/do"
	"github.com/Cospk/go-mall/internal/logic/domain"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/Cospk/go-mall/pkg/utils"
	"go.uber.org/zap"
)

// 演示Demo，后期使用删除
// service --> 服务层，只负责服务的调用;业务用例的入口，负责流程编排和技术整合。

type DemoSvc struct {
	ctx        context.Context
	demoDomain *domain.DemoDomain
}

func NewDemoSvc(ctx context.Context) *DemoSvc {
	return &DemoSvc{
		ctx:        ctx,
		demoDomain: domain.NewDemoDomain(ctx),
	}
}

func (s *DemoSvc) GetDemoList() ([]*do.DemoOrder, error) {
	demos, err := s.demoDomain.GetDemos()
	if err != nil {
		return nil, err
	}
	return demos, nil
}

// CreateDemoOrder 创建demo订单
func (s *DemoSvc) CreateDemoOrder(orderRequest *request.DemoOrderCreate) (*response.DemoOrder, error) {
	// 创建模型可以使用new，可以
	domainOrder := new(do.DemoOrder)
	err := utils.CopyStruct(&domainOrder, orderRequest)
	if err != nil {

		return nil, errcode.Wrap("请求数据转换领域对象失败", err)
	}
	domainOrder, err = s.demoDomain.CreateDemoOrder(domainOrder)
	if err != nil {
		return nil, err
	}

	// TODO 做一些创建订单后其他的操作，比如发生通知、设置订单状态等等

	// 将订单信息写入到缓存，并读取出来，没有意义但可以演示一下
	cache.SetDemoOrder(s.ctx, domainOrder)
	cacheOrder, _ := cache.GetDemoOrder(s.ctx, domainOrder.OrderNo)
	logger.NewLogger(s.ctx).Info("缓存订单信息", zap.Any("cacheOrder", cacheOrder))

	// 将领域对象转换为响应对象
	responseOrder := new(response.DemoOrder)
	err = utils.CopyStruct(responseOrder, domainOrder)
	if err != nil {
		return nil, errcode.Wrap("demoOrderDo转换成响应体失败", err)
	}

	return responseOrder, nil
}
