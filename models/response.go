package models

import (
	"time"
)

// FormattedShortLink 格式化后的短链接响应结构
type FormattedShortLink struct {
	ID          int64  `json:"id"`
	ShortCode   string `json:"shortCode"`
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
func (db *DBShortLink) ToFormattedShortLink() FormattedShortLink {
	return FormattedShortLink{
		ID:          db.ID,
		ShortCode:   db.ShortCode,
		OriginalURL: db.OriginalURL,
		CreatedAt:   FormatTime(db.CreatedAt),
		ExpiresAt:   FormatTime(db.ExpiresAt),
		AccessCount: db.AccessCount,
		LastAccess:  FormatTime(db.LastAccess),
	}
}