package mysql

import (
	"context"
	"errors"
	"github.com/Cospk/go-mall/internal/user/domain/entity"
	"github.com/Cospk/go-mall/internal/user/domain/repository"
	"github.com/Cospk/go-mall/internal/user/infrastructure/model"
	"gorm.io/gorm"
)

// AddressRepositoryImpl MySQL实现的地址仓储
type AddressRepositoryImpl struct {
	db *gorm.DB
}

// NewAddressRepositoryImpl 创建地址仓储
func NewAddressRepositoryImpl(db *gorm.DB) repository.AddressRepository {
	return &AddressRepositoryImpl{db: db}
}

// Create 创建地址
func (r *AddressRepositoryImpl) Create(ctx context.Context, address *entity.Address) (int64, error) {
	// 如果是默认地址，先将该用户的所有地址设为非默认
	if address.IsDefault == 1 {
		r.db.WithContext(ctx).Model(&model.Address{}).Where("user_id = ?", address.UserID).Update("is_default", false)
	}

	addressModel := &model.Address{}
	addressModel.FromEntity(address)

	result := r.db.WithContext(ctx).Create(addressModel)
	if result.Error != nil {
		return 0, result.Error
	}

	return addressModel.ID, nil
}

// GetByID 根据ID获取地址
func (r *AddressRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Address, error) {
	addressModel := &model.Address{}
	result := r.db.WithContext(ctx).First(addressModel, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return addressModel.ToEntity(), nil
}

// GetByUserID 获取用户的所有地址
func (r *AddressRepositoryImpl) GetByUserID(ctx context.Context, userID int64) ([]*entity.Address, error) {
	var addressModels []*model.Address
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC, id DESC").Find(&addressModels)
	if result.Error != nil {
		return nil, result.Error
	}

	addresses := make([]*entity.Address, 0, len(addressModels))
	for _, addressModel := range addressModels {
		addresses = append(addresses, addressModel.ToEntity())
	}

	return addresses, nil
}

// Update 更新地址
func (r *AddressRepositoryImpl) Update(ctx context.Context, address *entity.Address) error {
	// 如果是默认地址，先将该用户的所有地址设为非默认
	if address.IsDefault == 1 {
		r.db.WithContext(ctx).Model(&model.Address{}).Where("user_id = ?", address.UserID).Update("is_default", false)
	}

	addressModel := &model.Address{}
	addressModel.FromEntity(address)

	result := r.db.WithContext(ctx).Model(&model.Address{}).Where("id = ?", address.ID).Updates(map[string]interface{}{
		"user_name":     address.UserName,
		"user_phone":    address.UserPhone,
		"is_default":    address.IsDefault,
		"province_name": address.ProvinceName,
		"city_name":     address.CityName,
		"region_name":   address.RegionName,
		"detail_addr":   address.DetailAddr,
		"updated_at":    address.UpdatedAt,
	})

	return result.Error
}

// Delete 删除地址
func (r *AddressRepositoryImpl) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&model.Address{}, id)
	return result.Error
}

// SetDefault 设置默认地址
func (r *AddressRepositoryImpl) SetDefault(ctx context.Context, id int64, userID int64) error {
	// 先将该用户的所有地址设为非默认
	r.db.WithContext(ctx).Model(&model.Address{}).Where("user_id = ?", userID).Update("is_default", false)

	// 将指定地址设为默认
	result := r.db.WithContext(ctx).Model(&model.Address{}).Where("id = ?", id).Update("is_default", true)
	return result.Error
}

// GetDefaultByUserID 获取用户的默认地址
func (r *AddressRepositoryImpl) GetDefaultByUserID(ctx context.Context, userID int64) (*entity.Address, error) {
	addressModel := &model.Address{}
	result := r.db.WithContext(ctx).Where("user_id = ? AND is_default = ?", userID, true).First(addressModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return addressModel.ToEntity(), nil
}
