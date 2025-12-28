package provider

import (
	"context"
	"time"

	"github.com/Boyuan-IT-Club/go-kit/logs"
	"github.com/Eagle233Fake/omniread/backend/application/service/auth"
	"github.com/Eagle233Fake/omniread/backend/application/service/book"
	"github.com/Eagle233Fake/omniread/backend/application/service/insight"
	"github.com/Eagle233Fake/omniread/backend/application/service/reading"
	"github.com/Eagle233Fake/omniread/backend/infra/cache"
	"github.com/Eagle233Fake/omniread/backend/infra/config"
	"github.com/Eagle233Fake/omniread/backend/infra/oss"
	"github.com/Eagle233Fake/omniread/backend/infra/repo"
	agentservice "github.com/Eagle233Fake/omniread/backend/internal/agent/service"
	"github.com/google/wire"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var provider *Provider

func Init() {
	var err error
	provider, err = NewProvider()
	if err != nil {
		panic(err)
	}
}

func Get() *Provider {
	return provider
}

// Provider 提供Handler依赖的对象
type Provider struct {
	Config         *config.Config
	AuthService    *auth.AuthService
	BookService    *book.BookService
	ReadingService *reading.ReadingService
	InsightService *insight.InsightService
	AgentService   *agentservice.AgentService
}

var ApplicationSet = wire.NewSet(
	auth.AuthServiceSet,
	book.BookServiceSet,
	reading.ReadingServiceSet,
	insight.InsightServiceSet,
	agentservice.NewAgentService,
)

var InfraSet = wire.NewSet(
	config.NewConfig,
	repo.UserRepoSet,
	repo.ReadingRepoSet,
	repo.AgentRepoSet,
	GetDB,
	GetRedis,
	cache.NewAuthCache,
	oss.NewOSSClient,
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	InfraSet,
)

func GetDB(cfg *config.Config) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.Mongo.URL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logs.Errorf("Failed to connect to MongoDB: %v", err)
		panic(err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		logs.Errorf("Failed to ping MongoDB: %v", err)
		panic(err)
	}

	logs.Infof("Connected to MongoDB: %s", cfg.Mongo.DB)
	return client.Database(cfg.Mongo.DB)
}

func GetRedis(cfg *config.Config) *redis.Redis {
	if cfg.Redis == nil {
		return nil
	}
	r := redis.MustNewRedis(*cfg.Redis)
	return r
}
