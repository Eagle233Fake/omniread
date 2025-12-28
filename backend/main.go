package main

import (
	"github.com/Eagle233Fake/omniread/backend/api/router"
	"github.com/Eagle233Fake/omniread/backend/infra/config"
	"github.com/Eagle233Fake/omniread/backend/provider"
	_ "github.com/Eagle233Fake/omniread/backend/types/errno" // Register error codes

	"github.com/Boyuan-IT-Club/go-kit/logs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Initialize config first
	if _, err := config.NewConfig(); err != nil {
		panic(err)
	}

	provider.Init()
	r := router.SetupRoutes()
	setLogLevel()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func setLogLevel() {
	level := config.GetConfig().Log.Level

	logs.Infof("log level: %s", level)
	switch level {
	case "trace":
		logs.SetLevel(logs.LevelTrace)
	case "debug":
		logs.SetLevel(logs.LevelDebug)
	case "info":
		logs.SetLevel(logs.LevelInfo)
	case "notice":
		logs.SetLevel(logs.LevelNotice)
	case "warn":
		logs.SetLevel(logs.LevelWarn)
	case "error":
		logs.SetLevel(logs.LevelError)
	case "fatal":
		logs.SetLevel(logs.LevelFatal)
	default:
		logs.SetLevel(logs.LevelInfo)
	}
}
