package config

import (
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

var config *Config

type Auth struct {
	SecretKey    string
	PublicKey    string
	AccessExpire int64
}

type WeApp struct {
	AppID     string
	AppSecret string
}

type OSS struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
}

type Bocha struct {
	APIKey string
}

type Model struct {
	BaseURL string
	APIKey  string
	Model   string
}

type Config struct {
	service.ServiceConf
	ListenOn string
	State    string
	Auth     Auth
	Mongo    struct {
		URL string
		DB  string
	}
	Cache cache.CacheConf
	Redis *redis.RedisConf
	WeApp WeApp
	OSS   OSS
	Bocha Bocha
	Model Model
}

func NewConfig() (*Config, error) {
	c := new(Config)
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "etc/config.yaml"
	}
	err := conf.Load(path, c)
	if err != nil {
		return nil, err
	}
	err = c.SetUp()
	if err != nil {
		return nil, err
	}
	config = c
	return c, nil
}

func GetConfig() *Config {
	return config
}
