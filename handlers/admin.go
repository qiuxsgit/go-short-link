package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qiuxsgit/go-short-link/conf"
	"github.com/qiuxsgit/go-short-link/models"
	"github.com/qiuxsgit/go-short-link/utils"
)

// AdminHandler 处理管理员相关的请求
type AdminHandler struct {
	db     *models.GormStore
	config *conf.Config
}

// NewAdminHandler 创建一个新的管理员处理器
func NewAdminHandler(db *models.GormStore, config *conf.Config) *AdminHandler {
	return &AdminHandler{
		db:     db,
		config: config,
	}
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	UserID   int64  `json:"userId"`
}

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=NewPassword"`
}

// Login 管理员登录
func (h *AdminHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 查询管理员
	var admin models.SysAdmin
	if err := h.db.GetDB().Where("username = ?", req.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if !models.CheckPassword(admin.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 更新最后登录时间
	h.db.GetDB().Model(&admin).Update("last_login", time.Now())

	// 生成JWT令牌
	expireDuration := time.Duration(h.config.JWT.ExpireHours) * time.Hour
	token, err := utils.GenerateToken(admin.ID, admin.Username, expireDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, LoginResponse{
		Token:    token,
		Username: admin.Username,
		UserID:   admin.ID,
	})
}

// GetShortLinks 获取短链接列表
func (h *AdminHandler) GetShortLinks(c *gin.Context) {
	var links []models.DBShortLink
	query := h.db.GetDB().Order("created_at DESC")

	// 分页参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	// 过滤参数
	if shortCode := c.Query("shortCode"); shortCode != "" {
		query = query.Where("short_code LIKE ?", "%"+shortCode+"%")
	}
	if originalURL := c.Query("originalUrl"); originalURL != "" {
		query = query.Where("original_url LIKE ?", "%"+originalURL+"%")
	}
	if status := c.Query("status"); status != "" {
		if status == "active" {
			query = query.Where("expires_at > ?", time.Now())
		} else if status == "expired" {
			query = query.Where("expires_at <= ?", time.Now())
		}
	}

	// 执行查询
	var total int64
	countResult := query.Model(&models.DBShortLink{}).Count(&total)
	if countResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "计数查询失败: " + countResult.Error.Error(),
		})
		return
	}

	findResult := query.Scopes(models.Paginate(page, pageSize)).Find(&links)
	if findResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询数据失败: " + findResult.Error.Error(),
		})
		return
	}

	// 格式化日期
	formattedLinks := make([]models.FormattedShortLink, len(links))
	for i, link := range links {
		formattedLinks[i] = link.ToFormattedShortLink()
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"links": formattedLinks,
	})
}

// GetHistoryLinks 获取历史短链接列表
func (h *AdminHandler) GetHistoryLinks(c *gin.Context) {
	// 获取月份参数
	month := c.DefaultQuery("month", time.Now().Format("0601")) // 默认当前月份，格式为YYMM

	// 打印调试信息
	c.Set("debug_month", month)

	// 构建历史表名
	historyTable := h.config.Tasks.CleanExpiredLinks.HistoryTablePrefix + month

	// 打印调试信息
	c.Set("debug_table", historyTable)

	// 检查表是否存在
	var count int64
	result := h.db.GetDB().Raw(`
		SELECT COUNT(1) FROM information_schema.tables 
		WHERE table_schema = ? AND table_name = ?
	`, h.config.Database.DBName, historyTable).Scan(&count)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "检查表是否存在时出错: " + result.Error.Error(),
			"debug_month": month,
			"debug_table": historyTable,
		})
		return
	}

	if count == 0 {
		c.JSON(http.StatusOK, gin.H{
			"total":        0,
			"links":        []interface{}{},
			"debug_month":  month,
			"debug_table":  historyTable,
			"debug_exists": false,
		})
		return
	}

	// 分页参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	// 查询历史表
	var links []models.DBShortLink
	query := h.db.GetDB().Table(historyTable).Order("created_at DESC")

	// 过滤参数
	if shortCode := c.Query("shortCode"); shortCode != "" {
		query = query.Where("short_code LIKE ?", "%"+shortCode+"%")
	}
	if originalURL := c.Query("originalUrl"); originalURL != "" {
		query = query.Where("original_url LIKE ?", "%"+originalURL+"%")
	}

	// 执行查询
	var total int64
	countResult := query.Count(&total)
	if countResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "计数查询失败: " + countResult.Error.Error(),
			"debug_month": month,
			"debug_table": historyTable,
		})
		return
	}

	findResult := query.Scopes(models.Paginate(page, pageSize)).Find(&links)
	if findResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "查询数据失败: " + findResult.Error.Error(),
			"debug_month": month,
			"debug_table": historyTable,
		})
		return
	}

	// 格式化日期
	formattedLinks := make([]models.FormattedShortLink, len(links))
	for i, link := range links {
		formattedLinks[i] = link.ToFormattedShortLink()
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"total":        total,
		"links":        formattedLinks,
		"debug_month":  month,
		"debug_table":  historyTable,
		"debug_exists": true,
		"debug_count":  len(links),
	})
}

// DeleteShortLink 删除短链接（移动到历史表）
func (h *AdminHandler) DeleteShortLink(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的短链接ID"})
		return
	}

	// 查询短链接
	var link models.DBShortLink
	if err := h.db.GetDB().Where("id = ?", id).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "短链接不存在"})
		return
	}

	// 获取当前月份
	currentMonth := time.Now().Format("0601") // 格式为YYMM

	// 构建历史表名
	historyTable := h.config.Tasks.CleanExpiredLinks.HistoryTablePrefix + currentMonth

	// 确保历史表存在
	h.db.GetDB().Exec(`CREATE TABLE IF NOT EXISTS ` + historyTable + ` LIKE short_links`)

	// 开始事务
	tx := h.db.GetDB().Begin()

	// 将短链接插入历史表
	if err := tx.Exec("INSERT INTO "+historyTable+" SELECT * FROM short_links WHERE id = ?", id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移动短链接到历史表失败"})
		return
	}

	// 从主表删除短链接
	if err := tx.Delete(&link).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除短链接失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "短链接已成功删除"})
}

// ChangePassword 修改密码
func (h *AdminHandler) ChangePassword(c *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 解析请求参数
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 查询用户
	var admin models.SysAdmin
	if err := h.db.GetDB().Where("id = ?", userID).First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 验证当前密码
	if !models.CheckPassword(admin.Password, req.CurrentPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "当前密码错误"})
		return
	}

	// 生成新密码的哈希
	hashedPassword, err := models.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 更新密码
	if err := h.db.GetDB().Model(&admin).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}
