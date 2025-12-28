package errno

import "github.com/Boyuan-IT-Club/go-kit/errorx/code"

// comment: 100 000 000 ~ 100 999 999

const (
	SuccessCode      = 0
	UnknownErrorCode = 100000001
	ParamErrorCode   = 100000002
)

func init() {
	code.Register(SuccessCode, "success")
	code.Register(UnknownErrorCode, "unknown error")
	code.Register(ParamErrorCode, "parameter error")
}
