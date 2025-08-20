package utils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// IDGenerator 使用Redis生成唯一ID
type IDGenerator struct {
	client        *redis.Client
	key           string
	step          int64
	currentID     int64
	remainingIDs  int64
	mutex         sync.Mutex
	retryInterval time.Duration
	maxRetries    int
}

// NewIDGenerator 创建一个新的ID生成器
func NewIDGenerator(client *redis.Client, key string, step int64) *IDGenerator {
	return &IDGenerator{
		client:        client,
		key:           key,
		step:          step,
		currentID:     0,
		remainingIDs:  0,
		retryInterval: 100 * time.Millisecond,
		maxRetries:    5,
	}
}

// NextID 生成下一个唯一ID
func (g *IDGenerator) NextID() (int64, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// 如果没有剩余ID，从Redis获取新的ID段
	if g.remainingIDs == 0 {
		if err := g.fetchNewIDSegment(); err != nil {
			return 0, err
		}
	}

	// 使用当前ID并减少剩余ID计数
	id := g.currentID
	g.currentID++
	g.remainingIDs--

	return id, nil
}

// fetchNewIDSegment 从Redis获取新的ID段
func (g *IDGenerator) fetchNewIDSegment() error {
	ctx := context.Background()
	var err error
	var newID int64

	// 重试逻辑
	for i := 0; i < g.maxRetries; i++ {
		// 使用INCRBY原子操作增加ID
		newID, err = g.client.IncrBy(ctx, g.key, g.step).Result()
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
	g.currentID = newID - g.step + 1
	g.remainingIDs = g.step

	return nil
}