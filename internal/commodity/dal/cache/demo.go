package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Cospk/go-mall/internal/demo/logic/do"
	"github.com/Cospk/go-mall/pkg/enum"
	"github.com/Cospk/go-mall/pkg/logger"
)

// SetDemoOrder 缓存demoOrder
func SetDemoOrder(ctx context.Context, demoOrder *do.DemoOrder) error {
	jsonDataBytes, _ := json.Marshal(demoOrder)
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, demoOrder.OrderNo)
	_, err := Redis().Set(ctx, redisKey, jsonDataBytes, 0).Result()
	if err != nil {
		logger.NewLogger(ctx).Error("redis error", "err", err)
		return err
	}
	return nil
}

// GetDemoOrder 获取demoOrder
func GetDemoOrder(ctx context.Context, orderNo string) (*do.DemoOrder, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, orderNo)
	jsonBytes, err := Redis().Get(ctx, redisKey).Bytes()
	if err != nil {
		logger.NewLogger(ctx).Error("redis error", "err", err)
		return nil, err
	}
	var demoOrder do.DemoOrder
	_ = json.Unmarshal(jsonBytes, &demoOrder)
	return &demoOrder, nil
}

// DummyDemoOrder 演示一下hash的使用,结构体必须加tag否则报错，go-redis使用hset存储结构体网上也有讨论：https://github.com/redis/go-redis/discussions/2454
type DummyDemoOrder struct {
	OrderNo string `redis:"orderNo"`
	UserId  int64  `redis:"userId"`
}

func SetDemoOrderStruct(ctx context.Context, demoOrder do.DemoOrder) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, demoOrder.OrderNo)
	data := DummyDemoOrder{
		OrderNo: demoOrder.OrderNo,
		UserId:  demoOrder.UserId,
	}
	_, err := Redis().HSet(ctx, redisKey, data).Result()
	if err != nil {
		logger.NewLogger(ctx).Error("redis error", "err", err)
		return err
	}
	return nil
}
