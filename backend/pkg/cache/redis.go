package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig Redis配置
type RedisConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Password        string        `mapstructure:"password"`
	DB              int           `mapstructure:"db"`
	PoolSize        int           `mapstructure:"pool_size"`
	MinIdleConns    int           `mapstructure:"min_idle_conns"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	DialTimeout     time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
}

// redisCache Redis缓存实现
type redisCache struct {
	client *redis.Client
	ttl    time.Duration
	stats  CacheStats
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(config RedisConfig, defaultTTL time.Duration) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:        config.Password,
		DB:              config.DB,
		PoolSize:        config.PoolSize,
		MinIdleConns:    config.MinIdleConns,
		ConnMaxIdleTime: config.ConnMaxIdleTime,
		DialTimeout:     config.DialTimeout,
		ReadTimeout:     config.ReadTimeout,
		WriteTimeout:    config.WriteTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &redisCache{
		client: client,
		ttl:    defaultTTL,
	}, nil
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	defer func(start time.Time) {
		r.stats.AvgRespTime = (r.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	value, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		r.stats.RecordMiss()
		return "", ErrCacheMiss
	}
	if err != nil {
		r.stats.RecordError()
		return "", NewCacheError(fmt.Sprintf("redis get error: %v", err))
	}

	r.stats.RecordHit()
	return value, nil
}

func (r *redisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	defer func(start time.Time) {
		r.stats.AvgRespTime = (r.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	if ttl == 0 {
		ttl = r.ttl
	}

	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		r.stats.RecordError()
		return NewCacheError(fmt.Sprintf("redis set error: %v", err))
	}

	return nil
}

func (r *redisCache) Del(ctx context.Context, keys ...string) error {
	defer func(start time.Time) {
		r.stats.AvgRespTime = (r.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		r.stats.RecordError()
		return NewCacheError(fmt.Sprintf("redis del error: %v", err))
	}

	return nil
}

func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	defer func(start time.Time) {
		r.stats.AvgRespTime = (r.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		r.stats.RecordError()
		return false, NewCacheError(fmt.Sprintf("redis exists error: %v", err))
	}

	return count > 0, nil
}

func (r *redisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	defer func(start time.Time) {
		r.stats.AvgRespTime = (r.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		r.stats.RecordError()
		return 0, NewCacheError(fmt.Sprintf("redis ttl error: %v", err))
	}

	return ttl, nil
}

func (r *redisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	defer func(start time.Time) {
		r.stats.AvgRespTime = (r.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	success, err := r.client.Expire(ctx, key, ttl).Result()
	if err != nil {
		r.stats.RecordError()
		return NewCacheError(fmt.Sprintf("redis expire error: %v", err))
	}

	if !success {
		return NewCacheError("key not found")
	}

	return nil
}

func (r *redisCache) Incr(ctx context.Context, key string) (int64, error) {
	defer func(start time.Time) {
		r.stats.AvgRespTime = (r.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	value, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		r.stats.RecordError()
		return 0, NewCacheError(fmt.Sprintf("redis incr error: %v", err))
	}

	return value, nil
}

func (r *redisCache) Decr(ctx context.Context, key string) (int64, error) {
	defer func(start time.Time) {
		r.stats.AvgRespTime = (r.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	value, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		r.stats.RecordError()
		return 0, NewCacheError(fmt.Sprintf("redis decr error: %v", err))
	}

	return value, nil
}

func (r *redisCache) MGet(ctx context.Context, keys ...string) ([]string, error) {
	defer func(start time.Time) {
		r.stats.AvgRespTime = (r.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	values, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		r.stats.RecordError()
		return nil, NewCacheError(fmt.Sprintf("redis mget error: %v", err))
	}

	results := make([]string, len(values))
	for i, v := range values {
		if v != nil {
			results[i] = fmt.Sprintf("%v", v)
		} else {
			results[i] = ""
		}
	}

	return results, nil
}

func (r *redisCache) MSet(ctx context.Context, pairs map[string]string, ttl time.Duration) error {
	defer func(start time.Time) {
		r.stats.AvgRespTime = (r.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	if ttl == 0 {
		ttl = r.ttl
	}

	// Redis MSet 不支持每个key单独设置TTL，所以这里使用管道
	pipe := r.client.Pipeline()
	for key, value := range pairs {
		pipe.Set(ctx, key, value, ttl)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		r.stats.RecordError()
		return NewCacheError(fmt.Sprintf("redis mset error: %v", err))
	}

	return nil
}

func (r *redisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *redisCache) Close() error {
	return r.client.Close()
}
