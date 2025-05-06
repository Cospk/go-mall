package router

import (
	"github.com/Cospk/go-mall/internal/api/infrastructure/rpc"
	"github.com/Cospk/go-mall/internal/api/interfaces/http/handler"
	"github.com/Cospk/go-mall/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func InitApiRouter() *gin.Engine {
	r := gin.Default()

	// 注册全局中间件
	r.Use(middleware.AuthMiddleware())
	r.Use(middleware.TraceMiddleware())
	//r.Use(middleware.Cors())
	//r.Use(middleware.Logger())
	//r.Use(middleware.Recovery())

	// 初始化RPC客户端
	// 用户服务客户端
	userClient := rpc.NewUserServiceClient()
	// 商品服务客户端
	commodityClient := rpc.NewCommodityServiceClient()
	// 订单服务客户端
	orderClient := rpc.NewOrderServiceClient()
	// 购物车服务客户端
	cartClient := rpc.NewCartServiceClient()

	// 初始化处理器
	// 用户处理器
	userHandler := handler.NewUserHandler(userClient)
	// 商品处理器
	commodityHandler := handler.NewCommodityHandler(commodityClient)
	// 订单处理器
	orderHandler := handler.NewOrderHandler(orderClient)
	// 购物车处理器
	cartHandler := handler.NewCartHandler(cartClient)

	// 注册API路由
	api := r.Group("/api")
	{
		// 用户相关路由
		userGroup := api.Group("/user")
		{
			userGroup.POST("/register", userHandler.Register)
			userGroup.POST("/login", userHandler.Login)
			userGroup.GET("/token/refresh", userHandler.RefreshUserToken)
			userGroup.POST("/password/apply-reset", userHandler.PasswordResetApply)
			userGroup.POST("/password/reset", userHandler.PasswordReset)

			// 需要认证的路由
			authUserGroup := userGroup.Group("/")
			authUserGroup.Use(middleware.AuthMiddleware())
			{
				authUserGroup.DELETE("/logout", userHandler.LogoutUser)
				authUserGroup.GET("/info", userHandler.GetUserInfo)
				authUserGroup.PATCH("/info", userHandler.UpdateUserInfo)
				authUserGroup.GET("/address", userHandler.GetUserAddresses)
				authUserGroup.POST("/address", userHandler.AddUserAddress)
				authUserGroup.GET("/address/:address_id", userHandler.GetUserAddress)
				authUserGroup.PATCH("/address/:address_id", userHandler.UpdateUserAddress)
				authUserGroup.DELETE("/address/:address_id", userHandler.DeleteUserAddress)
			}
		}

		// 商品相关路由
		commodityGroup := api.Group("/commodity")
		{
			commodityGroup.GET("/category-hierarchy", commodityHandler.GetCommodityList)
			commodityGroup.GET("/category", commodityHandler.GetCommodityDetail)
			commodityGroup.GET("/commodity-in-cate", commodityHandler.GetCommodityDetail)
			commodityGroup.GET("/search", commodityHandler.GetCommodityDetail)
			commodityGroup.GET(":commodity_id/info", commodityHandler.GetCommodityDetail)
		}

		// 订单相关路由
		orderGroup := api.Group("/order")
		orderGroup.Use(middleware.AuthMiddleware())
		{
			orderGroup.POST("/create", orderHandler.CreateOrder)
			orderGroup.GET("/user-order", orderHandler.GetOrderList)
			orderGroup.GET("/:order_no/info", orderHandler.GetOrderDetail)
			orderGroup.PATCH("/:order_no/cancel", orderHandler.GetOrderDetail)
			orderGroup.POST("create-pay", orderHandler.GetOrderDetail)
		}

		// 购物车相关路由
		cartGroup := api.Group("/cart")
		cartGroup.Use(middleware.AuthMiddleware())
		{
			cartGroup.POST("/add-item", cartHandler.AddToCart)
			cartGroup.GET("/update-item", cartHandler.GetCartList)
			cartGroup.DELETE("/item", cartHandler.RemoveFromCart)
		}

		// 评价相关路由

		// 消息相关路由

		// 活动相关路由
	}

	return r
}
