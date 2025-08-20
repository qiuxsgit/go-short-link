package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/qiuxsgit/go-short-link/app"
	"github.com/qiuxsgit/go-short-link/server"
)

func main() {
	// 初始化应用程序
	application, err := app.Initialize()
	if err != nil {
		log.Fatalf("初始化应用程序失败: %v", err)
	}
	defer application.Cleanup()

	// 创建并初始化服务器
	srv := server.NewServer(application.Config, application.Store)
	srv.Initialize()

	// 启动定时任务调度器
	if application.TaskScheduler != nil {
		application.TaskScheduler.Start()
	}

	// 启动服务器
	srv.Start()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 关闭服务器
	if err := srv.Shutdown(); err != nil {
		log.Fatalf("关闭服务器失败: %v", err)
	}
}
