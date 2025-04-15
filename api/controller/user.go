package controller

import (
	"errors"
	"github.com/Cospk/go-mall/api/request"
	"github.com/Cospk/go-mall/internal/logic/service"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/Cospk/go-mall/pkg/resp"
	"github.com/Cospk/go-mall/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 注册用户
func RegisterUser(ctx *gin.Context) {
	userRequest := new(request.UserRegister)
	if err := ctx.ShouldBind(userRequest); err != nil {
		resp.NewResponse(ctx).Error(errcode.ErrParams.WithCause(err))
		return
	}
	if !utils.PasswordComplexityVerify(userRequest.Password) {
		// Validator验证通过后再应用 密码复杂度这样的特殊验证
		logger.NewLogger(ctx).Warn("RegisterUserError", "err", "密码复杂度不满足", "password", userRequest.Password)
		resp.NewResponse(ctx).Error(errcode.ErrParams)
		return
	}
	// 注册用户
	userSvc := service.NewUserService(ctx)
	err := userSvc.UserRegister(userRequest)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNameOccupied) {
			resp.NewResponse(ctx).Error(errcode.ErrUserNameOccupied)
		} else {
			resp.NewResponse(ctx).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	resp.NewResponse(ctx).SuccessOk()
	return
}

// LoginUser 登录
func LoginUser(c *gin.Context) {
	// 绑定请求体的参数
	var userLogin request.UserLogin
	if err := c.ShouldBindJSON(&userLogin.Body); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	// 绑定请求头的参数
	if err := c.ShouldBindHeader(&userLogin.Header); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	// 查询用户信息是否存在
	userSvc := service.NewUserService(c)
	token, err := userSvc.UserLogin(&userLogin)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNotFound) {
			resp.NewResponse(c).Error(errcode.ErrUserNotFound)
		} else if errors.Is(err, errcode.ErrUserPasswordError) {
			resp.NewResponse(c).Error(errcode.ErrUserPasswordError)
		}
		logger.NewLogger(c).Error("UserLoginError", "err", err)
		return
	}
	resp.NewResponse(c).Success(token)
	return
}

func LogoutUser(c *gin.Context) {
	userId := c.GetInt64("userId")
	platform := c.GetString("platform")
	userSvc := service.NewUserService(c)
	err := userSvc.UserLogout(userId, platform)
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	resp.NewResponse(c).SuccessOk()
}

func RefreshUserToken(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	if refreshToken == "" {
		resp.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	userSvc := service.NewUserService(c)
	token, err := userSvc.TokenRefresh(refreshToken)
	if err != nil {
		if errors.Is(err, errcode.ErrTooManyRequests) {
			// 客户端有并发刷新token
			resp.NewResponse(c).Error(errcode.ErrTooManyRequests)
		} else {
			appErr := err.(*errcode.AppError)
			resp.NewResponse(c).Error(appErr)
		}
		return
	}
	resp.NewResponse(c).Success(token)
}

// PasswordResetApply 申请重置密码
func PasswordResetApply(c *gin.Context) {
	request := new(request.PasswordResetApply)
	if err := c.ShouldBindJSON(request); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	userSvc := service.NewUserAppSvc(c)
	reply, err := userSvc.PasswordResetApply(request)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNotRight) {
			resp.NewResponse(c).Error(errcode.ErrUserNotRight)
		} else {
			resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	resp.NewResponse(c).Success(reply)
}

func PasswordReset(c *gin.Context) {
	request := new(request.PasswordReset)
	if err := c.ShouldBindJSON(request); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	if !utils.PasswordComplexityVerify(request.Password) {
		// Validator验证通过后再应用 密码复杂度这样的特殊验证
		logger.NewLogger(c).Warn("RegisterUserError", "err", "密码复杂度不满足", "password", request.Password)
		resp.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	userSvc := service.NewUserAppSvc(c)
	err := userSvc.PasswordReset(request)
	if err != nil {
		if errors.Is(err, errcode.ErrParams) {
			resp.NewResponse(c).Error(errcode.ErrParams)
		} else if errors.Is(err, errcode.ErrUserInvalid) {
			resp.NewResponse(c).Error(errcode.ErrUserInvalid)
		} else {
			resp.NewResponse(c).Error(errcode.ErrServer)
		}
		return
	}

	resp.NewResponse(c).SuccessOk()
}

// UserInfo 个人信息查询
func UserInfo(c *gin.Context) {
	userId := c.GetInt64("userId")
	userSvc := service.NewUserAppSvc(c)
	userInfoReply := userSvc.UserInfo(userId)
	if userInfoReply == nil {
		resp.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	resp.NewResponse(c).Success(userInfoReply)
}

// UpdateUserInfo 个人信息更新
func UpdateUserInfo(c *gin.Context) {
	request := new(request.UserInfoUpdate)
	if err := c.ShouldBindJSON(request); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	userSvc := service.NewUserAppSvc(c)
	err := userSvc.UserInfoUpdate(request, c.GetInt64("userId"))
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	resp.NewResponse(c).SuccessOk()
}
