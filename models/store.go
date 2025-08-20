package models

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrLinkNotFound = errors.New("短链接不存在或已过期")
)

// Store 是短链接存储的接口
type Store interface {
	Save(shortLink *ShortLink) error
	Get(shortCode string) (*ShortLink, error)
}

// MemoryStore 是一个基于内存的短链接存储实现
type MemoryStore struct {
	links map[string]*ShortLink
	mutex sync.RWMutex
}

// NewMemoryStore 创建一个新的内存存储
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		links: make(map[string]*ShortLink),
	}
}

// Save 保存短链接到存储中
func (s *MemoryStore) Save(shortLink *ShortLink) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.links[shortLink.ShortCode] = shortLink
	return nil
}

// Get 根据短码获取短链接
func (s *MemoryStore) Get(shortCode string) (*ShortLink, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	link, exists := s.links[shortCode]
	if !exists {
		return nil, ErrLinkNotFound
	}

	// 检查链接是否过期
	if time.Now().After(link.ExpiresAt) {
		return nil, ErrLinkNotFound
	}

	return link, nil
}
