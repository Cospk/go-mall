package repository

import (
	"context"
	"github.com/Cospk/go-mall/internal/user/domain/entity"
)

// AddressRepository 地址仓储接口
type AddressRepository interface {
	// 创建地址
	Create(ctx context.Context, address *entity.Address) (int64, error)

	// 根据ID获取地址
	GetByID(ctx context.Context, id int64) (*entity.Address, error)

	// 获取用户的所有地址
	GetByUserID(ctx context.Context, userID int64) ([]*entity.Address, error)

	// 更新地址
	Update(ctx context.Context, address *entity.Address) error

	// 删除地址
	Delete(ctx context.Context, id int64) error

	// 设置默认地址
	SetDefault(ctx context.Context, id int64, userID int64) error

	// 获取用户的默认地址
	GetDefaultByUserID(ctx context.Context, userID int64) (*entity.Address, error)
}
