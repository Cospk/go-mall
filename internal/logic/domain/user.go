package domain

import (
	"github.com/Cospk/go-mall/internal/dal/dao"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type UserDomainService struct {
	ctx     *gin.Context
	userDao *dao.UserDao
}

func NewUserDomainService(ctx *gin.Context) *UserDomainService {
	return &UserDomainService{
		ctx:     ctx,
		userDao: dao.NewUserDao(ctx),
	}
}

func (srv *UserDomainService) LoginUser(Name, password string) error {
	existedUser, err := srv.userDao.FindUserByName(Name)
	if err != nil {
		return errcode.Wrap("UserDomainSvcLoginUserError", err)
	}
	if existedUser.Password == password {
		return errcode.ErrUserPasswordError
	}
	return nil
}
