package model

import (
	"github.com/Cospk/go-mall/internal/user/domain/entity"
	"time"
)

// Address 地址数据库模型
type Address struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int64     `gorm:"column:user_id;index"`
	UserName     string    `gorm:"column:user_name;type:varchar(64)"`
	UserPhone    string    `gorm:"column:user_phone;type:varchar(32)"`
	IsDefault    int32     `gorm:"column:is_default;type:tinyint(1);default:0"`
	ProvinceName string    `gorm:"column:province_name;type:varchar(64)"`
	CityName     string    `gorm:"column:city_name;type:varchar(64)"`
	RegionName   string    `gorm:"column:region_name;type:varchar(64)"`
	DetailAddr   string    `gorm:"column:detail_addr;type:varchar(255)"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP"`
}

func (Address) TableName() string {
	return "addresses"
}

// ToEntity 转换为实体
func (a *Address) ToEntity() *entity.Address {
	return &entity.Address{
		ID:           a.ID,
		UserID:       a.UserID,
		UserName:     a.UserName,
		UserPhone:    a.UserPhone,
		IsDefault:    a.IsDefault,
		ProvinceName: a.ProvinceName,
		CityName:     a.CityName,
		RegionName:   a.RegionName,
		DetailAddr:   a.DetailAddr,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

// FromEntity 从实体转换
func (a *Address) FromEntity(address *entity.Address) {
	a.ID = address.ID
	a.UserID = address.UserID
	a.UserName = address.UserName
	a.UserPhone = address.UserPhone
	a.IsDefault = address.IsDefault
	a.ProvinceName = address.ProvinceName
	a.CityName = address.CityName
	a.RegionName = address.RegionName
	a.DetailAddr = address.DetailAddr
	a.CreatedAt = address.CreatedAt
	a.UpdatedAt = address.UpdatedAt
}
