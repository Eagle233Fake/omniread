package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type AuthCache struct {
	rds *redis.Redis
}

func NewAuthCache(rds *redis.Redis) *AuthCache {
	return &AuthCache{
		rds: rds,
	}
}

func (c *AuthCache) SetSession(ctx context.Context, token string, uid string, expiration time.Duration) error {
	key := fmt.Sprintf("auth:token:%s", token)
	return c.rds.SetexCtx(ctx, key, uid, int(expiration.Seconds()))
}

func (c *AuthCache) GetSession(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("auth:token:%s", token)
	return c.rds.GetCtx(ctx, key)
}

func (c *AuthCache) DelSession(ctx context.Context, token string) error {
	key := fmt.Sprintf("auth:token:%s", token)
	_, err := c.rds.DelCtx(ctx, key)
	return err
}

// StoreUserStatus stores user active status
func (c *AuthCache) SetUserStatus(ctx context.Context, uid string, status string) error {
	key := fmt.Sprintf("auth:user:%s:status", uid)
	// Status cache lasts longer, e.g., 24 hours
	return c.rds.SetexCtx(ctx, key, status, 24*60*60)
}
