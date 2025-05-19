package repository

import (
	"context"
	"github.com/Cospk/go-mall/internal/user/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *entity.User) (int64, error)

	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id int64) (*entity.User, error)

	// GetByLoginName 根据登录名获取用户
	GetByLoginName(ctx context.Context, loginName string) (*entity.User, error)

	// Update 更新用户
	Update(ctx context.Context, user *entity.User) error
}
