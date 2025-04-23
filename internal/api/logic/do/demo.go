package do

import "time"

// do 领域对象

// DemoOrder 演示demo
type DemoOrder struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	BillMoney int64     `json:"bill_money"`
	OrderNo   string    `json:"order_no"`
	State     int8      `json:"state"`
	IsDel     uint      `json:"is_del"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
