package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qiuxsgit/go-short-link/models"
	"github.com/qiuxsgit/go-short-link/utils"
)

// ShortLinkHandler 处理短链接相关的请求
type ShortLinkHandler struct {
	store   models.Store
	baseURL string
}

// NewShortLinkHandler 创建一个新的短链接处理器
func NewShortLinkHandler(store models.Store, baseURL string) *ShortLinkHandler {
	return &ShortLinkHandler{
		store:   store,
		baseURL: baseURL,
	}
}

// CreateShortLink 创建短链接
func (h *ShortLinkHandler) CreateShortLink(c *gin.Context) {
	var req models.CreateShortLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 生成短链接代码
	shortCode := utils.GenerateShortCode(req.Link)

	// 创建短链接记录
	shortLink := &models.ShortLink{
		OriginalURL: req.Link,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(time.Duration(req.Expire) * time.Second),
	}

	// 保存到存储
	if err := h.store.Save(shortLink); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建短链接失败"})
		return
	}

	// 构建完整的短链接URL
	fullShortLink := utils.BuildShortLink(h.baseURL, shortCode)

	// 返回响应
	c.JSON(http.StatusOK, models.CreateShortLinkResponse{
		ShortLink: fullShortLink,
	})
}

// RedirectShortLink 重定向短链接到原始URL
func (h *ShortLinkHandler) RedirectShortLink(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title":   "Page Not Found",
			"message": "无效的短链接代码",
		})
		return
	}

	// 从存储中获取短链接
	shortLink, err := h.store.Get(shortCode)
	if err != nil {
		// 返回404页面
		c.HTML(http.StatusNotFound, "", gin.H{})
		// 如果没有模板，则直接返回HTML
		if c.Writer.Status() != http.StatusNotFound {
			return
		}

		// 直接返回简单的HTML页面
		notFoundHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Page Not Found</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            text-align: center;
            padding-top: 100px;
            background-color: #f7f7f7;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #fff;
            border-radius: 5px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #e74c3c;
        }
        p {
            color: #7f8c8d;
            font-size: 18px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>404 - Page Not Found</h1>
        <p>抱歉，您访问的短链接不存在或已过期。</p>
    </div>
</body>
</html>`
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusNotFound, notFoundHTML)
		return
	}

	// 重定向到原始URL
	c.Redirect(http.StatusTemporaryRedirect, shortLink.OriginalURL)
}
