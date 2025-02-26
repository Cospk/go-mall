package service

import (
	"context"
	"github.com/Cospk/go-mall/internal/logic/domain"
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

func (s *DemoSvc) GetDemoList() ([]*domain.DemoDo, error) {
	demos, err := s.demoDomain.GetDemos()
	if err != nil {
		return nil, err
	}
	return demos, nil
}
