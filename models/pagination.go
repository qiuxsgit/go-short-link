package models

import (
	"strconv"

	"gorm.io/gorm"
)

// Paginate 分页查询
func Paginate(page, pageSize string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 转换页码和每页数量为整数
		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			pageInt = 1
		}

		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil || pageSizeInt < 1 {
			pageSizeInt = 10
		}

		// 限制每页最大数量
		if pageSizeInt > 100 {
			pageSizeInt = 100
		}

		// 计算偏移量
		offset := (pageInt - 1) * pageSizeInt

		// 应用分页
		return db.Offset(offset).Limit(pageSizeInt)
	}
}
