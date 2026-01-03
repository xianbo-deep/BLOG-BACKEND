package common

import "Blog-Backend/consts"

type PageRequest struct {
	Page     int `form:"page" json:"page"`           // 当前页
	PageSize int `form:"page_size" json:"page_size"` // 每页显示数量
}

type PageResponse[T any] struct {
	List      []T   `json:"list"`       // 数据列表
	Total     int64 `json:"total"`      // 总条数
	Page      int   `json:"page"`       // 当前页
	PageSize  int   `json:"page_size"`  // 每页大小
	TotalPage int   `json:"total_page"` // 总页数
}

func (p *PageRequest) GetPage() int {
	if p.Page <= 0 {
		return consts.DefaultPage
	}
	return p.Page
}

func (p *PageRequest) GetPageSize() int {
	if p.PageSize <= 0 || p.PageSize > consts.MaxPageSize {
		return consts.DefaultPageSize
	}
	return p.PageSize
}
