package main

import "github.com/Cospk/go-mall/initialize"

func main() {

	// 初始化配置
	initialize.InitConfig()

	// 初始化日志
	initialize.InitLogger()

	// 初始化路由
	Router := initialize.InitWebRouter()

	Router.Run(":8080")
}
