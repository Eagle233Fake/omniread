package middleware

import (
	"strings"
	"time"

	"github.com/Boyuan-IT-Club/go-kit/errorx"
	"github.com/Eagle233Fake/omniread/backend/api/handler"
	"github.com/Eagle233Fake/omniread/backend/infra/config"
	"github.com/Eagle233Fake/omniread/backend/types/errno"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handler.PostProcess(c, nil, nil, errorx.New(errno.ErrAuthFailed))
			c.Abort()
			return
		}

		// Support "Bearer <token>" or just "<token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)

		cfg := config.GetConfig()
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Auth.SecretKey), nil
		})

		if err != nil || !token.Valid {
			handler.PostProcess(c, nil, nil, errorx.New(errno.ErrAuthFailed))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handler.PostProcess(c, nil, nil, errorx.New(errno.ErrAuthFailed))
			c.Abort()
			return
		}

		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				handler.PostProcess(c, nil, nil, errorx.New(errno.ErrAuthFailed))
				c.Abort()
				return
			}
		}

		uid, ok := claims["uid"].(string)
		if !ok {
			handler.PostProcess(c, nil, nil, errorx.New(errno.ErrAuthFailed))
			c.Abort()
			return
		}

		c.Set("uid", uid)
		c.Next()
	}
}
