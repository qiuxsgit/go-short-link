package app

import (
	"context"
	"log"
	"os"

	"github.com/qiuxsgit/go-short-link/conf"
	"github.com/qiuxsgit/go-short-link/models"
	"github.com/qiuxsgit/go-short-link/utils"
	"github.com/redis/go-redis/v9"
)

// App 表示应用程序
type App struct {
	Config            *conf.Config
	Store             models.Store
	RedisClient       *redis.Client
	IDGeneratorPlugin *utils.RedisIDGenerator
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
	store, err := models.NewGormStore(config.Database.GetDSN(), config.Cache.Capacity, idGeneratorPlugin)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:            config,
		Store:             store,
		RedisClient:       redisClient,
		IDGeneratorPlugin: idGeneratorPlugin,
	}, nil
}

// Cleanup 清理应用程序资源
func (a *App) Cleanup() {
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
