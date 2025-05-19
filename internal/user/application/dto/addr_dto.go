package dto

// AddAddressRequest 添加地址请求
type AddAddressRequest struct {
	UserName     string
	UserPhone    string
	Default      int32
	ProvinceName string
	CityName     string
	RegionName   string
	DetailAddr   string
}

// AddAddressResponse 添加地址响应
type AddAddressResponse struct {
	AddressID int64
	Message   string
}

// AddressDTO 地址DTO
type AddressDTO struct {
	ID              int64
	UserName        string
	UserPhone       string
	MaskedUserName  string
	MaskedUserPhone string
	Default         int32
	ProvinceName    string
	CityName        string
	RegionName      string
	DetailAddress   string
	CreatedAt       string
}

// UpdateAddressRequest 更新地址请求
type UpdateAddressRequest struct {
	AddressID     int64
	UserName      string
	UserPhone     string
	Default       int32
	ProvinceName  string
	CityName      string
	RegionName    string
	DetailAddress string
}

// UpdateAddressResponse 更新地址响应
type UpdateAddressResponse struct {
	Success bool
	Message string
}

// DeleteAddressResponse 删除地址响应
type DeleteAddressResponse struct {
	Success bool
	Message string
}
