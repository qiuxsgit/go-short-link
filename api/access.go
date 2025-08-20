package api

import (
	"github.com/gin-gonic/gin"
	"github.com/qiuxsgit/go-short-link/handlers"
)

// SetupAccessRoutes 设置访问API路由
func SetupAccessRoutes(router *gin.Engine, handler *handlers.ShortLinkHandler) {
	// 注册重定向路由
	router.GET("/s/:code", handler.RedirectShortLink)
}
