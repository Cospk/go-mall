package service

import (
	"context"
	"github.com/Cospk/go-mall/internal/api/application/dto"
	//"github.com/Cospk/go-mall/internal/api/domain/service"
	pb "github.com/Cospk/go-mall/api/rpc/gen/go/user"
	"github.com/Cospk/go-mall/internal/api/infrastructure/rpc"
)

type UserService struct {
	userClient rpc.UserServiceClient
	//userDomainService *service.UserDomainService
}

func NewUserService(userClient interface{}) *UserService {
	return &UserService{
		userClient: userClient.(rpc.UserServiceClient),
		//userDomainService: service.NewUserDomainService(),
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req dto.UserRegister) (int64, error) {
	// 调用RPC服务
	resp, err := s.userClient.Register(ctx, &pb.RegisterRequest{
		LoginName:       req.LoginName,
		Password:        req.Password,
		PasswordConfirm: req.PasswordConfirm,
		NickName:        req.Nickname,
		Slogan:          req.Slogan,
		Avatar:          req.Avatar,
	})
	if err != nil {
		return 0, err
	}

	return resp.Id, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req dto.UserLogin) (string, error) {
	// 调用RPC服务
	resp, err := s.userClient.Login(ctx, &pb.LoginRequest{
		LoginName: req.LoginName,
		Password:  req.Password,
		Platform:  req.Platform,
	})
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

// Logout 用户登出
func (s *UserService) Logout(ctx context.Context, id int64, platform string) error {
	_, err := s.userClient.Logout(ctx, &pb.LogoutRequest{
		Id:       id,
		Platform: platform,
	})
	return err
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, userID int64) (*dto.UserInfoReply, error) {
	// 调用RPC服务
	resp, err := s.userClient.GetUserInfo(ctx, &pb.GetUserInfoRequest{
		Id: userID,
	})
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return &dto.UserInfoReply{
		ID:        resp.Id,
		Nickname:  resp.NickName,
		LoginName: resp.LoginName,
		Verified:  int(resp.Verified),
		Avatar:    resp.Avatar,
		Slogan:    resp.Slogan,
		IsBlocked: resp.IsBlocked,
		CreatedAt: resp.CreatedAt,
	}, nil
}

// RefreshUserToken 刷新用户token
func (s *UserService) RefreshUserToken(ctx context.Context, refreshToken string) (*dto.TokenReply, error) {
	resp, err := s.userClient.RefreshUserToken(ctx, &pb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}

	return &dto.TokenReply{
		AccessToken:   resp.AccessToken,
		RefreshToken:  resp.RefreshToken,
		Duration:      resp.Duration,
		SrvCreateTime: resp.SrvCreateTime,
	}, nil
}

// PasswordResetApply 申请重置登录密码
func (s *UserService) PasswordResetApply(ctx context.Context, req dto.PasswordResetApply) (*dto.PasswordResetApplyReply, error) {
	resp, err := s.userClient.PasswordResetApply(ctx, &pb.PasswordResetApplyRequest{
		LoginName: req.LoginName,
	})
	if err != nil {
		return nil, err
	}

	return &dto.PasswordResetApplyReply{
		PasswordResetToken: resp.PasswordResetToken,
	}, nil
}

// PasswordReset 重置登录密码
func (s *UserService) PasswordReset(ctx context.Context, req dto.PasswordReset) error {
	_, err := s.userClient.PasswordReset(ctx, &pb.PasswordResetRequest{
		Token:           req.Token,
		Password:        req.Password,
		ConfirmPassword: req.PasswordConfirm,
		Code:            req.Code,
	})
	if err != nil {
		return err
	}
	return nil
}

// GetUserInfoById 获取用户信息-ID
func (s *UserService) GetUserInfoById(ctx context.Context, id int64) (*dto.UserInfoReply, error) {
	resp, err := s.userClient.GetUserInfoById(ctx, &pb.GetUserInfoByIdRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &dto.UserInfoReply{
		ID:        resp.Id,
		Nickname:  resp.NickName,
		LoginName: resp.LoginName,
		Verified:  int(resp.Verified),
		Avatar:    resp.Avatar,
		Slogan:    resp.Slogan,
		IsBlocked: resp.IsBlocked,
		CreatedAt: resp.CreatedAt,
	}, nil
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(ctx context.Context, req dto.UserInfoUpdate) error {
	_, err := s.userClient.UpdateUserInfo(ctx, &pb.UpdateUserInfoRequest{
		NickName: req.Nickname,
		Slogan:   req.Slogan,
		Avatar:   req.Avatar,
	})
	return err
}

// AddUserAddressInfo 添加用户收货地址
func (s *UserService) AddUserAddressInfo(ctx context.Context, req dto.UserAddress) (int64, error) {
	resp, err := s.userClient.AddUserAddressInfo(ctx, &pb.AddUserAddressInfoRequest{
		UserName:     req.UserName,
		UserPhone:    req.UserPhone,
		Default:      req.Default,
		ProvinceName: req.ProvinceName,
		CityName:     req.CityName,
		RegionName:   req.RegionName,
		DetailAddr:   req.DetailAddress,
	})
	if err != nil {
		return 0, err
	}
	return resp.AddressId, nil
}

// GetUserAddressList 获取用户收货地址列表
func (s *UserService) GetUserAddressList(ctx context.Context, id int64) ([]*dto.UserAddressReply, error) {
	resp, err := s.userClient.GetUserAddressList(ctx, &pb.GetUserAddressListRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	addresses := make([]*dto.UserAddressReply, 0, len(resp.AddressList))

	for _, addr := range resp.AddressList {
		addresses = append(addresses, &dto.UserAddressReply{
			ID:              addr.Id,
			UserName:        addr.UserName,
			UserPhone:       addr.UserPhone,
			MaskedUserName:  addr.MaskedUserName,
			MaskedUserPhone: addr.MaskedUserPhone,
			Default:         addr.Default,
			ProvinceName:    addr.ProvinceName,
			CityName:        addr.CityName,
			RegionName:      addr.RegionName,
			DetailAddress:   addr.DetailAddress,
			CreatedAt:       addr.CreatedAt,
		})
	}
	return addresses, nil
}

// GetUserAddressInfo 获取单个收货地址信息
func (s *UserService) GetUserAddressInfo(ctx context.Context, address_id int64) (*dto.UserAddressReply, error) {
	address, err := s.userClient.GetUserAddressInfo(ctx, &pb.GetUserAddressInfoRequest{
		AddressId: address_id,
	})
	if err != nil {
		return nil, err
	}

	return &dto.UserAddressReply{
		ID:              address.Id,
		UserName:        address.UserName,
		UserPhone:       address.UserPhone,
		MaskedUserName:  address.MaskedUserName,
		MaskedUserPhone: address.MaskedUserPhone,
		Default:         address.Default,
		ProvinceName:    address.ProvinceName,
		CityName:        address.CityName,
		RegionName:      address.RegionName,
		DetailAddress:   address.DetailAddress,
		CreatedAt:       address.CreatedAt,
	}, nil
}

// UpdateUserAddressInfo 更新用户收货地址信息
func (s *UserService) UpdateUserAddressInfo(ctx context.Context, address_id int64, req dto.UserAddress) error {
	_, err := s.userClient.UpdateUserAddressInfo(ctx, &pb.UpdateUserAddressInfoRequest{
		AddressId:     address_id,
		UserName:      req.UserName,
		UserPhone:     req.UserPhone,
		Default:       req.Default,
		ProvinceName:  req.ProvinceName,
		CityName:      req.CityName,
		RegionName:    req.RegionName,
		DetailAddress: req.DetailAddress,
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserAddressInfo 删除用户收货地址
func (s *UserService) DeleteUserAddressInfo(ctx context.Context, address_id int64) error {
	_, err := s.userClient.DeleteUserAddressInfo(ctx, &pb.DeleteUserAddressInfoRequest{
		AddressId: address_id,
	})
	if err != nil {
		return err
	}
	return nil
}
