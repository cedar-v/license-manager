# Cache Package

基于Clean Architecture设计的缓存模块，支持内存缓存和Redis缓存两种实现方式。

## 特性

- ✅ 统一的缓存接口抽象
- ✅ 支持内存和Redis缓存
- ✅ 配置驱动，支持运行时切换
- ✅ 优雅降级，缓存异常不影响业务
- ✅ 序列化支持（JSON）
- ✅ 缓存键构建器
- ✅ 性能监控和统计
- ✅ 连接池管理和超时控制

## 快速开始

### 1. 初始化缓存

```go
import "license-manager/pkg/cache"

// 内存缓存
memCache, _ := cache.NewCache(cache.CacheConfig{
    Type:    "memory",
    TTL:     time.Hour,
    Enabled: true,
    Memory: cache.MemoryConfig{
        MaxSize: 10000,
    },
})

// Redis缓存
redisCache, _ := cache.NewCache(cache.CacheConfig{
    Type:    "redis",
    TTL:     time.Hour,
    Enabled: true,
    Redis: cache.RedisConfig{
        Host: "localhost",
        Port: 6379,
        PoolSize: 10,
    },
})
```

### 2. 基本操作

```go
ctx := context.Background()

// 设置缓存
cache.Set(ctx, "key", "value", time.Hour)

// 获取缓存
value, err := cache.Get(ctx, "key")
if err == cache.ErrCacheMiss {
    // 缓存未命中
}

// 删除缓存
cache.Del(ctx, "key1", "key2")

// 检查存在
exists, _ := cache.Exists(ctx, "key")
```

### 3. 高级功能

```go
// TTL操作
ttl, _ := cache.TTL(ctx, "key")
cache.Expire(ctx, "key", time.Hour)

// 数值操作
count, _ := cache.Incr(ctx, "counter")

// 批量操作
cache.MSet(ctx, map[string]string{
    "k1": "v1",
    "k2": "v2",
}, time.Hour)

values, _ := cache.MGet(ctx, "k1", "k2")
```

### 4. 序列化支持

```go
// 使用序列化包装器
cached := &cache.Cached{
    Cache: cache,
    Serializer: &cache.JSONSerializer{},
}

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// 存储对象
user := User{ID: 1, Name: "Alice"}
cached.SetObject(ctx, "user:1", user, time.Hour)

// 获取对象
var retrieved User
cached.GetObject(ctx, "user:1", &retrieved)
```

### 5. 键构建器

```go
keyBuilder := cache.NewKeyBuilder("myapp")

userKey := keyBuilder.User("123")        // "myapp:user:123"
sessionKey := keyBuilder.Session("abc")  // "myapp:session:abc"
smsKey := keyBuilder.SMSCode("13800138000") // "myapp:sms:code:13800138000"
```

## 配置

### 内存缓存配置

```yaml
cache:
  enabled: true
  type: memory
  ttl: 30m
  memory:
    max_size: 10000         # 最大缓存条目数
    cleanup_interval: 10m   # 清理间隔
```

### Redis缓存配置

```yaml
cache:
  enabled: true
  type: redis
  ttl: 30m
  redis:
    host: localhost
    port: 6379
    password: ""
    db: 0
    pool_size: 10
    min_idle_conns: 5
    conn_max_idle_time: 30m
    dial_timeout: 5s
    read_timeout: 3s
    write_timeout: 3s
```

## 接口定义

```go
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
```

## 错误处理

```go
value, err := cache.Get(ctx, "key")
if err != nil {
    if err == cache.ErrCacheMiss {
        // 缓存未命中，执行数据库查询
        return loadFromDatabase()
    }
    // 其他错误，记录日志但不影响业务
    log.Printf("Cache error: %v", err)
}
```

## Redis支持说明

### 当前状态 ✅
Redis缓存支持现已完全可用！支持内存缓存和Redis缓存无缝切换。

### 启用Redis缓存

1. **修改配置**：
   ```yaml
   cache:
     enabled: true
     type: redis  # 切换到redis
     ttl: 30m
     redis:
       host: localhost
       port: 6379
       password: ""
       db: 0
       pool_size: 10
       min_idle_conns: 5
       conn_max_idle_time: 30m
       dial_timeout: 5s
       read_timeout: 3s
       write_timeout: 3s
   ```

2. **重启应用**：
   应用启动时会自动连接Redis并测试连接。

### 切换回内存缓存

```yaml
cache:
  type: memory  # 切换回内存缓存
```

### Redis优势
- **分布式缓存**：支持多实例共享缓存
- **持久化**：数据不会因重启丢失
- **丰富的数据结构**：支持List、Set、Hash等
- **原子操作**：Incr/Decr等操作原子性
- **高性能**：连接池和管道操作
- **监控完善**：详细的性能统计

### 依赖版本说明
- Redis客户端：`github.com/redis/go-redis/v9 v9.17.2`
- 阿里云SMS：`github.com/alibabacloud-go/dysmsapi-20170525/v5 v5.4.0`

## 最佳实践

1. **合理设置TTL**：根据业务场景设置合适的过期时间
2. **使用键前缀**：通过KeyBuilder避免键冲突
3. **优雅降级**：使用`safeCacheOperation`确保业务连续性
4. **监控指标**：定期检查缓存命中率和性能指标
5. **序列化优化**：对复杂对象使用JSON序列化

## 性能考虑

- **内存缓存**：适合单机环境，响应速度快
- **Redis缓存**：适合分布式环境，支持持久化
- **连接池**：Redis使用连接池提高并发性能
- **批量操作**：优先使用MGet/MSet减少网络往返

## 扩展

模块设计支持轻松扩展：

1. **新增缓存类型**：实现Cache接口即可
2. **自定义序列化器**：实现Serializer接口
3. **缓存策略**：添加LRU、LFU等淘汰策略
4. **分布式锁**：基于Redis实现分布式锁
