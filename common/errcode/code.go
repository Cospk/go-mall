package errcode

import "net/http"

var codes = map[int]struct{}{}

// 根据模块来定义错误码

// 公共错误码,预留 10000 ~ 10099间的100个错误码
var (
	Success            = NewError(0, "成功")
	ErrServer          = NewError(10000, "服务内部错误")
	ErrParams          = NewError(10001, "入参错误")
	ErrNotFound        = NewError(10002, "找不到")
	ErrPanic           = NewError(10003, "(*^__^*)系统开小差了,请稍后重试")
	ErrToken           = NewError(10004, "鉴权失败，Token错误")
	ErrForbid          = NewError(10005, "禁止访问")
	ErrTooManyRequests = NewError(10006, "请求过多")
	ErrCoverData       = NewError(10007, "数据转换错误")
)

// 用户模块错误码， 预留11000 ~ 11099间的100个错误码
var (
	ErrUserNotFound            = NewError(11000, "用户不存在")
	ErrUserNameOrPasswordError = NewError(11001, "用户名或密码错误")
	ErrUserRegisterFailed      = NewError(11002, "用户注册失败")
	ErrUserLoginFailed         = NewError(11003, "用户登录失败")
	ErrUserUpdateFailed        = NewError(11004, "用户更新失败")
)

// 其他。。。

func (e *AppError) HttpStatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ErrParams.Code():
		return http.StatusBadRequest
	case ErrToken.Code():
		return http.StatusUnauthorized
	case ErrForbid.Code():
		return http.StatusForbidden
	case ErrNotFound.Code():
		return http.StatusNotFound
	case ErrTooManyRequests.Code():
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
