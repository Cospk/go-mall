package dao

import (
	"context"
	"github.com/Cospk/go-mall/internal/dal/model"
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
func (d DemoDao) GetAllDemo() (demos []*model.Demo, err error) {
	err = DB().WithContext(d.ctx).Find(&demos).Error
	if err != nil {
		return nil, err
	}
	return demos, nil
}
