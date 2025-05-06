package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
)

// InitConfig 初始化配置信息 (配置通用化)
// serviceName: 服务名称，用于指定配置文件名称
// configs: 配置映射表，用于将配置信息映射到结构体中
func InitConfig(serviceName string, configs map[string]interface{}) {
	// 使用viper加载配置信息
	config := viper.New()
	dir, _ := os.Getwd()

	config.AddConfigPath(dir + "/configs")
	switch serviceName {
	case "api":
		config.SetConfigName("api-server")
	case "user":
		config.SetConfigName("user-server")
	default:
		config.SetConfigName("api-server")
	}
	config.SetConfigType("yaml")

	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 读取的配置信息写到结构体中
	err := parseConfig(config, configs)
	if err != nil {
		// 解析配置文件失败,应该直接panic并抛出错误，否则其他初始化基本无法进行
		panic(err)
	}

	// 监听配置文件变化
	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已修改:", e.Name)
		err = parseConfig(config, configs)
		if err != nil {
			panic(err)
		}
	})

}

func parseConfig(config *viper.Viper, configs map[string]interface{}) error {
	// 如果提供了配置映射表，则使用映射表
	if configs != nil && len(configs) > 0 {
		for key, val := range configs {
			// 检查是否是指针类型
			if reflect.TypeOf(val).Kind() != reflect.Ptr {
				return fmt.Errorf("配置必须是指针类型: %s", key)
			}

			if err := config.UnmarshalKey(key, val); err != nil {
				return fmt.Errorf("解析配置失败 [%s]: %w", key, err)
			}
		}
		return nil
	}
	return error(nil)
}
