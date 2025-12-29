package model

type PathStatRes struct {
	Path       string `gorm:"column:path;primaryKey"`
	UV         int64  `gorm:"column:uv"`
	PV         int64  `gorm:"column:pv"`
	AvgLatency int64  `gorm:"column:avg_latency"`
}

// 绑定视图名
func (p *PathStatRes) TableName() string {
	return "article_ranking"
}
