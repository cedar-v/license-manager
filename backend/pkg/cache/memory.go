package cache

import (
	"context"
	"sync"
	"time"
)

// memoryCache 内存缓存实现
type memoryCache struct {
	data   sync.Map
	ttlMap sync.Map
	maxSize int
	stats  CacheStats
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(maxSize int) Cache {
	cache := &memoryCache{
		data:    sync.Map{},
		ttlMap:  sync.Map{},
		maxSize: maxSize,
	}

	// 启动清理协程
	go cache.cleanupRoutine()

	return cache
}

// cleanupRoutine 定期清理过期数据
func (m *memoryCache) cleanupRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.cleanupExpired()
	}
}

// cleanupExpired 清理过期数据
func (m *memoryCache) cleanupExpired() {
	now := time.Now()
	m.ttlMap.Range(func(key, value interface{}) bool {
		if expireTime, ok := value.(time.Time); ok && now.After(expireTime) {
			m.data.Delete(key)
			m.ttlMap.Delete(key)
		}
		return true
	})
}

func (m *memoryCache) Get(ctx context.Context, key string) (string, error) {
	defer func(start time.Time) {
		m.stats.AvgRespTime = (m.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	value, ok := m.data.Load(key)
	if !ok {
		m.stats.RecordMiss()
		return "", ErrCacheMiss
	}

	// 检查是否过期
	if ttl, exists := m.ttlMap.Load(key); exists {
		if expireTime, ok := ttl.(time.Time); ok && time.Now().After(expireTime) {
			m.data.Delete(key)
			m.ttlMap.Delete(key)
			m.stats.RecordMiss()
			return "", ErrCacheMiss
		}
	}

	m.stats.RecordHit()
	if str, ok := value.(string); ok {
		return str, nil
	}

	m.stats.RecordError()
	return "", NewCacheError("invalid cache value type")
}

func (m *memoryCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	defer func(start time.Time) {
		m.stats.AvgRespTime = (m.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	// 简单的容量控制（生产环境建议使用LRU库）
	if m.maxSize > 0 {
		// 简单的清理策略：当存储量接近最大值时，清理一些过期数据
		size := 0
		m.data.Range(func(key, value interface{}) bool {
			size++
			return size < m.maxSize
		})

		if size >= m.maxSize {
			m.cleanupExpired()
		}
	}

	m.data.Store(key, value)
	if ttl > 0 {
		m.ttlMap.Store(key, time.Now().Add(ttl))
	} else {
		m.ttlMap.Delete(key) // 永不过期
	}

	return nil
}

func (m *memoryCache) Del(ctx context.Context, keys ...string) error {
	defer func(start time.Time) {
		m.stats.AvgRespTime = (m.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	for _, key := range keys {
		m.data.Delete(key)
		m.ttlMap.Delete(key)
	}
	return nil
}

func (m *memoryCache) Exists(ctx context.Context, key string) (bool, error) {
	defer func(start time.Time) {
		m.stats.AvgRespTime = (m.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	_, exists := m.data.Load(key)
	if !exists {
		return false, nil
	}

	// 检查是否过期
	if ttl, exists := m.ttlMap.Load(key); exists {
		if expireTime, ok := ttl.(time.Time); ok && time.Now().After(expireTime) {
			m.data.Delete(key)
			m.ttlMap.Delete(key)
			return false, nil
		}
	}

	return true, nil
}

func (m *memoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	defer func(start time.Time) {
		m.stats.AvgRespTime = (m.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	ttl, exists := m.ttlMap.Load(key)
	if !exists {
		return -1, nil // 没有设置过期时间
	}

	if expireTime, ok := ttl.(time.Time); ok {
		remaining := time.Until(expireTime)
		if remaining <= 0 {
			m.data.Delete(key)
			m.ttlMap.Delete(key)
			return -2, nil // 已过期
		}
		return remaining, nil
	}

	return -1, nil
}

func (m *memoryCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	defer func(start time.Time) {
		m.stats.AvgRespTime = (m.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	if ttl <= 0 {
		m.ttlMap.Delete(key)
		return nil
	}

	// 检查key是否存在
	if _, exists := m.data.Load(key); !exists {
		return NewCacheError("key not found")
	}

	m.ttlMap.Store(key, time.Now().Add(ttl))
	return nil
}

func (m *memoryCache) Incr(ctx context.Context, key string) (int64, error) {
	defer func(start time.Time) {
		m.stats.AvgRespTime = (m.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	// 内存缓存的数值操作实现比较简单
	return 0, NewCacheError("Incr not implemented for memory cache")
}

func (m *memoryCache) Decr(ctx context.Context, key string) (int64, error) {
	defer func(start time.Time) {
		m.stats.AvgRespTime = (m.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	return 0, NewCacheError("Decr not implemented for memory cache")
}

func (m *memoryCache) MGet(ctx context.Context, keys ...string) ([]string, error) {
	defer func(start time.Time) {
		m.stats.AvgRespTime = (m.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	results := make([]string, len(keys))
	for i, key := range keys {
		if value, err := m.Get(ctx, key); err == nil {
			results[i] = value
		} else {
			results[i] = ""
		}
	}
	return results, nil
}

func (m *memoryCache) MSet(ctx context.Context, pairs map[string]string, ttl time.Duration) error {
	defer func(start time.Time) {
		m.stats.AvgRespTime = (m.stats.AvgRespTime + time.Since(start).Nanoseconds()) / 2
	}(time.Now())

	for key, value := range pairs {
		if err := m.Set(ctx, key, value, ttl); err != nil {
			return err
		}
	}
	return nil
}

func (m *memoryCache) Ping(ctx context.Context) error {
	return nil
}

func (m *memoryCache) Close() error {
	return nil
}

// GetStats 获取缓存统计信息
func (m *memoryCache) GetStats() CacheStats {
	return m.stats
}
