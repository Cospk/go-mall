package redis

import (
	"context"

	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

// 默认配置常量
const (
	defaultBatchSize       = 50 // 默认批量处理大小
	defaultConcurrentLimit = 3  // 默认并发限制数
)

// RedisShardManager 用于分片和处理Redis键的管理器
type RedisShardManager struct {
	redisClient redis.UniversalClient
	config      *Config
}

// Config 包含分片处理的配置参数,用于实现灵活配置模式（函数式选项模式）
type Config struct {
	batchSize       int
	continueOnError bool
	concurrentLimit int
}

// Option 是用于配置Config的函数类型
type Option func(c *Config)

// groupKeysBySlot  集群键分槽算法,将输入键按Redis集群哈希槽分组，返回槽位到键列表的映射
func groupKeysBySlot(ctx context.Context, redisClient redis.UniversalClient, keys []string) (map[int64][]string, error) {
	// 创建哈希槽到键的映射
	slots := make(map[int64][]string)
	// 检查是否为集群客户端
	clusterClient, isCluster := redisClient.(*redis.ClusterClient)
	if isCluster && len(keys) > 1 {
		// 创建管道用于批量操作
		pipe := clusterClient.Pipeline()
		// 创建命令切片
		cmds := make([]*redis.IntCmd, len(keys))
		for i, key := range keys {
			cmds[i] = pipe.ClusterKeySlot(ctx, key)
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			return nil, err
		}

		//  解析槽位分配结果
		for i, cmd := range cmds {
			slot, err := cmd.Result()
			if err != nil {
				logger.NewLogger(ctx).Error("some key get slot err", "err", err, "key", keys[i])
				return nil, err
			}
			slots[slot] = append(slots[slot], keys[i])
		}
	} else {
		// 单机模式处理：所有键分配到虚拟槽位0
		slots[0] = keys
	}

	return slots, nil
}

// splitIntoBatches 一个字符串切片 keys 按照指定的 batchSize 分割成多个子切片，最终返回一个二维字符串切片，其中每个子切片代表一个批次。
func splitIntoBatches(keys []string, batchSize int) [][]string {
	var batches [][]string
	for batchSize < len(keys) {
		keys, batches = keys[batchSize:], append(batches, keys[0:batchSize:batchSize])
	}
	return append(batches, keys)
}

// ProcessKeysBySlot 按照 Redis 集群的哈希槽进行分组，然后将每个分组内的键再按照指定的批量大小进行分批处理，最后使用用户提供的处理函数对每个批次的键进行处理
func ProcessKeysBySlot(
	ctx context.Context,
	redisClient redis.UniversalClient,
	keys []string,
	processFunc func(ctx context.Context, slot int64, keys []string) error,
	opts ...Option,
) error {

	config := &Config{
		batchSize:       defaultBatchSize,
		continueOnError: false,
		concurrentLimit: defaultConcurrentLimit,
	}
	for _, opt := range opts {
		opt(config)
	}

	// Group keys by slot
	slots, err := groupKeysBySlot(ctx, redisClient, keys)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(config.concurrentLimit)

	// Process keys in each slot using the provided function
	for slot, singleSlotKeys := range slots {
		batches := splitIntoBatches(singleSlotKeys, config.batchSize)
		for _, batch := range batches {
			slot, batch := slot, batch // Avoid closure capture issue
			g.Go(func() error {
				err := processFunc(ctx, slot, batch)
				if err != nil {
					logger.NewLogger(ctx).Error("Batch processFunc failed", "err", err, "slot", slot, "keys", batch)
					if !config.continueOnError {
						return err
					}
				}
				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

// DeleteCacheBySlot 删除缓存
func DeleteCacheBySlot(ctx context.Context, redisClient redis.UniversalClient, keys []string) error {
	switch len(keys) {
	case 0:
		return nil
	case 1:
		return redisClient.Del(ctx, keys[0]).Err()
	default:
		return ProcessKeysBySlot(ctx, redisClient, keys, func(ctx context.Context, slot int64, keys []string) error {
			return redisClient.Del(ctx, keys...).Err()
		})
	}
}
