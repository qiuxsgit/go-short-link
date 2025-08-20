package models

import (
	"time"

	"github.com/qiuxsgit/go-short-link/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormStore 使用GORM和Redis ID生成器的存储实现
type GormStore struct {
	db    *gorm.DB
	cache *LRUCache
}

// NewGormStore 创建新的GORM存储
func NewGormStore(dsn string, cacheSize int, idGenerator *utils.RedisIDGenerator) (*GormStore, error) {
	// 连接MySQL数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// 注册ID生成器插件
	if err := idGenerator.Initialize(db); err != nil {
		return nil, err
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(&DBShortLink{}); err != nil {
		return nil, err
	}

	return &GormStore{
		db:    db,
		cache: NewLRUCache(cacheSize),
	}, nil
}

// Save 保存短链接到存储中
func (s *GormStore) Save(shortLink *ShortLink) error {
	// 保存到数据库 - ID会由GORM插件自动生成
	dbLink := FromShortLink(shortLink)
	if err := s.db.Create(dbLink).Error; err != nil {
		return err
	}

	// 更新ID (从数据库获取自动生成的ID)
	shortLink.ID = dbLink.ID

	// 保存到缓存
	s.cache.Put(shortLink.ShortCode, shortLink)
	return nil
}

// Get 根据短码获取短链接
func (s *GormStore) Get(shortCode string) (*ShortLink, error) {
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
func (s *GormStore) updateAccessCount(shortCode string) {
	s.db.Model(&DBShortLink{}).
		Where("short_code = ?", shortCode).
		Updates(map[string]interface{}{
			"access_count": gorm.Expr("access_count + 1"),
			"last_access":  time.Now(),
		})
}

// Close 关闭数据库连接
func (s *GormStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
