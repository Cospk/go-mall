package service

import (
	"context"
	"github.com/Cospk/go-mall/internal/user/domain/entity"
	"github.com/Cospk/go-mall/internal/user/domain/repository"
)

// UserService 用户应用服务
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, userID int64) (*entity.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
