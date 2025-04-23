package dao

import (
	"context"
	"github.com/Cospk/go-mall/internal/demo/dal/model"
	"github.com/Cospk/go-mall/internal/demo/logic/do"
	"github.com/Cospk/go-mall/pkg/utils"
)

// 演示Demo，后期使用删除
// Dao层 --> 数据操作层，只负责数据的增删改查

type DemoDao struct {
	ctx context.Context
}

func NewDemoDao(ctx context.Context) *DemoDao {
	return &DemoDao{
		ctx: ctx,
	}
}

// GetAllDemo 获取所有demo
func (d DemoDao) GetAllDemo() (demos []*model.DemoOrder, err error) {
	err = DB().WithContext(d.ctx).Find(&demos).Error
	if err != nil {
		return nil, err
	}
	return demos, nil
}

func (d DemoDao) CreateDemoOrder(order *do.DemoOrder) (*model.DemoOrder, error) {
	var demoOrder model.DemoOrder
	err := utils.CopyStruct(&demoOrder, order)
	if err != nil {
		return nil, err
	}
	err = DB().WithContext(d.ctx).Create(&demoOrder).Error
	return &demoOrder, err
}
