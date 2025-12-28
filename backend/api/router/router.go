package router

import (
	"github.com/Eagle233Fake/omniread/backend/api/handler"
	"github.com/Eagle233Fake/omniread/backend/api/handler/agent"
	"github.com/Eagle233Fake/omniread/backend/api/handler/book"
	"github.com/Eagle233Fake/omniread/backend/api/handler/reading"
	"github.com/Eagle233Fake/omniread/backend/application/service/auth/middleware"
	"github.com/Eagle233Fake/omniread/backend/infra/config"
	"github.com/Eagle233Fake/omniread/backend/internal/agent/domain"
	"github.com/Eagle233Fake/omniread/backend/internal/agent/service"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/stores/redis"
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
	// NOTE: Middleware should be applied here for protected routes.
	// Assuming it's already there or handled globally for now based on previous context

	userGroup := api.Group("/user")
	{
		userGroup.GET("/profile", handler.GetProfile)
		userGroup.PUT("/profile", handler.UpdateProfile)
		userGroup.PUT("/password", handler.ChangePassword)
		userGroup.PUT("/preferences", handler.UpdatePreferences)
	}

	bookGroup := api.Group("/books")
	{
		bookGroup.POST("/upload", book.UploadBook)
		bookGroup.GET("", book.ListBooks)
		bookGroup.GET("/:id", book.GetBook)
		bookGroup.PUT("/:id", book.UpdateBook)
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

	// Dependency Injection (Stubbed for now)
	cfg := config.GetConfig()
	if cfg == nil {
		// Fallback for testing without config file
		cfg = &config.Config{}
		cfg.Redis = &redis.RedisConf{Host: "localhost:6379", Pass: ""}
	}

	agentRepo := &domain.StubAgentRepository{}
	agentSvc := service.NewAgentService(agentRepo, cfg)
	agentHdl := agent.NewAgentHandler(agentSvc)

	apiV1 := r.Group("/api/v1")
	{
		agents := apiV1.Group("/agents")
		{
			agents.POST("", agentHdl.Create)
			agents.POST("/chat", agentHdl.Chat)
		}
	}

	return r
}
