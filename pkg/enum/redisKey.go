package enum

// Redis Key的格式规范
// 项目名:模块名:键名

// DEMO模块

const (
	REDIS_KEY_DEMO_ORDER_DETAIL = "GOMALL:DEMO:ORDER_DETAIL_%s"
)

// 用户模块

const (
	REDIS_KEY_ACCESS_TOKEN = "GOMALL:USER:ACCESS_TOKEN_%s"
)
