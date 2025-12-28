package assembler

import (
	"time"

	"github.com/Eagle233Fake/omniread/backend/application/dto"
	"github.com/Eagle233Fake/omniread/backend/infra/model"
)

func RegisterReqToUser(req *dto.RegisterReq, hashedPassword string) *model.User {
	var birthdate *time.Time
	if req.Birthdate != "" {
		if t, err := time.Parse("2006-01-02", req.Birthdate); err == nil {
			birthdate = &t
		}
	}
	return &model.User{
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashedPassword,
		Gender:    req.Gender,
		Birthdate: birthdate,
		Avatar:    req.Avatar,
		Bio:       req.Bio,
	}
}

func UserToLoginResp(user *model.User, token string) *dto.LoginResp {
	return &dto.LoginResp{
		Token: token,
		User:  user,
	}
}
