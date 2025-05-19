package entity

import "time"

// Address 地址实体
type Address struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	UserName     string    `json:"user_name"`
	UserPhone    string    `json:"user_phone"`
	IsDefault    int32     `json:"is_default"`
	ProvinceName string    `json:"province_name"`
	CityName     string    `json:"city_name"`
	RegionName   string    `json:"region_name"`
	DetailAddr   string    `json:"detail_addr"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewAddress 创建新地址
func NewAddress(userID int64, userName, userPhone string, isDefault int32, provinceName, cityName, regionName, detailAddr string) *Address {
	now := time.Now()
	return &Address{
		UserID:       userID,
		UserName:     userName,
		UserPhone:    userPhone,
		IsDefault:    isDefault,
		ProvinceName: provinceName,
		CityName:     cityName,
		RegionName:   regionName,
		DetailAddr:   detailAddr,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Update 更新地址信息
func (a *Address) Update(userName, userPhone string, isDefault int32, provinceName, cityName, regionName, detailAddr string) {
	if userName != "" {
		a.UserName = userName
	}
	if userPhone != "" {
		a.UserPhone = userPhone
	}
	a.IsDefault = isDefault
	if provinceName != "" {
		a.ProvinceName = provinceName
	}
	if cityName != "" {
		a.CityName = cityName
	}
	if regionName != "" {
		a.RegionName = regionName
	}
	if detailAddr != "" {
		a.DetailAddr = detailAddr
	}
	a.UpdatedAt = time.Now()
}

// SetDefault 设置为默认地址
func (a *Address) SetDefault() {
	a.IsDefault = 1
	a.UpdatedAt = time.Now()
}

// UnsetDefault 取消默认地址
func (a *Address) UnsetDefault() {
	a.IsDefault = 0
	a.UpdatedAt = time.Now()
}
