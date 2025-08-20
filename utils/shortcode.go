package utils

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// GenerateShortCode 生成短链接代码
func GenerateShortCode(originalURL string) string {
	// 使用MD5哈希原始URL
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hasher.Write([]byte(uuid.New().String())) // 添加UUID以确保唯一性
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	// 取前8个字符作为短码
	shortCode := hash[:8]
	return shortCode
}

// BuildShortLink 构建完整的短链接URL
func BuildShortLink(baseURL, shortCode string) string {
	// 确保baseURL以/结尾
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}
	return fmt.Sprintf("%ss/%s", baseURL, shortCode)
}
