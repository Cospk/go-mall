package handler

import (
	"errors"
	"github.com/Cospk/go-mall/internal/api/application/dto"
	"github.com/Cospk/go-mall/internal/api/application/service"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/resp"
	"github.com/Cospk/go-mall/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userClient interface{}) *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(userClient),
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.UserRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	// 调用应用服务
	userID, err := h.userService.Register(c, req)
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	resp.NewResponse(c).Success(gin.H{"user_id": userID})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	// 调用应用服务
	token, err := h.userService.Login(c, req)
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	resp.NewResponse(c).Success(gin.H{"token": token})
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := c.GetInt64("userID")

	// 调用应用服务
	user, err := h.userService.GetUserInfo(c, userID)
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	resp.NewResponse(c).Success(gin.H{"user": user})

}

func (h *UserHandler) RefreshUserToken(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	if refreshToken == "" {
		resp.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	token, err := h.userService.RefreshUserToken(c, refreshToken)
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

func (h *UserHandler) LoginUser(c *gin.Context) {
	loginRequest := new(dto.UserLogin)
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	if err := c.ShouldBindHeader(&loginRequest); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	// 登录用户
	token, err := h.userService.Login(c, *loginRequest)
	if err != nil {
		if errors.Is(err, errcode.ErrUserInvalid) {
			resp.NewResponse(c).Error(errcode.ErrUserInvalid)
			return
		}
	}

	resp.NewResponse(c).Success(token)
	return
}

func (h *UserHandler) LogoutUser(c *gin.Context) {
	userId := c.GetInt64("userId")
	platform := c.GetString("platform")
	err := h.userService.Logout(c, userId, platform)
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	resp.NewResponse(c).SuccessOk()
}

// PasswordResetApply 申请重置密码
func (h *UserHandler) PasswordResetApply(c *gin.Context) {
	request := new(dto.PasswordResetApply)
	if err := c.ShouldBindJSON(request); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	apply, err := h.userService.PasswordResetApply(c, *request)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNotRight) {
			resp.NewResponse(c).Error(errcode.ErrUserNotRight)
		} else {
			resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	resp.NewResponse(c).Success(apply)
}

func (h *UserHandler) PasswordReset(c *gin.Context) {
	request := new(dto.PasswordReset)
	if err := c.ShouldBindJSON(request); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	if !utils.PasswordComplexityVerify(request.Password) {
		// Validator验证通过后再应用 密码复杂度这样的特殊验证
		resp.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	err := h.userService.PasswordReset(c, *request)
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
func (h *UserHandler) UserInfo(c *gin.Context) {
	userId := c.GetInt64("userId")
	info, err := h.userService.GetUserInfo(c, userId)
	if err == nil {
		resp.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	resp.NewResponse(c).Success(info)
}

// UpdateUserInfo 个人信息更新
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	request := new(dto.UserInfoUpdate)
	if err := c.ShouldBindJSON(request); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	err := h.userService.UpdateUserInfo(c, *request)
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	resp.NewResponse(c).SuccessOk()
}

// AddUserAddress 新增收货地址
func (h *UserHandler) AddUserAddress(c *gin.Context) {
	request := new(dto.UserAddress)
	if err := c.ShouldBindJSON(request); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	_, err := h.userService.AddUserAddressInfo(c, *request)
	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	resp.NewResponse(c).SuccessOk()
}

// GetUserAddresses 获取用户的收货信息列表
func (h *UserHandler) GetUserAddresses(c *gin.Context) {
	list, err := h.userService.GetUserAddressList(c, c.GetInt64("userId"))

	if err != nil {
		resp.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	resp.NewResponse(c).Success(list)
}

// GetUserAddress 获取用户单个收货信息
func (h *UserHandler) GetUserAddress(c *gin.Context) {
	addressId, _ := strconv.ParseInt(c.Param("address_id"), 10, 64)
	if addressId <= 0 {
		resp.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	replyData, err := h.userService.GetUserAddressInfo(c, c.GetInt64("userId"))

	if err != nil {
		if errors.Is(err, errcode.ErrParams) {
			resp.NewResponse(c).Error(errcode.ErrParams)
		} else {
			resp.NewResponse(c).Error(errcode.ErrServer)
		}
		return
	}
	resp.NewResponse(c).Success(replyData)
}

// UpdateUserAddress 修改用户地址信息
func (h *UserHandler) UpdateUserAddress(c *gin.Context) {
	// 验证URL中的参数
	addressId, _ := strconv.ParseInt(c.Param("address_id"), 10, 64)
	if addressId <= 0 {
		resp.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	// 验证请求Body中的参数
	requestData := new(dto.UserAddress)
	if err := c.ShouldBindJSON(requestData); err != nil {
		resp.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	err := h.userService.UpdateUserAddressInfo(c, addressId, *requestData)

	if err != nil {
		if errors.Is(err, errcode.ErrParams) {
			resp.NewResponse(c).Error(errcode.ErrParams)
		} else {
			resp.NewResponse(c).Error(errcode.ErrServer)
		}
		return
	}
	resp.NewResponse(c).SuccessOk()
}

// DeleteUserAddress 删除用户地址信息
func (h *UserHandler) DeleteUserAddress(c *gin.Context) {
	// 验证URL中的参数
	addressId, _ := strconv.ParseInt(c.Param("address_id"), 10, 64)
	if addressId <= 0 {
		resp.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	err := h.userService.DeleteUserAddressInfo(c, addressId)
	if err != nil {
		if errors.Is(err, errcode.ErrParams) {
			resp.NewResponse(c).Error(errcode.ErrParams)
		} else {
			resp.NewResponse(c).Error(errcode.ErrServer)
		}
		return
	}
	resp.NewResponse(c).SuccessOk()
}
