package utils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// RedisIDGenerator GORM插件，使用Redis生成唯一ID
type RedisIDGenerator struct {
	client        *redis.Client
	keyPrefix     string
	step          int64
	idCache       map[string]*idSegment
	mutex         sync.Mutex
	retryInterval time.Duration
	maxRetries    int
}

// idSegment 表示一个ID段
type idSegment struct {
	currentID    int64
	remainingIDs int64
	mutex        sync.Mutex
}

// NewRedisIDGenerator 创建一个新的Redis ID生成器插件
func NewRedisIDGenerator(client *redis.Client, keyPrefix string, step int64) *RedisIDGenerator {
	if keyPrefix == "" {
		keyPrefix = "seq:"
	}
	return &RedisIDGenerator{
		client:        client,
		keyPrefix:     keyPrefix,
		step:          step,
		idCache:       make(map[string]*idSegment),
		retryInterval: 100 * time.Millisecond,
		maxRetries:    5,
	}
}

// Name 返回插件名称
func (g *RedisIDGenerator) Name() string {
	return "RedisIDGenerator"
}

// Initialize 初始化插件
func (g *RedisIDGenerator) Initialize(db *gorm.DB) error {
	// 注册回调函数，在创建记录前生成ID
	err := db.Callback().Create().Before("gorm:create").Register("redis_id_generator:before_create", g.beforeCreate)
	if err != nil {
		return err
	}
	return nil
}

// beforeCreate 在创建记录前生成ID
func (g *RedisIDGenerator) beforeCreate(db *gorm.DB) {
	// 只有当主键为零值时才生成ID
	if db.Statement.Schema != nil && db.Statement.Schema.PrioritizedPrimaryField != nil {
		field := db.Statement.Schema.PrioritizedPrimaryField
		if _, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue); isZero {
			tableName := db.Statement.Table
			id, err := g.NextID(tableName)
			if err != nil {
				db.AddError(err)
				return
			}

			// 设置ID值
			if err := field.Set(db.Statement.Context, db.Statement.ReflectValue, id); err != nil {
				db.AddError(err)
			}
		}
	}
}

// NextID 为指定表生成下一个唯一ID
func (g *RedisIDGenerator) NextID(tableName string) (int64, error) {
	g.mutex.Lock()
	segment, exists := g.idCache[tableName]
	if !exists {
		segment = &idSegment{
			currentID:    0,
			remainingIDs: 0,
		}
		g.idCache[tableName] = segment
	}
	g.mutex.Unlock()

	segment.mutex.Lock()
	defer segment.mutex.Unlock()

	// 如果没有剩余ID，从Redis获取新的ID段
	if segment.remainingIDs == 0 {
		if err := g.fetchNewIDSegment(tableName, segment); err != nil {
			return 0, err
		}
	}

	// 使用当前ID并减少剩余ID计数
	id := segment.currentID
	segment.currentID++
	segment.remainingIDs--

	return id, nil
}

// fetchNewIDSegment 从Redis获取新的ID段
func (g *RedisIDGenerator) fetchNewIDSegment(tableName string, segment *idSegment) error {
	ctx := context.Background()
	var err error
	var newID int64
	key := g.keyPrefix + tableName

	// 重试逻辑
	for i := 0; i < g.maxRetries; i++ {
		// 使用INCRBY原子操作增加ID
		newID, err = g.client.IncrBy(ctx, key, g.step).Result()
		if err == nil {
			break
		}

		// 如果失败，等待一段时间后重试
		time.Sleep(g.retryInterval)
	}

	if err != nil {
		return fmt.Errorf("无法从Redis获取新的ID段: %v", err)
	}

	// 设置新的ID段
	segment.currentID = newID - g.step + 1
	segment.remainingIDs = g.step

	return nil
}
