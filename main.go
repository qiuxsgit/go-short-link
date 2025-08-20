package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/qiuxsgit/go-short-link/handlers"
	"github.com/qiuxsgit/go-short-link/models"
)

func main() {
	// 创建短链接存储
	store := models.NewMemoryStore()

	// 获取基础URL，默认为localhost:8080
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080/"
	}

	// 创建短链接处理器
	shortLinkHandler := handlers.NewShortLinkHandler(store, baseURL)

	// 创建Gin路由
	router := gin.Default()

	// 注册API路由
	api := router.Group("/short-link")
	{
		api.POST("/create", shortLinkHandler.CreateShortLink)
	}

	// 注册重定向路由
	router.GET("/s/:code", shortLinkHandler.RedirectShortLink)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("短链接服务启动在 %s...\n", baseURL)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
