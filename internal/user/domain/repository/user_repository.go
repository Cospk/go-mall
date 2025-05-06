package repository

import (
	"context"
	"github.com/Cospk/go-mall/internal/user/domain/model"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id int64) (*model.User, error)

	// 其他方法...
}
