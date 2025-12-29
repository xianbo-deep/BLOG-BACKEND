package model

import "time"

/* 声明表结构 */
type DailyArticleStat struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Date         time.Time `gorm:"column:date;type:date;not null" json:"date"`
	Path         string    `gorm:"column:path;type:text;not null" json:"path"`
	UV           int64     `gorm:"column:uv;not null;type:int" json:"uv"`
	PV           int64     `gorm:"column:pv;not null;type:int" json:"pv"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;type:timestamp" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;type:timestamp" json:"updated_at"`
	TotalLatency int64     `gorm:"column:total_latency;type:bigint" json:"total_latency"`
	LatencyCount int64     `gorm:"column:latency_count;type:bigint" json:"latency_count"`
}

/* 指定表名 */
func (d DailyArticleStat) TableName() string {
	return "daily_article_stats"
}
