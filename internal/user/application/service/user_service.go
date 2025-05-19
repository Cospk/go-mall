package service

import (
	"context"
	"errors"
	"github.com/Cospk/go-mall/internal/user/application/dto"
	"github.com/Cospk/go-mall/internal/user/domain/entity"
	"github.com/Cospk/go-mall/internal/user/domain/repository"
	"time"
)

// UserService 用户服务接口
type UserService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	Logout(ctx context.Context, req dto.LogoutRequest) (dto.LogoutResponse, error)
	GetUserInfo(ctx context.Context, id int64) (dto.UserInfoDTO, error)
	RefreshToken(ctx context.Context, refreshToken string) (dto.RefreshTokenResponse, error)
	PasswordResetApply(ctx context.Context, loginName string) (dto.PasswordResetApplyResponse, error)
	PasswordReset(ctx context.Context, req dto.PasswordResetRequest) (dto.PasswordResetResponse, error)
	UpdateUserInfo(ctx context.Context, req dto.UpdateUserInfoRequest) (dto.UpdateUserInfoResponse, error)
}

// UserService 用户应用服务
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Register 注册用户
func (s *userService) Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// 检查密码是否匹配
	if req.Password != req.PasswordConfirm {
		return dto.RegisterResponse{
			Message: "密码不匹配",
		}, errors.New("密码不匹配")
	}

	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByLoginName(ctx, req.LoginName)
	if err != nil {
		return dto.RegisterResponse{
			Message: "检查用户名失败",
		}, err
	}
	if exists {
		return dto.RegisterResponse{
			Message: "用户名已存在",
		}, errors.New("用户名已存在")
	}

	// 创建用户实体
	user := entity.User{
		LoginName: req.LoginName,
		Password:  req.Password, // 注意：实际应用中应该对密码进行哈希处理
		NickName:  req.NickName,
		Slogan:    req.Slogan,
		Avatar:    req.Avatar,
		CreatedAt: time.Now(),
	}

	// 保存用户
	id, err := s.userRepo.Save(ctx, user)
	if err != nil {
		return dto.RegisterResponse{
			Message: "注册失败",
		}, err
	}

	return dto.RegisterResponse{
		ID:      id,
		Message: "注册成功",
	}, nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	// 根据登录名查找用户
	user, err := s.userRepo.FindByLoginName(ctx, req.LoginName)
	if err != nil {
		return dto.LoginResponse{
			Message: "用户不存在",
		}, err
	}

	// 验证密码
	if user.Password != req.Password { // 注意：实际应用中应该对密码进行哈希比较
		return dto.LoginResponse{
			Message: "密码错误",
		}, errors.New("密码错误")
	}

	// 生成token（实际应用中应该使用JWT等方式）
	token := "mock_token_" + user.LoginName + "_" + time.Now().String()

	// 更新用户登录状态
	err = s.userRepo.UpdateLoginStatus(ctx, user.ID, req.Platform, token)
	if err != nil {
		return dto.LoginResponse{
			Message: "登录失败",
		}, err
	}

	return dto.LoginResponse{
		ID:      user.ID,
		Token:   token,
		Message: "登录成功",
	}, nil
}

// Logout 用户登出
func (s *userService) Logout(ctx context.Context, req dto.LogoutRequest) (dto.LogoutResponse, error) {
	// 清除用户登录状态
	err := s.userRepo.ClearLoginStatus(ctx, req.ID, req.Platform)
	if err != nil {
		return dto.LogoutResponse{
			Success: false,
			Message: "登出失败",
		}, err
	}

	return dto.LogoutResponse{
		Success: true,
		Message: "登出成功",
	}, nil
}

// GetUserInfo 获取用户信息
func (s *userService) GetUserInfo(ctx context.Context, id int64) (dto.UserInfoDTO, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return dto.UserInfoDTO{}, err
	}

	return dto.UserInfoDTO{
		ID:        user.ID,
		NickName:  user.NickName,
		LoginName: user.LoginName,
		Verified:  user.Verified,
		Avatar:    user.Avatar,
		Slogan:    user.Slogan,
		IsBlocked: user.IsBlocked,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// RefreshToken 刷新用户token
func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (dto.RefreshTokenResponse, error) {
	// 验证刷新token
	userID, err := s.userRepo.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	// 生成新的访问token和刷新token
	accessToken := "new_access_token_" + time.Now().String()
	newRefreshToken := "new_refresh_token_" + time.Now().String()
	duration := int64(3600) // token有效期，单位秒

	// 更新token
	err = s.userRepo.UpdateTokens(ctx, userID, accessToken, newRefreshToken)
	if err != nil {
		return dto.RefreshTokenResponse{}, err
	}

	return dto.RefreshTokenResponse{
		AccessToken:   accessToken,
		RefreshToken:  newRefreshToken,
		Duration:      duration,
		SrvCreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// PasswordResetApply 申请重置密码
func (s *userService) PasswordResetApply(ctx context.Context, loginName string) (dto.PasswordResetApplyResponse, error) {
	// 检查用户是否存在
	exists, err := s.userRepo.ExistsByLoginName(ctx, loginName)
	if err != nil {
		return dto.PasswordResetApplyResponse{
			Success: false,
		}, err
	}
	if !exists {
		return dto.PasswordResetApplyResponse{
			Success: false,
		}, errors.New("用户不存在")
	}

	// 生成密码重置token
	resetToken := "reset_token_" + loginName + "_" + time.Now().String()

	// 保存重置token
	err = s.userRepo.SavePasswordResetToken(ctx, loginName, resetToken)
	if err != nil {
		return dto.PasswordResetApplyResponse{
			Success: false,
		}, err
	}

	// 实际应用中应该发送邮件

	return dto.PasswordResetApplyResponse{
		Success:            true,
		PasswordResetToken: resetToken,
	}, nil
}

// PasswordReset 重置密码
func (s *userService) PasswordReset(ctx context.Context, req dto.PasswordResetRequest) (dto.PasswordResetResponse, error) {
	// 验证密码
	if req.Password != req.ConfirmPassword {
		return dto.PasswordResetResponse{
			Success: false,
			Message: "密码不匹配",
		}, errors.New("密码不匹配")
	}

	// 验证token和验证码
	valid, loginName, err := s.userRepo.ValidatePasswordResetToken(ctx, req.Token, req.Code)
	if err != nil || !valid {
		return dto.PasswordResetResponse{
			Success: false,
			Message: "无效的重置token或验证码",
		}, errors.New("无效的重置token或验证码")
	}

	// 更新密码
	err = s.userRepo.UpdatePassword(ctx, loginName, req.Password)
	if err != nil {
		return dto.PasswordResetResponse{
			Success: false,
			Message: "密码重置失败",
		}, err
	}

	return dto.PasswordResetResponse{
		Success: true,
		Message: "密码重置成功",
	}, nil
}

// UpdateUserInfo 更新用户信息
func (s *userService) UpdateUserInfo(ctx context.Context, req dto.UpdateUserInfoRequest) (dto.UpdateUserInfoResponse, error) {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(ctx, req.ID)
	if err != nil {
		return dto.UpdateUserInfoResponse{
			Success: false,
			Message: "用户不存在",
		}, err
	}

	// 更新用户信息
	user.NickName = req.NickName
	user.Slogan = req.Slogan
	user.Avatar = req.Avatar

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return dto.UpdateUserInfoResponse{
			Success: false,
			Message: "更新失败",
		}, err
	}

	return dto.UpdateUserInfoResponse{
		Success: true,
		Message: "更新成功",
	}, nil
}
