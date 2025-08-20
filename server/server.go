package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qiuxsgit/go-short-link/api"
	"github.com/qiuxsgit/go-short-link/conf"
	"github.com/qiuxsgit/go-short-link/handlers"
	"github.com/qiuxsgit/go-short-link/models"
)

// Server 表示短链接服务器
type Server struct {
	config       *conf.Config
	store        models.Store
	adminServer  *http.Server
	accessServer *http.Server
}

// NewServer 创建一个新的服务器实例
func NewServer(config *conf.Config, store models.Store) *Server {
	return &Server{
		config: config,
		store:  store,
	}
}

// Initialize 初始化服务器
func (s *Server) Initialize() {
	// 创建管理API处理器
	adminHandler := handlers.NewShortLinkHandler(s.store, s.config.Server.Admin.BaseURL)

	// 创建访问API处理器
	accessHandler := handlers.NewShortLinkHandler(s.store, s.config.Server.Access.BaseURL)

	// 创建管理API路由
	adminRouter := gin.Default()
	api.SetupAdminRoutes(adminRouter, adminHandler, &s.config.Server.Admin)

	// 创建访问API路由
	accessRouter := gin.Default()
	api.SetupAccessRoutes(accessRouter, accessHandler)

	// 创建管理API服务器
	s.adminServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Server.Admin.Port),
		Handler: adminRouter,
	}

	// 创建访问API服务器
	s.accessServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Server.Access.Port),
		Handler: accessRouter,
	}
}

// Start 启动服务器
func (s *Server) Start() {
	// 启动管理API服务器
	go func() {
		log.Printf("管理API服务启动在 %s (端口: %d)...\n",
			s.config.Server.Admin.BaseURL, s.config.Server.Admin.Port)
		if err := s.adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("管理API服务器启动失败: %v", err)
		}
	}()

	// 启动访问API服务器
	go func() {
		log.Printf("访问API服务启动在 %s (端口: %d)...\n",
			s.config.Server.Access.BaseURL, s.config.Server.Access.Port)
		if err := s.accessServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("访问API服务器启动失败: %v", err)
		}
	}()
}

// Shutdown 关闭服务器
func (s *Server) Shutdown() error {
	log.Println("正在关闭服务器...")

	// 创建关闭超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭管理API服务器
	if err := s.adminServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("管理API服务器关闭失败: %v", err)
	}

	// 关闭访问API服务器
	if err := s.accessServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("访问API服务器关闭失败: %v", err)
	}

	log.Println("服务器已关闭")
	return nil
}
