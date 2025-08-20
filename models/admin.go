package models

import (
	"crypto/rand"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SysAdmin 系统管理员模型
type SysAdmin struct {
	ID        int64     `gorm:"primaryKey;type:bigint(20);not null;auto_increment:false"`
	Username  string    `gorm:"uniqueIndex;type:varchar(50);not null"`
	Password  string    `gorm:"type:varchar(100);not null"`
	LastLogin time.Time `gorm:"type:datetime"`
	CreatedAt time.Time `gorm:"type:datetime;not null"`
	UpdatedAt time.Time `gorm:"type:datetime;not null"`
}

// TableName 设置表名
func (SysAdmin) TableName() string {
	return "sys_admin"
}

// GenerateRandomPassword 生成随机6位数字密码
func GenerateRandomPassword() string {
	b := make([]byte, 3)
	rand.Read(b)
	// 将随机字节转换为6位数字
	num := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	num = num%900000 + 100000 // 确保是6位数字
	return fmt.Sprintf("%d", num)
}

// HashPassword 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 检查密码是否正确
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// EnsureAdminExists 确保至少存在一个管理员账户
func EnsureAdminExists(db *gorm.DB) (string, error) {
	var count int64
	db.Model(&SysAdmin{}).Count(&count)

	if count == 0 {
		// 生成随机密码
		password := GenerateRandomPassword()
		hashedPassword, err := HashPassword(password)
		if err != nil {
			return "", err
		}

		// 创建默认管理员
		admin := &SysAdmin{
			Username:  "admin",
			Password:  hashedPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(admin).Error; err != nil {
			return "", err
		}

		return password, nil
	}

	return "", nil
}
