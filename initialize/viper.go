package initialize

import (
	"fmt"
	"github.com/Cospk/go-mall/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

func InitConfig() {
	// 使用viper加载配置信息
	config := viper.New()
	dir, _ := os.Getwd()
	config.AddConfigPath(dir + "/config")
	config.SetConfigName("application.env")
	config.SetConfigType("yaml")

	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已修改:", e.Name)
		if err := config.Unmarshal(&global.Config); err != nil {
			fmt.Println(err)
		}
	})

	// 将配置文件内容解析到config结构体中
	if err := config.Unmarshal(&global.Config); err != nil {
		fmt.Println(err)
	}
}
