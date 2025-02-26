package domain

import (
	"context"
	"github.com/Cospk/go-mall/internal/dal/dao"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/utils"
	"time"
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

func (d DemoDomain) GetDemos() ([]*DemoDo, error) {
	demos, err := d.DemoDao.GetAllDemo()
	if err != nil {
		return nil, errcode.Wrap("查询 Demo 出错了", err)
	}

	demosDoS := make([]*DemoDo, 0, len(demos))
	for _, demo := range demos {
		demoDo := new(DemoDo)
		_ = utils.CopyStruct(demoDo, demo)
		demosDoS = append(demosDoS, demoDo)
	}
	return demosDoS, nil
}

type DemoDo struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	State     int8      `json:"state"`
	IsDel     uint      `json:"is_del"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
