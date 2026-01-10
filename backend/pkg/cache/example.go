// Package cache 缓存模块使用示例
package cache

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ExampleUsage 缓存使用示例
func ExampleUsage() {
	ctx := context.Background()

	// 1. 初始化内存缓存
	memCache, _ := NewCache(CacheConfig{
		Type:    "memory",
		TTL:     time.Hour,
		Enabled: true,
		Memory: MemoryConfig{
			MaxSize: 1000,
		},
	})

	// 2. 基本操作
	memCache.Set(ctx, "user:123", "John Doe", time.Hour)
	user, _ := memCache.Get(ctx, "user:123")
	fmt.Printf("User: %s\n", user)

	// 3. 批量操作
	memCache.MSet(ctx, map[string]string{
		"config:app":  "v1.0.0",
		"config:env":  "production",
		"config:host": "api.example.com",
	}, time.Hour)

	configs, _ := memCache.MGet(ctx, "config:app", "config:env", "config:host")
	fmt.Printf("Configs: %v\n", configs)

	// 4. 使用键构建器
	keyBuilder := NewKeyBuilder("myapp")
	userKey := keyBuilder.User("123")
	sessionKey := keyBuilder.Session("abc123")
	smsKey := keyBuilder.SMSCode("+8613800012345")

	memCache.Set(ctx, userKey, "user_data", time.Hour)
	memCache.Set(ctx, sessionKey, "session_data", 30*time.Minute)
	memCache.Set(ctx, smsKey, "123456", 5*time.Minute)

	// 5. 使用序列化包装器
	cached := &Cached{
		Cache:      memCache,
		Serializer: &JSONSerializer{},
	}

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	userObj := User{ID: 1, Name: "Alice", Age: 30}
	cached.SetObject(ctx, "user:object:1", userObj, time.Hour)

	var retrievedUser User
	cached.GetObject(ctx, "user:object:1", &retrievedUser)
	fmt.Printf("Retrieved user: %+v\n", retrievedUser)

	// 6. 数值操作
	count, _ := memCache.Incr(ctx, "counter:visits")
	fmt.Printf("Visit count: %d\n", count)

	// 7. TTL操作
	memCache.Set(ctx, "temp:key", "temp_value", 10*time.Second)
	ttl, _ := memCache.TTL(ctx, "temp:key")
	fmt.Printf("TTL: %v\n", ttl)

	// 8. 统计信息
	if memCacheWithStats, ok := memCache.(interface{ GetStats() CacheStats }); ok {
		stats := memCacheWithStats.GetStats()
		fmt.Printf("Cache stats - Hits: %d, Misses: %d, HitRate: %.2f%%\n",
			stats.Hits, stats.Misses, stats.HitRate*100)
	}

	// 9. 优雅关闭
	memCache.Close()
}

// ExampleRedisUsage Redis缓存使用示例
func ExampleRedisUsage() {
	// 初始化Redis缓存
	redisCache, err := NewCache(CacheConfig{
		Type:    "redis",
		TTL:     time.Hour,
		Enabled: true,
		Redis: RedisConfig{
			Host:            "localhost",
			Port:            6379,
			Password:        "",
			DB:              0,
			PoolSize:        10,
			MinIdleConns:    5,
			ConnMaxIdleTime: 30 * time.Minute,
			DialTimeout:     5 * time.Second,
			ReadTimeout:     3 * time.Second,
			WriteTimeout:    3 * time.Second,
		},
	})

	if err != nil {
		log.Printf("Failed to create Redis cache: %v", err)
		return
	}

	ctx := context.Background()

	// 基础操作
	err = redisCache.Set(ctx, "user:redis:123", "Redis User", time.Hour)
	if err != nil {
		log.Printf("Redis set error: %v", err)
		return
	}

	user, err := redisCache.Get(ctx, "user:redis:123")
	if err != nil {
		log.Printf("Redis get error: %v", err)
	} else {
		fmt.Printf("Redis user: %s\n", user)
	}

	// Redis特有的批量操作会更高效
	data := make(map[string]string)
	for i := 0; i < 10; i++ {
		data[fmt.Sprintf("batch:key:%d", i)] = fmt.Sprintf("redis_value_%d", i)
	}

	redisCache.MSet(ctx, data, time.Hour)

	// 数值操作在Redis中非常高效
	counter, _ := redisCache.Incr(ctx, "redis:global:counter")
	fmt.Printf("Redis global counter: %d\n", counter)

	// TTL查询
	ttl, _ := redisCache.TTL(ctx, "user:redis:123")
	fmt.Printf("TTL remaining: %v\n", ttl)

	redisCache.Close()
}

// ExampleSafeCacheOperation 优雅降级的缓存操作示例
func ExampleSafeCacheOperation(cache Cache) error {
	return safeCacheOperation(cache, func() error {
		ctx := context.Background()

		// 正常的缓存操作
		cache.Set(ctx, "safe:key", "safe_value", time.Hour)
		value, err := cache.Get(ctx, "safe:key")
		if err != nil {
			return err
		}

		fmt.Printf("Safe cache value: %s\n", value)
		return nil
	})
}
