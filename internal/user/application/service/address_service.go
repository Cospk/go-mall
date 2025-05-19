package service

import (
	"context"
	"github.com/Cospk/go-mall/internal/user/application/dto"
	"github.com/Cospk/go-mall/internal/user/domain/entity"
	"github.com/Cospk/go-mall/internal/user/domain/repository"
	"time"
)

// AddressService 地址服务接口
type AddressService interface {
	AddAddress(ctx context.Context, userID int64, req dto.AddAddressRequest) (dto.AddAddressResponse, error)
	GetAddressList(ctx context.Context, userID int64) ([]dto.AddressDTO, error)
	GetAddressInfo(ctx context.Context, addressID int64) (dto.AddressDTO, error)
	UpdateAddress(ctx context.Context, req dto.UpdateAddressRequest) (dto.UpdateAddressResponse, error)
	DeleteAddress(ctx context.Context, addressID int64) (dto.DeleteAddressResponse, error)
}

// addressService 地址服务实现
type addressService struct {
	addressRepo repository.AddressRepository
}

// NewAddressService 创建地址服务
func NewAddressService(addressRepo repository.AddressRepository) AddressService {
	return &addressService{
		addressRepo: addressRepo,
	}
}

// AddAddress 添加地址
func (s *addressService) AddAddress(ctx context.Context, userID int64, req dto.AddAddressRequest) (dto.AddAddressResponse, error) {
	// 如果设置为默认地址，需要将其他地址设置为非默认
	if req.Default == 1 {
		err := s.addressRepo.ClearDefaultAddress(ctx, userID)
		if err != nil {
			return dto.AddAddressResponse{
				Message: "设置默认地址失败",
			}, err
		}
	}

	// 创建地址实体
	address := entity.Address{
		UserID:       userID,
		UserName:     req.UserName,
		UserPhone:    req.UserPhone,
		Default:      req.Default,
		ProvinceName: req.ProvinceName,
		CityName:     req.CityName,
		RegionName:   req.RegionName,
		DetailAddr:   req.DetailAddr,
		CreatedAt:    time.Now(),
	}

	// 保存地址
	id, err := s.addressRepo.Save(ctx, address)
	if err != nil {
		return dto.AddAddressResponse{
			Message: "添加地址失败",
		}, err
	}

	return dto.AddAddressResponse{
		AddressID: id,
		Message:   "添加地址成功",
	}, nil
}

// GetAddressList 获取地址列表
func (s *addressService) GetAddressList(ctx context.Context, userID int64) ([]dto.AddressDTO, error) {
	addresses, err := s.addressRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []dto.AddressDTO
	for _, addr := range addresses {
		result = append(result, dto.AddressDTO{
			ID:              addr.ID,
			UserName:        addr.UserName,
			UserPhone:       addr.UserPhone,
			MaskedUserName:  maskName(addr.UserName),
			MaskedUserPhone: maskPhone(addr.UserPhone),
			Default:         addr.Default,
			ProvinceName:    addr.ProvinceName,
			CityName:        addr.CityName,
			RegionName:      addr.RegionName,
			DetailAddress:   addr.DetailAddr,
			CreatedAt:       addr.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

// GetAddressInfo 获取地址信息
func (s *addressService) GetAddressInfo(ctx context.Context, addressID int64) (dto.AddressDTO, error) {
	addr, err := s.addressRepo.FindByID(ctx, addressID)
	if err != nil {
		return dto.AddressDTO{}, err
	}

	return dto.AddressDTO{
		ID:              addr.ID,
		UserName:        addr.UserName,
		UserPhone:       addr.UserPhone,
		MaskedUserName:  maskName(addr.UserName),
		MaskedUserPhone: maskPhone(addr.UserPhone),
		Default:         addr.Default,
		ProvinceName:    addr.ProvinceName,
		CityName:        addr.CityName,
		RegionName:      addr.RegionName,
		DetailAddress:   addr.DetailAddr,
		CreatedAt:       addr.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// UpdateAddress 更新地址
func (s *addressService) UpdateAddress(ctx context.Context, req dto.UpdateAddressRequest) (dto.UpdateAddressResponse, error) {
	// 检查地址是否存在
	addr, err := s.addressRepo.FindByID(ctx, req.AddressID)
	if err != nil {
		return dto.UpdateAddressResponse{
			Success: false,
			Message: "地址不存在",
		}, err
	}

	// 如果设置为默认地址，需要将其他地址设置为非默认
	if req.Default == 1 && addr.Default != 1 {
		err := s.addressRepo.ClearDefaultAddress(ctx, addr.UserID)
		if err != nil {
			return dto.UpdateAddressResponse{
				Success: false,
				Message: "设置默认地址失败",
			}, err
		}
	}

	// 更新地址信息
	addr.UserName = req.UserName
	addr.UserPhone = req.UserPhone
	addr.Default = req.Default
	addr.ProvinceName = req.ProvinceName
	addr.CityName = req.CityName
	addr.RegionName = req.RegionName
	addr.DetailAddr = req.DetailAddress

	err = s.addressRepo.Update(ctx, addr)
	if err != nil {
		return dto.UpdateAddressResponse{
			Success: false,
			Message: "更新地址失败",
		}, err
	}

	return dto.UpdateAddressResponse{
		Success: true,
		Message: "更新地址成功",
	}, nil
}

// DeleteAddress 删除地址
func (s *addressService) DeleteAddress(ctx context.Context, addressID int64) (dto.DeleteAddressResponse, error) {
	err := s.addressRepo.Delete(ctx, addressID)
	if err != nil {
		return dto.DeleteAddressResponse{
			Success: false,
			Message: "删除地址失败",
		}, err
	}

	return dto.DeleteAddressResponse{
		Success: true,
		Message: "删除地址成功",
	}, nil
}

// 辅助函数：掩码处理姓名
func maskName(name string) string {
	if len(name) <= 1 {
		return name
	}
	return name[:1] + "**"
}

// 辅助函数：掩码处理电话
func maskPhone(phone string) string {
	if len(phone) <= 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}
