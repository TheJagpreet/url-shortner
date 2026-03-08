// Redis CLI CRUD operations examples:
// List all keys:
//	redis-cli KEYS '*'
// Get value for a key:
//	redis-cli GET <code>
// Add or update a key:
//	redis-cli SET <code> <url>
// Delete a key:
//	redis-cli DEL <code>
// Check TTL for a key:
//	redis-cli TTL <code>
// Example: redis-cli GET abc123
// These commands let you inspect and manage URL shortener keys directly in Redis.
package storage

import (
	"context"
	"time"
	"url-shortner/shortener"

	"github.com/redis/go-redis/v9"
)

// RedisStorage implements Storage interface using Redis

type RedisStorage struct {
	client     *redis.Client
	ttlSeconds int64
}

var _ Storage = (*RedisStorage)(nil) // Ensure RedisStorage implements Storage

const defaultTTLSeconds int64 = 86400 // 24 hours

func NewRedisStorage(addr string, ttlSeconds ...int64) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	ttl := defaultTTLSeconds
	if len(ttlSeconds) > 0 && ttlSeconds[0] > 0 {
		ttl = ttlSeconds[0]
	}
	return &RedisStorage{client: client, ttlSeconds: ttl}
}

func (s *RedisStorage) Shorten(url string, ttlSeconds int64) string {
	code := shortener.GenerateCode()
	ctx := context.Background()
	for {
		exists, _ := s.client.Exists(ctx, code).Result()
		if exists == 0 {
			break
		}
		code = shortener.GenerateCode()
	}
	ttl := ttlSeconds
	if ttl <= 0 {
		ttl = s.ttlSeconds
	}
	s.client.Set(ctx, code, url, time.Duration(ttl)*time.Second)
	return code
}

func (s *RedisStorage) Resolve(code string) (string, bool) {
	ctx := context.Background()
	url, err := s.client.Get(ctx, code).Result()
	if err == redis.Nil || err != nil {
		return "", false
	}
	return url, true
}

func generateCode() string {
	return shortener.GenerateCode()
}
