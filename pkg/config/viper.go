package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

func InitConfig() {
	// 使用viper加载配置信息
	config := viper.New()
	dir, _ := os.Getwd()
	config.AddConfigPath(dir + "/pkg/config")
	config.SetConfigName("application.env")
	config.SetConfigType("yaml")

	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 读取的配置信息写到结构体中
	err := parseConfig(config)
	if err != nil {
		// 解析配置文件失败,应该直接panic并抛出错误，否则其他初始化基本无法进行
		panic(err)
	}

	// 监听配置文件变化
	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已修改:", e.Name)
		err = parseConfig(config)
		if err != nil {
			panic(err)
		}
	})

}

func parseConfig(config *viper.Viper) error {
	if err := config.UnmarshalKey("app", &AppConfig); err != nil {
		return fmt.Errorf("解析AppConfig失败: %w", err)
	}
	if err := config.UnmarshalKey("redis", &Redis); err != nil {
		return fmt.Errorf("解析Redis配置失败: %w", err)
	}
	if err := config.UnmarshalKey("database", &Database); err != nil {
		return fmt.Errorf("解析Database配置失败: %w", err)
	}
	return nil
}
