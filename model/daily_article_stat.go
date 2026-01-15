package model

import "time"

/* 声明表结构 */
type DailyArticleStat struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Date         time.Time `gorm:"column:date;type:date;not null"`
	Path         string    `gorm:"column:path;type:text;not null"`
	UV           int64     `gorm:"column:uv;not null;type:int"`
	PV           int64     `gorm:"column:pv;not null;type:int"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;type:timestamp"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;type:timestamp"`
	TotalLatency int64     `gorm:"column:total_latency;type:bigint"`
	LatencyCount int64     `gorm:"column:latency_count;type:bigint"`
}

/* 指定表名 */
func (d DailyArticleStat) TableName() string {
	return "daily_article_stats"
}
