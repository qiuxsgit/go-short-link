package models

import (
	"time"

	"github.com/qiuxsgit/go-short-link/utils"
)

// FormattedShortLink 格式化后的短链接响应结构
type FormattedShortLink struct {
	ID          int64  `json:"id"`
	ShortCode   string `json:"shortCode"`
	ShortLink   string `json:"shortLink"`
	OriginalURL string `json:"originalUrl"`
	CreatedAt   string `json:"createdAt"`
	ExpiresAt   string `json:"expiresAt"`
	AccessCount int64  `json:"accessCount"`
	LastAccess  string `json:"lastAccess"`
}

// FormatTime 将时间格式化为指定格式
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05.000")
}

// ToFormattedShortLink 将DBShortLink转换为FormattedShortLink
// baseURL 是访问API服务的BaseURL，用于构建完整的短链接URL
func (db *DBShortLink) ToFormattedShortLink(baseURL string) FormattedShortLink {
	return FormattedShortLink{
		ID:          db.ID,
		ShortCode:   db.ShortCode,
		ShortLink:   utils.BuildShortLink(baseURL, db.ShortCode),
		OriginalURL: db.OriginalURL,
		CreatedAt:   FormatTime(db.CreatedAt),
		ExpiresAt:   FormatTime(db.ExpiresAt),
		AccessCount: db.AccessCount,
		LastAccess:  FormatTime(db.LastAccess),
	}
}
