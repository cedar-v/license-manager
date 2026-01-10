package cache

import (
	"context"
	"time"
)

// Cache 缓存接口
type Cache interface {
	// 基础操作
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)

	// TTL操作
	TTL(ctx context.Context, key string) (time.Duration, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error

	// 数值操作
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)

	// 批量操作
	MGet(ctx context.Context, keys ...string) ([]string, error)
	MSet(ctx context.Context, pairs map[string]string, ttl time.Duration) error

	// 工具方法
	Ping(ctx context.Context) error
	Close() error
}

// CacheWithTTL 带TTL查询的缓存接口
type CacheWithTTL interface {
	Cache
	GetWithTTL(ctx context.Context, key string) (string, time.Duration, error)
}

// Serializer 序列化接口
type Serializer interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// 错误定义
var (
	ErrCacheMiss = NewCacheError("cache miss")
)

// CacheError 缓存错误
type CacheError struct {
	Message string
}

func (e *CacheError) Error() string {
	return e.Message
}

func NewCacheError(message string) *CacheError {
	return &CacheError{Message: message}
}

// CacheStats 缓存统计信息
type CacheStats struct {
	Hits        int64   // 命中次数
	Misses      int64   // 未命中次数
	Errors      int64   // 错误次数
	HitRate     float64 // 命中率
	TotalOps    int64   // 总操作数
	AvgRespTime int64   // 平均响应时间(纳秒)
}

func (s *CacheStats) RecordHit() {
	s.Hits++
	s.updateHitRate()
	s.TotalOps++
}

func (s *CacheStats) RecordMiss() {
	s.Misses++
	s.updateHitRate()
	s.TotalOps++
}

func (s *CacheStats) RecordError() {
	s.Errors++
	s.TotalOps++
}

func (s *CacheStats) updateHitRate() {
	total := s.Hits + s.Misses
	if total > 0 {
		s.HitRate = float64(s.Hits) / float64(total)
	}
}
