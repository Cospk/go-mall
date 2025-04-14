package external

import (
	"context"
	"encoding/json"
	"github.com/Cospk/go-mall/api/request"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/Cospk/go-mall/pkg/utils/httptool"
)

// external 层 对接外部系统api,比如服务端可能会调用服务的api查询用户在线数量

type DemoLib struct {
	ctx context.Context
}

// NewDemoLib 创建时上层通过ctx 把 gin.Ctx传递过来
func NewDemoLib(ctx context.Context) *DemoLib {
	return &DemoLib{ctx: ctx}
}

type OrderCreateResult struct {
	UserId    int64  `json:"user_id"`
	BillMoney int64  `json:"bill_money"`
	OrderNo   string `json:"order_no"`
	State     int8   `json:"state"`
	PaidAt    string `json:"paid_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// 用http调自己项目里的POST接口, 演示用, 实际使用时不要这么干

func (lib *DemoLib) TestPostCreateOrder() (*OrderCreateResult, error) {
	data := &request.DemoOrderCreate{
		UserId:       12345,
		BillMoney:    20,
		OrderGoodsId: 1111110,
	}
	jsonReq, _ := json.Marshal(data)
	httCode, respBody, err := httptool.Post(lib.ctx, "http://localhost:8080/building/create-demo-order", jsonReq)

	logger.NewLogger(lib.ctx).Info("create-demo-order api reply ", "code", httCode, "data", respBody, "err", err)

	if err != nil {
		return nil, err
	}

	reply := &struct {
		Code int `json:"code"`
		Data *OrderCreateResult
	}{}
	json.Unmarshal(respBody, reply)
	return reply.Data, nil
}
