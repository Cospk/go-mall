package errcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"runtime"
)

type AppError struct {
	code     int    `json:"code"`
	msg      string `json:"msg"`
	cause    error  `json:"cause"`
	occurred string `json:"occurred"`
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	formatedErr := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Cause    string `json:"cause"`
		Occurred string `json:"occurred"`
	}{
		Code:     e.code,
		Msg:      e.msg,
		Occurred: e.occurred,
	}
	if e.cause != nil {
		formatedErr.Cause = e.cause.Error()
	}
	errByte, _ := json.Marshal(formatedErr)
	return string(errByte)
}

func (e *AppError) String() string {
	return e.Error()
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Msg() string {
	return e.msg
}

func (e *AppError) UnWrap() error {
	return e.cause
}

// Is 与上面的UnWrap一起让 *AppError 支持 errors.Is(err, target)
func (e *AppError) Is(target error) bool {
	var targetErr *AppError

	ok := errors.As(target, &targetErr)
	if !ok {
		return false
	}
	return targetErr.Code() == e.Code()
}

// NewError 创建AppError实例
func NewError(code int, msg string) *AppError {
	if code > -1 {
		if _, duplicate := codes[code]; duplicate {
			panic(fmt.Sprintf("code %d already exist", code))
		}
		codes[code] = struct{}{}
	}
	return &AppError{
		code: code,
		msg:  msg,
	}
}

// WithCause和Wrap都是记录错误发生的位置，一个是直接附加错误原因，一个是包装error为AppError

// WithCause 在原有的AppError实例附加错误原因，并记录错误发生的位置
func (e *AppError) WithCause(err error) *AppError {
	newErr := e.Clone()
	newErr.cause = err
	newErr.occurred = getAppErrOccurredInfo()
	return newErr
}

// Clone 克隆AppError,保留之前的错误信息
func (e *AppError) Clone() *AppError {
	return &AppError{
		code:     e.code,
		msg:      e.msg,
		cause:    e.cause,
		occurred: e.occurred,
	}
}

// Wrap 包装error为AppError
func Wrap(msg string, err error) *AppError {
	if err == nil {
		return nil
	}
	appErr := &AppError{
		code:     -1,
		msg:      msg,
		cause:    err,
		occurred: getAppErrOccurredInfo(),
	}
	return appErr
}

// getAppErrOccurredInfo 获取项目中调用Wrap或者WithCause方法时的程序位置, 方便排查问题
func getAppErrOccurredInfo() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	file = path.Base(file)
	funcName := runtime.FuncForPC(pc).Name()
	triggerInfo := fmt.Sprintf("func: %s, file: %s, line: %d", funcName, file, line)
	return triggerInfo
}
