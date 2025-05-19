package entity

import "time"

// User 用户领域模型
type User struct {
	ID        int64     `json:"id"`
	LoginName string    `json:"login_name"`
	Password  string    `json:"-"` // 不序列化密码
	NickName  string    `json:"nick_name"`
	Slogan    string    `json:"slogan"`
	Avatar    string    `json:"avatar"`
	Verified  int64     `json:"verified"`
	IsBlocked int32     `json:"is_blocked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser 创建新用户
func NewUser(loginName, password, nickName, slogan, avatar string) *User {
	now := time.Now()
	return &User{
		LoginName: loginName,
		Password:  password, // 注意：这里应该是加密后的密码
		NickName:  nickName,
		Slogan:    slogan,
		Avatar:    avatar,
		Verified:  0,
		IsBlocked: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateInfo 更新用户信息
func (u *User) UpdateInfo(nickName, slogan, avatar string) {
	if nickName != "" {
		u.NickName = nickName
	}
	if slogan != "" {
		u.Slogan = slogan
	}
	if avatar != "" {
		u.Avatar = avatar
	}
	u.UpdatedAt = time.Now()
}

// UpdatePassword 更新密码
func (u *User) UpdatePassword(password string) {
	u.Password = password
	u.UpdatedAt = time.Now()
}

// Block 封禁用户
func (u *User) Block() {
	u.IsBlocked = 1
	u.UpdatedAt = time.Now()
}

// Unblock 解封用户
func (u *User) Unblock() {
	u.IsBlocked = 0
	u.UpdatedAt = time.Now()
}

// Verify 验证用户
func (u *User) Verify() {
	u.Verified = 1
	u.UpdatedAt = time.Now()
}
