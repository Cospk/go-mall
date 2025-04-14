package service

import (
	"github.com/Cospk/go-mall/api/reply"
	"github.com/Cospk/go-mall/api/request"
	"github.com/Cospk/go-mall/internal/logic/do"
	"github.com/Cospk/go-mall/internal/logic/domain"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/Cospk/go-mall/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	ctx        *gin.Context
	userDomain *domain.UserDomain
}

func NewUserService(ctx *gin.Context) *UserService {
	return &UserService{
		ctx:        ctx,
		userDomain: domain.NewUserDomain(ctx),
	}
}

func (svc *UserService) GetToken() (*reply.TokenReply, error) {
	token, err := svc.userDomain.GenAuthToken(12345678, "h5", "")
	if err != nil {
		return nil, err
	}
	logger.NewLogger(svc.ctx).Info("generate token success", "tokenData", token)
	tokenReply := new(reply.TokenReply)
	_ = utils.CopyStruct(tokenReply, token)

	return tokenReply, err
}

func (svc *UserService) TokenRefresh(refreshToken string) (*reply.TokenReply, error) {
	token, err := svc.userDomain.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	logger.NewLogger(svc.ctx).Info("refresh token success", "tokenData", token)
	tokenReply := new(reply.TokenReply)
	_ = utils.CopyStruct(tokenReply, token)
	return tokenReply, err
}

// UserRegister 用户注册
func (svc *UserService) UserRegister(userRegisterReq *request.UserRegister) error {
	userInfo := new(do.UserBaseInfo)
	utils.CopyStruct(userInfo, userRegisterReq)

	// 领域服务注册用户
	_, err := svc.userDomain.RegisterUser(userInfo, userRegisterReq.Password)
	if errcode.Is(err, errcode.ErrUserNameOccupied) {
		return err
	}

	// TODO 写注册成功后的外围辅助逻辑, 比如注册成功后给用户发确认邮件|短信

	// TODO 如果产品逻辑是注册后帮用户登录, 那这里再掉登录的逻辑

	return nil
}

// UserLogin 用户登录
func (svc *UserService) UserLogin(userLoginReq *request.UserLogin) (*reply.TokenReply, error) {
	tokenInfo, err := svc.userDomain.LoginUser(userLoginReq.Body.LoginName, userLoginReq.Body.Password, userLoginReq.Header.Platform)
	if err != nil {
		return nil, err

	}
	tokenReply := new(reply.TokenReply)
	utils.CopyStruct(tokenReply, tokenInfo)

	// TODO 执行登录后的业务逻辑
	return tokenReply, nil
}

func (svc *UserService) UserLogout(userId int64, platform string) error {
	err := svc.userDomain.LogoutUser(userId, platform)
	return err
}
