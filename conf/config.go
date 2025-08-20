package conf

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 应用程序配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Cache    CacheConfig    `yaml:"cache"`
	Tasks    TasksConfig    `yaml:"tasks"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Admin   AdminServerConfig  `yaml:"admin"`
	Access  AccessServerConfig `yaml:"access"`
	GinMode string             `yaml:"ginMode"` // gin模式: debug或release
}

// AdminServerConfig 管理API服务配置
type AdminServerConfig struct {
	Port        int      `yaml:"port"`
	BaseURL     string   `yaml:"baseURL"`
	IPWhitelist []string `yaml:"ipWhitelist"`
}

// AccessServerConfig 访问API服务配置
type AccessServerConfig struct {
	Port    int    `yaml:"port"`
	BaseURL string `yaml:"baseURL"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `yaml:"driver"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	Charset         string `yaml:"charset"`
	ParseTime       bool   `yaml:"parseTime"`
	Loc             string `yaml:"loc"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr        string `yaml:"addr"`
	Password    string `yaml:"password"`
	DB          int    `yaml:"db"`
	PoolSize    int    `yaml:"poolSize"`
	IDKeyPrefix string `yaml:"idKeyPrefix"`
	IDStep      int64  `yaml:"idStep"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Type     string `yaml:"type"`
	Capacity int    `yaml:"capacity"`
}

// TasksConfig 定时任务配置
type TasksConfig struct {
	CleanExpiredLinks CleanExpiredLinksConfig `yaml:"cleanExpiredLinks"`
}

// CleanExpiredLinksConfig 清理过期短链接任务配置
type CleanExpiredLinksConfig struct {
	Cron               string `yaml:"cron"`
	Enabled            bool   `yaml:"enabled"`
	BatchSize          int    `yaml:"batchSize"`
	HistoryTablePrefix string `yaml:"historyTablePrefix"`
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 如果未指定配置文件路径，则使用默认路径
	if configPath == "" {
		configPath = "conf/config.yaml"
	}

	// 确保配置文件存在
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("无法获取配置文件的绝对路径: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", absPath)
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件: %v", err)
	}

	// 解析YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("无法解析配置文件: %v", err)
	}

	return &config, nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username, c.Password, c.Host, c.Port, c.DBName, c.Charset, c.ParseTime, c.Loc)
}

// IsIPAllowed 检查IP是否在白名单中
func (c *AdminServerConfig) IsIPAllowed(ip string) bool {
	// 如果白名单为空，允许所有IP
	if len(c.IPWhitelist) == 0 {
		return true
	}

	// 解析客户端IP
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}

	// 检查IP是否在白名单中
	for _, allowedIP := range c.IPWhitelist {
		// 检查是否是CIDR格式
		if _, ipNet, err := net.ParseCIDR(allowedIP); err == nil {
			if ipNet.Contains(clientIP) {
				return true
			}
			continue
		}

		// 检查是否是单个IP
		if net.ParseIP(allowedIP).Equal(clientIP) {
			return true
		}
	}

	return false
}
