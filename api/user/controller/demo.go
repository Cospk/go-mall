package controller

import (
	"errors"
	"github.com/Cospk/go-mall/api/demo/request"
	service2 "github.com/Cospk/go-mall/internal/demo/logic/service"
	"github.com/Cospk/go-mall/pkg/config"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/Cospk/go-mall/pkg/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 演示Demo，后期使用删除
// controller层 --> 路由层，只负责处理请求和响应，然后进行参数校验，不处理业务逻辑
// HTTP 请求 → controller → service（应用服务） → domain（领域对象/服务） → dal（数据访问）
// 上面的流程中的应用服务和domain的分层价值：避免业务逻辑散落在 controller 或 dao 中，提升代码可维护性和可测试性

// TestPing 测试接口是否正常
func TestPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
	return
}

// TestConfigRead 测试读取配置文件
func TestConfigRead(c *gin.Context) {
	database := config.Database
	c.JSON(200, gin.H{
		"type":          database.Master.Type,
		"dsn":           database.Master.DSN,
		"max_open_conn": database.Master.MaxOpenConn,
		"max_idle_conn": database.Master.MaxIdleConn,
	})
	return
}

// TestLog 测试日志
func TestLog(c *gin.Context) {
	logger.NewLogger(c).Info("test log", "key", "value", "val", 2)
	c.JSON(200, gin.H{
		"status": "ok",
	})
	return
}

// TestAccessLog 测试访问日志
func TestAccessLog(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
	return
}

// TestPanicLog 测试panic日志
func TestPanicLog(c *gin.Context) {
	var a map[string]string
	a["test"] = "test"
	c.JSON(200, gin.H{
		"status": "ok",
		"data":   a,
	})
	return
}

// TestAppError 测试应用错误
func TestAppError(c *gin.Context) {
	// 使用Wrap()函数包装错误
	err := errors.New("生成一个错误")
	appErr := errcode.Wrap("对错误进行包装", err)
	appErr2 := errcode.Wrap("再次对错误进行包装", appErr)
	logger.NewLogger(c).Error("记录错误测试", "err", appErr2)

	// 使用WithDetails()函数添加错误详情
	err = errors.New("生成一个错误")
	apiErr := errcode.ErrServer.WithCause(err)
	logger.NewLogger(c).Error("执行中出现错误", "err", apiErr)

	c.JSON(apiErr.HttpStatusCode(), gin.H{
		"code": apiErr.Code(),
		"msg":  apiErr.Msg(),
	})
}

// TestResponseObj 测试响应
func TestResponseObj(c *gin.Context) {
	data := map[string]any{
		"a": "test",
		"b": 12,
	}
	resp.NewResponse(c).Success(data)
	return
}

// TestResponseList 测试响应列表
func TestResponseList(c *gin.Context) {
	pageInfo := resp.GetPageInfo(c)

	data := []struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		{Name: "张三", Age: 12},
		{Name: "李四", Age: 13},
	}
	pageInfo.Total = 2
	resp.NewResponse(c).SetPageInfo(pageInfo).Success(data)
	return

}

// TestResponseError 测试响应错误
func TestResponseError(c *gin.Context) {
	baseErr := errors.New("测试错误")
	// 下面这个在正式开发时写在service层
	err := errcode.Wrap("encountered an error when xxx service did xxx", baseErr)
	resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
	return
}

// TestGormLogger 测试gorm日志是否会输出到项目日志中
func TestGormLogger(c *gin.Context) {
	svc := service2.NewDemoSvc(c)
	list, err := svc.GetDemoList()
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	resp.NewResponse(c).Success(list)
}

// TestCreateDemoOrder 创建demo订单
func TestCreateDemoOrder(c *gin.Context) {
	req := new(request.DemoOrderCreate)
	err := c.ShouldBindJSON(req)
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	// TODO 验证token，这个后面集成jwt后再写

	order, err2 := service2.NewDemoSvc(c).CreateDemoOrder(req)
	if err2 != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err2))
		return
	}
	resp.NewResponse(c).Success(order)
}

func TestGetToken(c *gin.Context) {
	userSvc := service2.NewUserService(c)
	token, err := userSvc.GetToken()
	if err != nil {
		if errors.Is(err, errcode.ErrUserInvalid) {
			logger.NewLogger(c).Error("invalid user is unable to generate token", err)
			//app.NewResponse(c).Error(errcode.ErrUserInvalid.WithCause(err)) 第一版的AppError会导Error循环引用，现在已解决
			resp.NewResponse(c).Error(errcode.ErrUserInvalid)
		} else {
			appErr := err.(*errcode.AppError)
			resp.NewResponse(c).Error(appErr)
		}
		return
	}
	resp.NewResponse(c).Success(token)
}

func TestRefreshToken(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	if refreshToken == "" {
		resp.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	userSvc := service2.NewUserService(c)
	token, err := userSvc.TokenRefresh(refreshToken)
	if err != nil {
		if errors.Is(err, errcode.ErrTooManyRequests) {
			// 客户端有并发刷新token
			resp.NewResponse(c).Error(errcode.ErrTooManyRequests)
			return
		} else {
			appErr := err.(*errcode.AppError)
			resp.NewResponse(c).Error(appErr)
		}
		return
	}
	resp.NewResponse(c).Success(token)
}

func TestVerifyToken(c *gin.Context) {
	resp.NewResponse(c).Success(gin.H{
		"user_id":    c.GetInt64("userId"),
		"session_id": c.GetString("sessionId"),
	})
	return
}
