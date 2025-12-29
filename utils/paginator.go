package utils

import (
	"Blog-Backend/dto/common"

	"gorm.io/gorm"
)

// 定义一个分页器，用了泛型
// 要注意没有指定表名称，这个泛型应当是数据表的实体，且绑定了数据表
func Paginate[T any](db *gorm.DB, pageReq common.PageRequest) (*common.PageResponse[T], error) {
	var result []T
	var total int64

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 获得页码，每页数据大小和偏移量
	page := pageReq.GetPage()
	pageSize := pageReq.GetPageSize()
	offset := (page - 1) * pageSize

	// 查结果
	if err := db.Offset(offset).Limit(pageSize).Find(&result).Error; err != nil {
		return nil, err
	}

	// 返回总页数，可以好好想想为什么这样
	totalPage := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &common.PageResponse[T]{
		List:      result,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}
