package tasks

import (
	"fmt"
	"log"
	"time"

	"github.com/qiuxsgit/go-short-link/conf"
	"gorm.io/gorm"
)

// CleanExpiredLinksTask 清理过期短链接的任务
type CleanExpiredLinksTask struct {
	config *conf.CleanExpiredLinksConfig
	db     *gorm.DB
}

// NewCleanExpiredLinksTask 创建一个新的清理过期短链接任务
func NewCleanExpiredLinksTask(config *conf.CleanExpiredLinksConfig, db *gorm.DB) *CleanExpiredLinksTask {
	return &CleanExpiredLinksTask{
		config: config,
		db:     db,
	}
}

// Name 返回任务名称
func (t *CleanExpiredLinksTask) Name() string {
	return "CleanExpiredLinks"
}

// IsEnabled 检查任务是否启用
func (t *CleanExpiredLinksTask) IsEnabled() bool {
	return t.config.Enabled
}

// Schedule 返回任务的调度表达式
func (t *CleanExpiredLinksTask) Schedule() string {
	return t.config.Cron
}

// Run 执行任务
func (t *CleanExpiredLinksTask) Run() error {
	log.Println("开始清理过期短链接...")

	// 获取当前时间
	now := time.Now()

	// 创建历史表名称，格式为：short_links_history_YYMM
	historyTableName := fmt.Sprintf("%s%s",
		t.config.HistoryTablePrefix,
		now.Format("0601")) // 格式为YYMM，如2508表示2025年8月

	// 确保历史表存在
	if err := t.ensureHistoryTableExists(historyTableName); err != nil {
		return fmt.Errorf("确保历史表存在失败: %v", err)
	}

	// 批量处理过期的短链接
	var processedCount int64
	batchSize := t.config.BatchSize
	if batchSize <= 0 {
		batchSize = 1000 // 默认批处理大小
	}

	// 开始事务
	tx := t.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开始事务失败: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("任务执行过程中发生panic: %v", r)
		}
	}()

	// 查询过期的短链接
	var expiredLinks []map[string]interface{}
	if err := tx.Table("short_links").
		Where("expires_at < ?", now).
		Limit(batchSize).
		Find(&expiredLinks).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("查询过期短链接失败: %v", err)
	}

	if len(expiredLinks) == 0 {
		tx.Rollback() // 没有需要处理的记录，回滚事务
		log.Println("没有找到过期的短链接")
		return nil
	}

	// 将过期的短链接移动到历史表
	for _, link := range expiredLinks {
		// 插入到历史表
		if err := tx.Table(historyTableName).Create(link).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("插入历史表失败: %v", err)
		}

		// 从短链接表中删除
		if err := tx.Table("short_links").
			Where("id = ?", link["id"]).
			Delete(nil).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("删除过期短链接失败: %v", err)
		}

		processedCount++
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	log.Printf("成功清理 %d 个过期短链接，移动到历史表 %s", processedCount, historyTableName)
	return nil
}

// ensureHistoryTableExists 确保历史表存在
func (t *CleanExpiredLinksTask) ensureHistoryTableExists(historyTableName string) error {
	// 检查表是否存在
	var count int64
	if err := t.db.Raw(
		"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?",
		historyTableName,
	).Count(&count).Error; err != nil {
		return err
	}

	// 如果表不存在，创建表
	if count == 0 {
		// 使用与short_links相同的表结构创建历史表
		sql := fmt.Sprintf(`
			CREATE TABLE %s LIKE short_links
		`, historyTableName)

		if err := t.db.Exec(sql).Error; err != nil {
			return fmt.Errorf("创建历史表失败: %v", err)
		}

		log.Printf("已创建历史表: %s", historyTableName)
	}

	return nil
}
