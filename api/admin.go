package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qiuxsgit/go-short-link/conf"
	"github.com/qiuxsgit/go-short-link/handlers"
	"github.com/qiuxsgit/go-short-link/utils"
)

// JWTAuthMiddleware 创建JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式无效"})
			c.Abort()
			return
		}

		// 解析令牌
		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			if err == utils.ErrExpiredToken {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌已过期"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌无效"})
			}
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

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
func SetupAdminRoutes(router *gin.Engine, shortLinkHandler *handlers.ShortLinkHandler, adminHandler *handlers.AdminHandler, config *conf.AdminServerConfig) {
	// 添加IP白名单中间件
	router.Use(IPWhitelistMiddleware(config))

	// 公共API路由（无需认证）
	publicAPI := router.Group("/api")
	{
		// 管理员登录
		publicAPI.POST("/login", adminHandler.Login)

		// 创建短链接（保留原有功能）
		publicAPI.POST("/short-link/create", shortLinkHandler.CreateShortLink)
	}

	// 需要认证的API路由
	privateAPI := router.Group("/api")
	privateAPI.Use(JWTAuthMiddleware())
	{
		// 短链接管理
		linkAPI := privateAPI.Group("/short-link")
		{
			// 获取有效短链接列表
			linkAPI.GET("/list", adminHandler.GetShortLinks)

			// 获取历史短链接列表
			linkAPI.GET("/history", adminHandler.GetHistoryLinks)

			// 删除短链接（移动到历史表）
			linkAPI.DELETE("/:id", adminHandler.DeleteShortLink)
		}
	}
}
