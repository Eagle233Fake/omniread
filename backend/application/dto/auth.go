package dto

import "github.com/Eagle233Fake/omniread/backend/infra/model"

type RegisterReq struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Gender    string `json:"gender"`
	Birthdate string `json:"birthdate"` // YYYY-MM-DD
	Avatar    string `json:"avatar"`
	Bio       string `json:"bio"`
}

type LoginReq struct {
	Username string `json:"username"` // Can be username, email, or phone
	Password string `json:"password"`
}

type LoginResp struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}
