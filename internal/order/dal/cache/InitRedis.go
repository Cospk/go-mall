package cache

import (
	"context"
	"github.com/Cospk/go-mall/pkg/config"
	"github.com/redis/go-redis/v9"
	"time"
)

var RedisClient *redis.Client

func Redis() *redis.Client {
	return RedisClient
}

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         config.Redis.Address,
		Password:     config.Redis.Password,
		PoolSize:     config.Redis.PoolSize, // 连接池大小
		DB:           config.Redis.DB,       // use default DB
		DialTimeout:  10 * time.Second,      // 连接超时
		ReadTimeout:  30 * time.Second,      // 读取超时
		WriteTimeout: 30 * time.Second,      // 写入超时
		PoolTimeout:  30 * time.Second,      // 当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒
	})

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		// 连接失败停止程序，若是缓存无影响可注释掉
		panic(err)
	}
}

// 注：
//	1、redis没有提供对应的接口，就无法像gorm定制Logger，只能在Redis的存取中添加日志
//	2、redis的hash功能（HSET命令）不是很完善，HSET存储结构体每一个字段必须加tag，结构体嵌套也不支持，为此一般是结构体json化后存储，读取后再解析出来
