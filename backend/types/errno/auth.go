package errno

import "github.com/Boyuan-IT-Club/go-kit/errorx/code"

// auth: 101 000 000 ~ 101 999 999

const (
	ErrAuthFailed       = 101000001
	ErrTokenInvalid     = 101000002
	ErrTokenExpired     = 101000003
	ErrPermissionDenied = 101000004
)

func init() {
	code.Register(
		ErrAuthFailed,
		"authentication failed",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrTokenInvalid,
		"token is invalid",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrTokenExpired,
		"token is expired",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrPermissionDenied,
		"permission denied",
		code.WithAffectStability(false),
	)
}
