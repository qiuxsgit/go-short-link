package models

import (
	"container/list"
	"errors"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrLinkNotFound = errors.New("短链接不存在或已过期")
)

// Store 是短链接存储的接口
type Store interface {
	Save(shortLink *ShortLink) error
	Get(shortCode string) (*ShortLink, error)
	Close() error
}

// DBShortLink 是数据库中短链接的模型
type DBShortLink struct {
	ID          string    `gorm:"primaryKey;type:varchar(64)"`
	ShortCode   string    `gorm:"uniqueIndex;type:varchar(16)"`
	OriginalURL string    `gorm:"type:text"`
	CreatedAt   time.Time
	ExpiresAt   time.Time
	AccessCount int64     `gorm:"default:0"`
	LastAccess  time.Time
}

// TableName 设置表名
func (DBShortLink) TableName() string {
	return "short_links"
}

// ToShortLink 转换为ShortLink模型
func (db *DBShortLink) ToShortLink() *ShortLink {
	return &ShortLink{
		ID:          db.ID,
		ShortCode:   db.ShortCode,
		OriginalURL: db.OriginalURL,
		CreatedAt:   db.CreatedAt,
		ExpiresAt:   db.ExpiresAt,
	}
}

// FromShortLink 从ShortLink模型转换
func FromShortLink(sl *ShortLink) *DBShortLink {
	return &DBShortLink{
		ID:          sl.ID,
		ShortCode:   sl.ShortCode,
		OriginalURL: sl.OriginalURL,
		CreatedAt:   sl.CreatedAt,
		ExpiresAt:   sl.ExpiresAt,
		LastAccess:  time.Now(),
	}
}

// LRUCache 实现LRU缓存
type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
	mutex    sync.RWMutex
}

// cacheItem 缓存项
type cacheItem struct {
	key   string
	value *ShortLink
}

// NewLRUCache 创建新的LRU缓存
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

// Get 从缓存获取值
func (c *LRUCache) Get(key string) (*ShortLink, bool) {
	c.mutex.RLock()
	if elem, ok := c.cache[key]; ok {
		c.mutex.RUnlock()
		c.mutex.Lock()
		c.list.MoveToFront(elem)
		c.mutex.Unlock()
		return elem.Value.(*cacheItem).value, true
	}
	c.mutex.RUnlock()
	return nil, false
}

// Put 放入缓存
func (c *LRUCache) Put(key string, value *ShortLink) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*cacheItem).value = value
		return
	}

	if c.list.Len() >= c.capacity {
		oldest := c.list.Back()
		if oldest != nil {
			delete(c.cache, oldest.Value.(*cacheItem).key)
			c.list.Remove(oldest)
		}
	}

	elem := c.list.PushFront(&cacheItem{key: key, value: value})
	c.cache[key] = elem
}

// HybridStore 混合存储实现（MySQL + 内存缓存）
type HybridStore struct {
	db    *gorm.DB
	cache *LRUCache
}

// NewHybridStore 创建新的混合存储
func NewHybridStore(dsn string, cacheSize int) (*HybridStore, error) {
	// 连接MySQL数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(&DBShortLink{}); err != nil {
		return nil, err
	}

	return &HybridStore{
		db:    db,
		cache: NewLRUCache(cacheSize),
	}, nil
}

// Save 保存短链接到存储中
func (s *HybridStore) Save(shortLink *ShortLink) error {
	// 保存到数据库
	dbLink := FromShortLink(shortLink)
	if err := s.db.Create(dbLink).Error; err != nil {
		return err
	}

	// 保存到缓存
	s.cache.Put(shortLink.ShortCode, shortLink)
	return nil
}

// Get 根据短码获取短链接
func (s *HybridStore) Get(shortCode string) (*ShortLink, error) {
	// 先从缓存获取
	if link, found := s.cache.Get(shortCode); found {
		// 检查链接是否过期
		if time.Now().After(link.ExpiresAt) {
			return nil, ErrLinkNotFound
		}

		// 异步更新访问计数
		go s.updateAccessCount(shortCode)

		return link, nil
	}

	// 从数据库获取
	var dbLink DBShortLink
	if err := s.db.Where("short_code = ?", shortCode).First(&dbLink).Error; err != nil {
		return nil, ErrLinkNotFound
	}

	// 检查链接是否过期
	if time.Now().After(dbLink.ExpiresAt) {
		return nil, ErrLinkNotFound
	}

	// 转换为ShortLink
	link := dbLink.ToShortLink()

	// 添加到缓存
	s.cache.Put(shortCode, link)

	// 异步更新访问计数
	go s.updateAccessCount(shortCode)

	return link, nil
}

// updateAccessCount 更新访问计数
func (s *HybridStore) updateAccessCount(shortCode string) {
	s.db.Model(&DBShortLink{}).
		Where("short_code = ?", shortCode).
		Updates(map[string]interface{}{
			"access_count": gorm.Expr("access_count + 1"),
			"last_access":  time.Now(),
		})
}

// Close 关闭数据库连接
func (s *HybridStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// MemoryStore 是一个基于内存的短链接存储实现（保留原有实现作为备用）
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

// Close 实现Store接口
func (s *MemoryStore) Close() error {
	return nil
}
