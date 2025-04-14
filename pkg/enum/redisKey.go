package enum

// Redis Key的格式规范
// 项目名:模块名:键名

// DEMO模块

const (
	REDIS_KEY_DEMO_ORDER_DETAIL = "GOMALL:DEMO:ORDER_DETAIL_%s"
)

// 用户模块

// Token相关的Redis键名模板
const (
	REDIS_KEY_ACCESS_TOKEN       = "token:access:%s"       // 访问Token的Redis键名模板
	REDIS_KEY_REFRESH_TOKEN      = "token:refresh:%s"      // 刷新Token的Redis键名模板
	REDIS_KEY_USER_SESSION       = "user:session:%d"       // 用户会话信息的Redis键名模板
	REDISKEY_TOKEN_REFRESH_LOCK  = "token:refresh:lock:%s" // 刷新Token的锁
	REDISKEY_PASSWORDRESET_TOKEN = "token:pwdreset:%s"     // 密码重置Token
)
