package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// CacheConfig 缓存配置
type CacheConfig struct {
	Type    string        `mapstructure:"type"`    // 缓存类型: memory, redis
	TTL     time.Duration `mapstructure:"ttl"`     // 默认TTL
	Enabled bool          `mapstructure:"enabled"` // 是否启用缓存
	Redis   RedisConfig   `mapstructure:"redis"`   // Redis配置
	Memory  MemoryConfig  `mapstructure:"memory"`  // 内存配置
}

// MemoryConfig 内存配置
type MemoryConfig struct {
	MaxSize         int           `mapstructure:"max_size"`         // 最大缓存条目数
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"` // 清理间隔
}

// NewCache 创建缓存实例
func NewCache(config CacheConfig) (Cache, error) {
	if !config.Enabled {
		return &noOpCache{}, nil
	}

	switch config.Type {
	case "memory":
		cache := NewMemoryCache(config.Memory.MaxSize)
		return cache, nil
	case "redis":
		// TODO: Redis支持暂时不可用，需要解决依赖版本冲突
		return nil, fmt.Errorf("Redis cache not available: dependency issue, use memory cache instead")
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", config.Type)
	}
}

// noOpCache 空操作缓存（当缓存被禁用时使用）
type noOpCache struct{}

func (n *noOpCache) Get(ctx context.Context, key string) (string, error) {
	return "", ErrCacheMiss
}

func (n *noOpCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return nil
}

func (n *noOpCache) Del(ctx context.Context, keys ...string) error {
	return nil
}

func (n *noOpCache) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (n *noOpCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return -1, nil
}

func (n *noOpCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return nil
}

func (n *noOpCache) Incr(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (n *noOpCache) Decr(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (n *noOpCache) MGet(ctx context.Context, keys ...string) ([]string, error) {
	results := make([]string, len(keys))
	return results, nil
}

func (n *noOpCache) MSet(ctx context.Context, pairs map[string]string, ttl time.Duration) error {
	return nil
}

func (n *noOpCache) Ping(ctx context.Context) error {
	return nil
}

func (n *noOpCache) Close() error {
	return nil
}

// JSONSerializer JSON序列化器
type JSONSerializer struct{}

func (j *JSONSerializer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JSONSerializer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Cached 带序列化的缓存包装器
type Cached struct {
	Cache
	Serializer Serializer
}

func (c *Cached) SetObject(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := c.Serializer.Marshal(value)
	if err != nil {
		return err
	}
	return c.Cache.Set(ctx, key, string(data), ttl)
}

func (c *Cached) GetObject(ctx context.Context, key string, dest interface{}) error {
	value, err := c.Cache.Get(ctx, key)
	if err != nil {
		return err
	}
	return c.Serializer.Unmarshal([]byte(value), dest)
}

// safeCacheOperation 优雅降级的缓存操作
func safeCacheOperation(cache Cache, operation func() error) error {
	defer func() {
		if r := recover(); r != nil {
			// 记录错误日志，但不影响业务
		}
	}()

	if cache == nil {
		return nil // 缓存未启用，直接跳过
	}

	return operation()
}
