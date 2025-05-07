package dto

import "time"

// http 请求结构体

// OrderCreate 订单创建请求
type OrderCreate struct {
	CartItemIdList []int64 `json:"cart_item_id_list" binding:"required"`
	UserAddressId  int64   `json:"user_address_id" binding:"required"`
}

// OrderPayCreate 订单发起支付请求
type OrderPayCreate struct {
	OrderNo string `json:"order_no" binding:"required"`
	PayType string `json:"pay_type" binding:"required,oneof= 1 2"`
}

// WxPayNotifyRequest 微信支付回调通知请求
// https://pay.weixin.qq.com/docs/merchant/apis/jsapi-payment/payment-notice.html
type WxPayNotifyRequest struct {
	Header struct {
		Timestamp string `json:"Wechatpay-Timestamp"`
		Nonce     string `json:"Wechatpay-Nonce"`
		Signature string `json:"Wechatpay-Signature""`
	}
	Body struct {
		ID           string    `json:"id"`
		CreateTime   time.Time `json:"create_time"`
		ResourceType string    `json:"resource_type"`
		EventType    string    `json:"event_type"`
		Summary      string    `json:"summary"`
		Resource     struct {
			OriginalType   string `json:"original_type"`
			Algorithm      string `json:"algorithm"`
			Ciphertext     string `json:"ciphertext"`
			AssociatedData string `json:"associated_data"`
			Nonce          string `json:"nonce"`
		} `json:"resource"`
	}
}

// http 响应结构体

type OrderCreateReply struct {
	OrderNo string `json:"order_no"`
}

type Order struct {
	OrderNo     string `json:"order_no"`
	PayTransId  string `json:"pay_trans_id"`
	PayType     int    `json:"pay_type"`
	BillMoney   int    `json:"bill_money"`
	PayMoney    int    `json:"pay_money"`
	PayState    int    `json:"pay_state"`
	OrderStatus int    `json:"-"`
	FrontStatus string `json:"status"`
	Address     struct {
		UserName      string `json:"user_name"`
		UserPhone     string `json:"user_phone"`
		ProvinceName  string `json:"province_name"`
		CityName      string `json:"city_name"`
		RegionName    string `json:"region_name"`
		DetailAddress string `json:"detail_address"`
	} `json:"address,omitempty"`
	Items []struct {
		CommodityId           int64  `json:"commodity_id"`
		CommodityName         string `json:"commodity_name"`
		CommodityImg          string `json:"commodity_img"`
		CommoditySellingPrice int    `json:"commodity_selling_price"`
		CommodityNum          int    `json:"commodity_num"`
	} `json:"items,omitempty"`
	CreatedAt string `json:"created_at"`
}
