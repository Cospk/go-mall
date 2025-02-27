package domain

import (
	"context"
	"github.com/Cospk/go-mall/internal/dal/dao"
	"github.com/Cospk/go-mall/internal/logic/do"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/utils"
)

// 演示Demo，后期使用删除
// Domain层 --> 业务逻辑层，只负责业务逻辑;业务逻辑的核心，负责规则实现和状态管理

type DemoDomain struct {
	ctx     context.Context
	DemoDao *dao.DemoDao
}

func NewDemoDomain(ctx context.Context) *DemoDomain {
	return &DemoDomain{
		ctx:     ctx,
		DemoDao: dao.NewDemoDao(ctx),
	}
}

func (d DemoDomain) GetDemos() ([]*do.DemoOrder, error) {
	demos, err := d.DemoDao.GetAllDemo()
	if err != nil {
		return nil, errcode.Wrap("查询 Demo 出错了", err)
	}

	demosDoS := make([]*do.DemoOrder, 0, len(demos))
	for _, demo := range demos {
		demoDo := new(do.DemoOrder)
		_ = utils.CopyStruct(demoDo, demo)
		demosDoS = append(demosDoS, demoDo)
	}
	return demosDoS, nil
}

func (d DemoDomain) CreateDemoOrder(demoOrder *do.DemoOrder) (*do.DemoOrder, error) {
	// 随便写一个订单号
	demoOrder.OrderNo = "20240627596615375920904456"

	demoOrderModel, err := d.DemoDao.CreateDemoOrder(demoOrder)
	if err != nil {
		return nil, errcode.Wrap("创建 Demo 出错了", err)
	}

	// 写订单快照，这里不演示了

	err = utils.CopyStruct(demoOrder, demoOrderModel)
	return demoOrder, nil
}
