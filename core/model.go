package core

import "time"

/* 声明表结构 */
type VisitLog struct {
	ID         int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt  time.Time `gorm:"column:created_at;not null" json:"created_at"`
	VisitTime  time.Time `gorm:"column:visit_time;not null" json:"visit_time"`
	ClientTime time.Time `gorm:"column:client_time" json:"client_time"`
	Path       string    `gorm:"column:path;type;not null" json:"path"`
	Method     string    `gorm:"column:method" json:"method"`
	IP         string    `gorm:"column:ip" json:"ip"`
	UserAgent  string    `gorm:"column:user_agent" json:"user_agent"`
	Country    string    `gorm:"column:country" json:"country"`
	Region     string    `gorm:"column:region" json:"region"`
	Referer    string    `gorm:"column:referer" json:"referer"`
}

type DailyArticleStat struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Date      time.Time `gorm:"column:date;type:date;not null" json:"date"`
	Path      string    `gorm:"column:path;type;not null" json:"path"`
	UV        int64     `gorm:"column:uv" json:"uv"`
	PV        int64     `gorm:"column:pv" json:"pv"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAT time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

/* 指定表名 */
func (v VisitLog) TableName() string {
	return "visit_logs"
}

func (d DailyArticleStat) TableName() string {
	return "daily_article_stats"
}
