package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// 密钥，实际应用中应该从配置文件或环境变量中获取
	jwtSecret = []byte("go-short-link-secret-key")

	// ErrInvalidToken 无效的令牌
	ErrInvalidToken = errors.New("无效的令牌")
	// ErrExpiredToken 令牌已过期
	ErrExpiredToken = errors.New("令牌已过期")
)

// JWTClaims 自定义JWT声明
type JWTClaims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID int64, username string, expireDuration time.Duration) (string, error) {
	// 设置过期时间
	expireTime := time.Now().Add(expireDuration)

	// 创建声明
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-short-link",
			Subject:   username,
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*JWTClaims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// 验证令牌
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// 获取声明
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// SetJWTSecret 设置JWT密钥
func SetJWTSecret(secret string) {
	if secret != "" {
		jwtSecret = []byte(secret)
	}
}
