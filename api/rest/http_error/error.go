package http_error

import (
	"github.com/huichen/kunlun/api/rest"
)

const (
	ErrGenericServerError   = -1
	ErrNoAccessPermission   = 1
	ErrLoginUsernameIsEmpty = 2
	ErrMissingParameter     = 3
	ErrIllegalParameter     = 4
	ErrDuplicateRecord      = 5
)

// 这里包含一些预定义的错误码和提示信息
var (
	GenericServerError = rest.ResponseBase{
		Code:    ErrGenericServerError,
		Message: "服务返回错误",
	}

	NoAccessPermission = rest.ResponseBase{
		Code:    ErrNoAccessPermission,
		Message: "没有访问权限",
	}

	LoginUsernameIsEmpty = rest.ResponseBase{
		Code:    ErrLoginUsernameIsEmpty,
		Message: "用户名不能为空",
	}

	MissingParameter = rest.ResponseBase{
		Code:    ErrMissingParameter,
		Message: "字段填写不完整",
	}

	IllegalParameter = rest.ResponseBase{
		Code:    ErrIllegalParameter,
		Message: "字段非法",
	}

	DuplicateRecord = rest.ResponseBase{
		Code:    ErrDuplicateRecord,
		Message: "重复数据",
	}
)

func GetError(err error) rest.ResponseBase {
	return rest.ResponseBase{
		Code:    ErrGenericServerError,
		Message: err.Error(),
	}
}

func GetErrorFromString(errMessage string) rest.ResponseBase {
	return rest.ResponseBase{
		Code:    ErrGenericServerError,
		Message: errMessage,
	}
}
