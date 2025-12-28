package router

import (
	"github.com/Eagle233Fake/omniread/backend/api/handler"
	"github.com/Eagle233Fake/omniread/backend/application/service/auth/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", handler.Register)
		authGroup.POST("/login", middleware.LoginRateLimitMiddleware(), handler.Login)
	}

	return r
}
