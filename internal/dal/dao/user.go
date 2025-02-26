package dao

import (
	"context"
	"github.com/Cospk/go-mall/internal/dal/model"
)

type UserDao struct {
	ctx context.Context
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{ctx: ctx}
}

func (dao *UserDao) FindUserById(id int64) *model.User {
	return &model.User{}
}

func (dao *UserDao) FindUserByName(name string) (model.User, error) {
	// TODO 执行sql查询数据库的数据
	return model.User{}, nil
}
