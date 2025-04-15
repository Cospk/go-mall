package dao

import (
	"context"
	"errors"
	"github.com/Cospk/go-mall/internal/dal/model"
	"github.com/Cospk/go-mall/internal/logic/do"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/utils"
	"gorm.io/gorm"
)

type UserDao struct {
	ctx context.Context
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{ctx: ctx}
}

func (dao *UserDao) FindUserById(id int64) (*model.User, error) {
	return &model.User{}, nil
}

func (dao *UserDao) FindUserByName(name string) (user model.User, err error) {
	// TODO 执行sql查询数据库的数据
	result := DB().WithContext(dao.ctx).Where("name = ?", name).First(&user)
	if result.RowsAffected == 0 {
		return model.User{}, result.Error
	}
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}

// FindUserByLoginName 根据登录名查询用户
func (dao *UserDao) FindUserByLoginName(loginName string) (*model.User, error) {
	user := new(model.User)
	err := DB().WithContext(dao.ctx).Where(model.User{LoginName: loginName}).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return user, nil
}

func (dao *UserDao) CreateUser(info *do.UserBaseInfo, userPasswordHash string) (*model.User, error) {
	userModel := new(model.User)
	err := utils.CopyStruct(userModel, info)
	if err != nil {
		err = errcode.Wrap("UserDaoCreateUserError", err)
		return nil, err
	}
	userModel.Password = userPasswordHash
	err = DBMaster().WithContext(dao.ctx).Create(userModel).Error
	if err != nil {
		err = errcode.Wrap("UserDaoCreateUserError", err)
		return nil, err
	}
	return userModel, nil

}

func (dao *UserDao) UpdateUser(user *model.User) error {
	err := DBMaster().WithContext(dao.ctx).Model(user).Updates(user).Error
	return err
}
