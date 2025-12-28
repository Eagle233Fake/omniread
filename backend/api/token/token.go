package token

import (
	"time"

	"github.com/Eagle233Fake/omniread/backend/infra/config"
	"github.com/Eagle233Fake/omniread/backend/infra/model"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user *model.User) (string, error) {
	cfg := config.GetConfig()
	claims := jwt.MapClaims{
		"uid":      user.ID.Hex(),
		"username": user.Username,
		"exp":      time.Now().Add(time.Duration(cfg.Auth.AccessExpire) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Auth.SecretKey))
}
