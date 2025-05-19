package dto

// RegisterRequest 注册请求
type RegisterRequest struct {
	LoginName       string
	Password        string
	PasswordConfirm string
	NickName        string
	Slogan          string
	Avatar          string
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	ID      int64
	Message string
}

// LoginRequest 登录请求
type LoginRequest struct {
	LoginName string
	Password  string
	Platform  string
}

// LoginResponse 登录响应
type LoginResponse struct {
	ID      int64
	Token   string
	Message string
}

// LogoutRequest 登出请求
type LogoutRequest struct {
	ID       int64
	Platform string
}

// LogoutResponse 登出响应
type LogoutResponse struct {
	Success bool
	Message string
}

// UserInfoDTO 用户信息DTO
type UserInfoDTO struct {
	ID        int64
	NickName  string
	LoginName string
	Verified  int64
	Avatar    string
	Slogan    string
	IsBlocked int32
	CreatedAt string
}

// RefreshTokenResponse 刷新token响应
type RefreshTokenResponse struct {
	AccessToken   string
	RefreshToken  string
	Duration      int64
	SrvCreateTime string
}

// PasswordResetApplyResponse 申请重置密码响应
type PasswordResetApplyResponse struct {
	Success            bool
	PasswordResetToken string
}

// PasswordResetRequest 重置密码请求
type PasswordResetRequest struct {
	Token           string
	Password        string
	ConfirmPassword string
	Code            string
}

// PasswordResetResponse 重置密码响应
type PasswordResetResponse struct {
	Success bool
	Message string
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	ID       int64
	NickName string
	Slogan   string
	Avatar   string
}

// UpdateUserInfoResponse 更新用户信息响应
type UpdateUserInfoResponse struct {
	Success bool
	Message string
}
