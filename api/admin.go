package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qiuxsgit/go-short-link/conf"
	"github.com/qiuxsgit/go-short-link/handlers"
)

// IPWhitelistMiddleware 创建IP白名单中间件
func IPWhitelistMiddleware(config *conf.AdminServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if !config.IsIPAllowed(clientIP) {
			c.JSON(http.StatusForbidden, gin.H{"error": "IP不在白名单中"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SetupAdminRoutes 设置管理API路由
func SetupAdminRoutes(router *gin.Engine, handler *handlers.ShortLinkHandler, config *conf.AdminServerConfig) {
	// 添加IP白名单中间件
	router.Use(IPWhitelistMiddleware(config))

	// 注册API路由
	api := router.Group("/short-link")
	{
		api.POST("/create", handler.CreateShortLink)
	}
}
