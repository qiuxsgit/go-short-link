package app

import (
	"log"
	"os"

	"github.com/qiuxsgit/go-short-link/conf"
	"github.com/qiuxsgit/go-short-link/models"
)

// App 表示应用程序
type App struct {
	Config *conf.Config
	Store  models.Store
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

	// 创建短链接存储
	store, err := models.NewHybridStore(config.Database.GetDSN(), config.Cache.Capacity)
	if err != nil {
		return nil, err
	}

	return &App{
		Config: config,
		Store:  store,
	}, nil
}

// Cleanup 清理应用程序资源
func (a *App) Cleanup() {
	if a.Store != nil {
		if err := a.Store.Close(); err != nil {
			log.Printf("关闭存储失败: %v", err)
		}
	}
}
