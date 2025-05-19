package model

import (
	"github.com/Cospk/go-mall/internal/user/domain/entity"
	"time"
)

// User 用户数据库模型
type User struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	LoginName string    `gorm:"column:login_name;type:varchar(64);uniqueIndex"`
	Password  string    `gorm:"column:password;type:varchar(128)"`
	NickName  string    `gorm:"column:nick_name;type:varchar(64)"`
	Slogan    string    `gorm:"column:slogan;type:varchar(255)"`
	Avatar    string    `gorm:"column:avatar;type:varchar(255)"`
	Verified  int64     `gorm:"column:verified;type:tinyint(1);default:0"`
	IsBlocked int32     `gorm:"column:is_blocked;type:tinyint(1);default:0"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP"`
}

func (u *User) TableName() string {
	return "gm_users"
}

// ToEntity 转换为实体
func (u *User) ToEntity() *entity.User {
	return &entity.User{
		ID:        u.ID,
		LoginName: u.LoginName,
		Password:  u.Password,
		NickName:  u.NickName,
		Slogan:    u.Slogan,
		Avatar:    u.Avatar,
		Verified:  u.Verified,
		IsBlocked: u.IsBlocked,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromEntity 从实体转换
func (u *User) FromEntity(user *entity.User) {
	u.ID = user.ID
	u.LoginName = user.LoginName
	u.Password = user.Password
	u.NickName = user.NickName
	u.Slogan = user.Slogan
	u.Avatar = user.Avatar
	u.Verified = user.Verified
	u.IsBlocked = user.IsBlocked
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
}
