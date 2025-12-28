package errno

import "github.com/Boyuan-IT-Club/go-kit/errorx/code"

// user: 102 000 000 ~ 102 999 999

const (
	ErrUserNotFound     = 102000001
	ErrUserExist        = 102000002
	ErrUserCreateFailed = 102000003
)

func init() {
	code.Register(
		ErrUserNotFound,
		"user not found",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrUserExist,
		"user already exists",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrUserCreateFailed,
		"failed to create user",
		code.WithAffectStability(false),
	)
}
