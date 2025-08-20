package models

import (
	"time"
)

// ShortLink 表示短链接的数据结构
type ShortLink struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"originalUrl"`
	ShortCode   string    `json:"shortCode"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// CreateShortLinkRequest 创建短链接的请求结构
type CreateShortLinkRequest struct {
	Link   string `json:"link" binding:"required"`
	Expire int    `json:"expire" binding:"required"`
}

// CreateShortLinkResponse 创建短链接的响应结构
type CreateShortLinkResponse struct {
	ShortLink string `json:"shortLink"`
}
