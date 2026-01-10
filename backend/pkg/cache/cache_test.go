package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache(100)

	ctx := context.Background()

	// Test Set and Get
	err := cache.Set(ctx, "test:key", "test_value", time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	value, err := cache.Get(ctx, "test:key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}

	// Test cache miss
	_, err = cache.Get(ctx, "nonexistent:key")
	if err != ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss, got %v", err)
	}

	// Test Exists
	exists, err := cache.Exists(ctx, "test:key")
	if err != nil || !exists {
		t.Errorf("Exists failed or returned false")
	}

	exists, err = cache.Exists(ctx, "nonexistent:key")
	if err != nil || exists {
		t.Errorf("Exists should return false for nonexistent key")
	}

	// Test TTL
	ttl, err := cache.TTL(ctx, "test:key")
	if err != nil {
		t.Fatalf("TTL failed: %v", err)
	}

	if ttl <= 0 {
		t.Errorf("TTL should be positive")
	}

	// Test Del
	err = cache.Del(ctx, "test:key")
	if err != nil {
		t.Fatalf("Del failed: %v", err)
	}

	// Verify deletion
	_, err = cache.Get(ctx, "test:key")
	if err != ErrCacheMiss {
		t.Errorf("Key should be deleted")
	}
}

func TestKeyBuilder(t *testing.T) {
	kb := NewKeyBuilder("myapp")

	userKey := kb.User("123")
	expected := "myapp:user:123"
	if userKey != expected {
		t.Errorf("Expected '%s', got '%s'", expected, userKey)
	}

	sessionKey := kb.Session("abc")
	expected = "myapp:session:abc"
	if sessionKey != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sessionKey)
	}

	smsKey := kb.SMSCode("13800138000")
	expected = "myapp:sms:code:13800138000"
	if smsKey != expected {
		t.Errorf("Expected '%s', got '%s'", expected, smsKey)
	}
}

func TestCachedWrapper(t *testing.T) {
	cache := NewMemoryCache(100)
	cached := &Cached{
		Cache:      cache,
		Serializer: &JSONSerializer{},
	}

	ctx := context.Background()

	type TestUser struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	user := TestUser{ID: 1, Name: "Alice"}

	// Test SetObject
	err := cached.SetObject(ctx, "user:1", user, time.Minute)
	if err != nil {
		t.Fatalf("SetObject failed: %v", err)
	}

	// Test GetObject
	var retrieved TestUser
	err = cached.GetObject(ctx, "user:1", &retrieved)
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}

	if retrieved.ID != user.ID || retrieved.Name != user.Name {
		t.Errorf("Retrieved object doesn't match original")
	}
}

func TestCacheFactory(t *testing.T) {
	// Test memory cache
	config := CacheConfig{
		Type:    "memory",
		TTL:     time.Hour,
		Enabled: true,
		Memory: MemoryConfig{
			MaxSize: 100,
		},
	}

	cache, err := NewCache(config)
	if err != nil {
		t.Fatalf("NewCache failed: %v", err)
	}

	if cache == nil {
		t.Fatal("Cache should not be nil")
	}

	// Test disabled cache
	config.Enabled = false
	cache, err = NewCache(config)
	if err != nil {
		t.Fatalf("NewCache with disabled should not fail: %v", err)
	}

	// Should be noOpCache
	ctx := context.Background()
	err = cache.Set(ctx, "test", "value", time.Minute)
	if err != nil {
		t.Errorf("noOpCache Set should not fail: %v", err)
	}

	_, err = cache.Get(ctx, "test")
	if err != ErrCacheMiss {
		t.Errorf("noOpCache Get should return ErrCacheMiss")
	}
}

func TestRedisCacheUnavailable(t *testing.T) {
	config := CacheConfig{
		Type:    "redis",
		TTL:     time.Hour,
		Enabled: true,
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
		},
	}

	_, err := NewCache(config)
	if err == nil {
		t.Error("Redis cache should not be available due to dependency issues")
	}
}
