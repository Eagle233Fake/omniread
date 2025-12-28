package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
)

const (
	defaultHistoryLimit = 20
	historyKeyPrefix    = "agent:history:"
	historyTTL          = 24 * time.Hour
)

type RedisHistoryManager struct {
	client *redis.Client
}

func NewRedisHistoryManager(addr, password string, db int) *RedisHistoryManager {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisHistoryManager{client: rdb}
}

// GetHistory 获取历史消息
func (m *RedisHistoryManager) GetHistory(ctx context.Context, sessionID string) ([]*schema.Message, error) {
	key := historyKeyPrefix + sessionID
	// 使用 LRange 获取列表所有元素
	val, err := m.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("redis lrange error: %w", err)
	}

	messages := make([]*schema.Message, 0, len(val))
	for _, v := range val {
		var msg schema.Message
		if err := json.Unmarshal([]byte(v), &msg); err != nil {
			// 忽略解析失败的消息，避免整个历史中断
			continue
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}

// AddMessage 追加消息
func (m *RedisHistoryManager) AddMessage(ctx context.Context, sessionID string, msg *schema.Message) error {
	key := historyKeyPrefix + sessionID
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message error: %w", err)
	}

	pipe := m.client.Pipeline()
	pipe.RPush(ctx, key, data)
	// 保持固定长度，移除旧消息
	pipe.LTrim(ctx, key, -defaultHistoryLimit, -1)
	pipe.Expire(ctx, key, historyTTL)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis pipeline error: %w", err)
	}

	return nil
}
