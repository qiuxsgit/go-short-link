package app

import (
	"context"
	"log"
	"os"

	"github.com/qiuxsgit/go-short-link/conf"
	"github.com/qiuxsgit/go-short-link/models"
	"github.com/qiuxsgit/go-short-link/tasks"
	"github.com/qiuxsgit/go-short-link/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// App 表示应用程序
type App struct {
	Config            *conf.Config
	Store             models.Store
	RedisClient       *redis.Client
	IDGeneratorPlugin *utils.RedisIDGenerator
	TaskScheduler     *tasks.Scheduler
	DB                *gorm.DB
}

// Initialize 初始化应用程序
func Initialize() (*App, error) {
	// 加载配置文件
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "conf/config.yaml"
	}

	config, err := conf.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// 创建Redis客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		PoolSize: config.Redis.PoolSize,
	})

	// 测试Redis连接
	ctx := context.Background()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	// 创建ID生成器插件
	idGeneratorPlugin := utils.NewRedisIDGenerator(redisClient, config.Redis.IDKeyPrefix, config.Redis.IDStep)

	// 创建短链接存储
	gormStore, err := models.NewGormStore(config.Database.GetDSN(), config.Cache.Capacity, idGeneratorPlugin)
	if err != nil {
		return nil, err
	}

	// 获取GORM DB实例
	db := gormStore.GetDB()

	// 创建定时任务调度器
	taskScheduler := tasks.NewScheduler(config)

	// 注册清理过期短链接任务
	if config.Tasks.CleanExpiredLinks.Enabled {
		cleanTask := tasks.NewCleanExpiredLinksTask(&config.Tasks.CleanExpiredLinks, db)
		taskScheduler.RegisterTask(cleanTask)
	}

	return &App{
		Config:            config,
		Store:             gormStore,
		RedisClient:       redisClient,
		IDGeneratorPlugin: idGeneratorPlugin,
		TaskScheduler:     taskScheduler,
		DB:                db,
	}, nil
}

// Cleanup 清理应用程序资源
func (a *App) Cleanup() {
	// 停止定时任务调度器
	if a.TaskScheduler != nil {
		a.TaskScheduler.Stop()
	}

	if a.Store != nil {
		if err := a.Store.Close(); err != nil {
			log.Printf("关闭存储失败: %v", err)
		}
	}

	if a.RedisClient != nil {
		if err := a.RedisClient.Close(); err != nil {
			log.Printf("关闭Redis连接失败: %v", err)
		}
	}
}
