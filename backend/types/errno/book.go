package errno

import "github.com/Boyuan-IT-Club/go-kit/errorx/code"

// book: 103 000 000 ~ 103 999 999

const (
	ErrBookNotFound = 103000001
	ErrBookExist    = 103000002
)

func init() {
	code.Register(
		ErrBookNotFound,
		"book not found",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrBookExist,
		"book already exists",
		code.WithAffectStability(false),
	)
}
