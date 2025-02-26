package service

import (
	"github.com/Cospk/go-mall/api/request"
	"github.com/Cospk/go-mall/internal/logic/domain"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	ctx           *gin.Context
	userDomainSvc *domain.UserDomainService
}

func NewUserService(ctx *gin.Context) *UserService {
	return &UserService{
		ctx:           ctx,
		userDomainSvc: domain.NewUserDomainService(ctx),
	}
}

// UserLogin 用户登录
func (srv *UserService) UserLogin(userLoginReq *request.UserLogin) error {
	err := srv.userDomainSvc.LoginUser(userLoginReq.Body.LoginName, userLoginReq.Body.Password)
	if err != nil {
		return err

	}
	// TODO 执行登录后的业务逻辑
	return nil
}
