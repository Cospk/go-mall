package mysql

import (
	"context"
	"errors"
	"github.com/Cospk/go-mall/internal/user/domain/entity"
	"github.com/Cospk/go-mall/internal/user/domain/repository"
	"github.com/Cospk/go-mall/internal/user/infrastructure/model"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// UserRepositoryImpl 用户仓储MySQL实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) (int64, error) {
	userModel := &model.User{}
	userModel.FromEntity(user)

	result := r.db.WithContext(ctx).Create(userModel)
	if result.Error != nil {
		return 0, result.Error
	}

	return userModel.ID, nil
}

// GetByID 根据ID获取用户
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	userModel := &model.User{}
	result := r.db.WithContext(ctx).First(userModel, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return userModel.ToEntity(), nil
}

// GetByLoginName 根据登录名获取用户
func (r *UserRepositoryImpl) GetByLoginName(ctx context.Context, loginName string) (*entity.User, error) {
	userModel := &model.User{}
	result := r.db.WithContext(ctx).Where("login_name = ?", loginName).First(userModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return userModel.ToEntity(), nil
}

// Update 更新用户信息
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	userModel := &model.User{}
	userModel.FromEntity(user)

	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"nick_name":  user.NickName,
		"slogan":     user.Slogan,
		"avatar":     user.Avatar,
		"password":   user.Password,
		"verified":   user.Verified,
		"is_blocked": user.IsBlocked,
		"updated_at": user.UpdatedAt,
	})

	return result.Error
}
