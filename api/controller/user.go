package controller

import (
	"errors"
	"github.com/Cospk/go-mall/api/request"
	"github.com/Cospk/go-mall/internal/logic/service"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/Cospk/go-mall/pkg/response"
	"github.com/gin-gonic/gin"
)

// LoginUser 登录
func LoginUser(c *gin.Context) {
	// 绑定请求体的参数
	var userLogin request.UserLogin
	if err := c.ShouldBindJSON(&userLogin.Body); err != nil {
		response.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	// 绑定请求头的参数
	if err := c.ShouldBindHeader(&userLogin.Header); err != nil {
		response.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	// 查询用户信息是否存在
	userSvc := service.NewUserService(c)
	err := userSvc.UserLogin(&userLogin)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNotFound) {
			response.NewResponse(c).Error(errcode.ErrUserNotFound)
		} else if errors.Is(err, errcode.ErrUserPasswordError) {
			response.NewResponse(c).Error(errcode.ErrUserPasswordError)
		}
		logger.NewLogger(c).Error("UserLoginError", "err", err)
		return
	}
	response.NewResponse(c).Success("登入成功，欢迎回来")
	return
}
