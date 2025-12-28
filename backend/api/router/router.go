package router

import (
	"github.com/Eagle233Fake/omniread/backend/api/handler"
	"github.com/Eagle233Fake/omniread/backend/api/handler/book"
	"github.com/Eagle233Fake/omniread/backend/api/handler/reading"
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

	// Serve static files for uploaded books
	r.Static("/uploads", "./uploads")

	// Protected routes
	// TODO: Add Auth Middleware here
	// For now, assuming middleware sets "uid"
	api := r.Group("/")
	// api.Use(middleware.Auth())

	bookGroup := api.Group("/books")
	{
		bookGroup.POST("/upload", book.UploadBook)
		bookGroup.GET("", book.ListBooks)
		bookGroup.GET("/:id", book.GetBook)
	}

	readingGroup := api.Group("/reading")
	{
		readingGroup.POST("/progress", reading.UpdateProgress)
		readingGroup.GET("/progress", reading.GetProgress)
		readingGroup.POST("/session", reading.RecordSession)
	}

	insightGroup := api.Group("/insight")
	{
		insightGroup.GET("/summary", reading.GetInsightSummary)
	}

	return r
}
