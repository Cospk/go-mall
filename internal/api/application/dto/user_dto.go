package dto

// ===================http 请求结构
type UserRegister struct {
	LoginName       string `json:"login_name" binding:"required,e164|email"` // 验证登录名必须为手机号或者邮箱地址
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
	Nickname        string `json:"nickname" binding:"max=30"`
	Slogan          string `json:"slogan" binding:"max=30"`
	Avatar          string `json:"avatar" binding:"max=100"`
}

// UserLogin 用户登录请求,需要同时验证和绑定Body和Header中的数据
// 使用Gin绑定RequestBoy和Header https://github.com/gin-gonic/gin/issues/2309#issuecomment-2020168668
type UserLogin struct {
	LoginName string `json:"login_name" binding:"required,e164|email"`
	Password  string `json:"password" binding:"required,min=8"`

	Platform string `json:"platform" header:"platform" binding:"required,oneof=H5 APP"`
}

type UserInfoUpdate struct {
	Nickname string `json:"nickname" binding:"max=30"`
	Slogan   string `json:"slogan" binding:"max=30"`
	Avatar   string `json:"avatar" binding:"max=100"`
}

type PasswordResetApply struct {
	LoginName string `json:"login_name" binding:"required,e164|email"` // 验证登录名必须为手机号或者邮箱地址
}

type PasswordReset struct {
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
	Token           string `json:"password_reset_token" binding:"required"`
	Code            string `json:"password_reset_code" binding:"required"`
}

type UserAddress struct {
	UserName      string `json:"user_name" binding:"required"`
	UserPhone     string `json:"user_phone" binding:"required"`
	Default       int32  `json:"default" binding:"oneof=0 1"`
	ProvinceName  string `json:"province_name" binding:"required"`
	CityName      string `json:"city_name" binding:"required"`
	RegionName    string `json:"region_name" binding:"required"`
	DetailAddress string `json:"detail_address" binding:"required"`
}

// ============ http 响应结构
type TokenReply struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	Duration      int64  `json:"duration"`
	SrvCreateTime string `json:"srv_create_time"`
}

type UserInfoReply struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	LoginName string `json:"login_name"`
	Verified  int    `json:"verified"`
	Avatar    string `json:"avatar"`
	Slogan    string `json:"slogan"`
	IsBlocked int32  `json:"is_blocked"`
	CreatedAt string `json:"created_at"`
}

// PasswordResetApplyReply 申请重置密码的响应
type PasswordResetApplyReply struct {
	PasswordResetToken string `json:"password_reset_token"`
}

type UserAddressReply struct {
	ID              int64  `json:"id"`
	UserName        string `json:"user_name"`
	UserPhone       string `json:"user_phone"`
	MaskedUserName  string `json:"masked_user_name"`  // 用于前台展示的脱敏后的用户姓名
	MaskedUserPhone string `json:"masked_user_phone"` // 用于前台展示的脱敏后的用户手机号
	Default         int32  `json:"default"`
	ProvinceName    string `json:"province_name"`
	CityName        string `json:"city_name"`
	RegionName      string `json:"region_name"`
	DetailAddress   string `json:"detail_address"`
	CreatedAt       string `json:"created_at"`
}

// rpc 请求/响应结构
