package controller

import (
	"errors"
	"github.com/Cospk/go-mall/api/response"
	"github.com/Cospk/go-mall/internal/logic/service"
	"github.com/Cospk/go-mall/pkg/config"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 演示Demo，后期使用删除
// controller层 --> 路由层，只负责处理请求和响应，然后进行参数校验，不处理业务逻辑
// HTTP 请求 → controller → service（应用服务） → domain（领域对象/服务） → dal（数据访问）
// 上面的流程中的应用服务和domain的分层价值：避免业务逻辑散落在 controller 或 dao 中，提升代码可维护性和可测试性

func TestPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
	return
}

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

func TestLog(c *gin.Context) {
	logger.NewLogger(c).Info("test log", "key", "value", "val", 2)
	c.JSON(200, gin.H{
		"status": "ok",
	})
	return
}
func TestAccessLog(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
	return
}
func TestPanicLog(c *gin.Context) {
	var a map[string]string
	a["test"] = "test"
	c.JSON(200, gin.H{
		"status": "ok",
		"data":   a,
	})
	return
}
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
func TestResponseObj(c *gin.Context) {
	data := map[string]any{
		"a": "test",
		"b": 12,
	}
	response.NewResponse(c).Success(data)
	return
}
func TestResponseList(c *gin.Context) {
	pageInfo := response.GetPageInfo(c)

	data := []struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		{Name: "张三", Age: 12},
		{Name: "李四", Age: 13},
	}
	pageInfo.Total = 2
	response.NewResponse(c).SetPageInfo(pageInfo).Success(data)
	return

}
func TestResponseError(c *gin.Context) {
	baseErr := errors.New("测试错误")
	// 下面这个在正式开发时写在service层
	err := errcode.Wrap("encountered an error when xxx service did xxx", baseErr)
	response.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
	return
}

func TestGormLogger(c *gin.Context) {
	svc := service.NewDemoSvc(c)
	list, err := svc.GetDemoList()
	if err != nil {
		response.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	response.NewResponse(c).Success(list)
}
