package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/static"
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
	// 设置Gin模式
	if s.config.Server.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 获取GORM存储实例
	gormStore, ok := s.store.(*models.GormStore)
	if !ok {
		log.Fatal("存储必须是GormStore类型")
	}

	// 创建管理API处理器
	adminHandler := handlers.NewShortLinkHandler(s.store, s.config.Server.Admin.BaseURL)

	// 创建访问API处理器
	accessHandler := handlers.NewShortLinkHandler(s.store, s.config.Server.Access.BaseURL)

	// 创建管理员处理器
	adminUserHandler := handlers.NewAdminHandler(gormStore, s.config)

	// 创建管理API路由
	adminRouter := gin.Default()

	// 添加静态文件服务
	adminRouter.Use(func() gin.HandlerFunc {
		// 检查static目录是否存在
		staticDir := "static"
		if _, err := os.Stat(staticDir); os.IsNotExist(err) {
			log.Printf("警告: 静态文件目录 %s 不存在，前端将不可用", staticDir)
			return func(c *gin.Context) {
				c.Next()
			}
		}

		// 使用静态文件中间件
		handler := static.Serve("/", static.LocalFile(staticDir, false))
		return func(c *gin.Context) {
			// 如果是API请求，跳过静态文件处理
			if strings.HasPrefix(c.Request.URL.Path, "/api") ||
				strings.HasPrefix(c.Request.URL.Path, "/s/") {
				c.Next()
				return
			}

			handler(c)

			// 如果找不到静态文件，返回index.html（支持SPA路由）
			if c.Writer.Status() == http.StatusNotFound {
				c.File(filepath.Join(staticDir, "index.html"))
			}
		}
	}())

	api.SetupAdminRoutes(adminRouter, adminHandler, adminUserHandler, &s.config.Server.Admin)

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
		// 创建一个通道来通知服务器已启动
		started := make(chan struct{})

		// 在单独的goroutine中启动服务器
		go func() {
			// 监听端口
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Server.Admin.Port))
			if err != nil {
				log.Fatalf("管理API服务器启动失败: %v", err)
				return
			}

			// 通知主goroutine服务器已准备好接受连接
			close(started)

			// 使用已创建的监听器提供服务
			if err := s.adminServer.Serve(listener); err != nil && err != http.ErrServerClosed {
				log.Fatalf("管理API服务器运行失败: %v", err)
			}
		}()

		// 等待服务器启动
		<-started
		log.Printf("Listening and serving HTTP on Port on %s (端口: %d)...\n",
			s.config.Server.Admin.BaseURL, s.config.Server.Admin.Port)
	}()

	// 启动访问API服务器
	go func() {
		// 创建一个通道来通知服务器已启动
		started := make(chan struct{})

		// 在单独的goroutine中启动服务器
		go func() {
			// 监听端口
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Server.Access.Port))
			if err != nil {
				log.Fatalf("访问API服务器启动失败: %v", err)
				return
			}

			// 通知主goroutine服务器已准备好接受连接
			close(started)

			// 使用已创建的监听器提供服务
			if err := s.accessServer.Serve(listener); err != nil && err != http.ErrServerClosed {
				log.Fatalf("访问API服务器运行失败: %v", err)
			}
		}()

		// 等待服务器启动
		<-started
		log.Printf("Listening and serving HTTP on Port %s (端口: %d)...\n",
			s.config.Server.Access.BaseURL, s.config.Server.Access.Port)
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
