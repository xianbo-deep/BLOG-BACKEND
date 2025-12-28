package core

import "time"

/* 声明表结构 */
type VisitLog struct {
	ID         int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;type:timestamp" json:"created_at"`
	VisitTime  time.Time `gorm:"column:visit_time;not null;type:timestamp" json:"visit_time"`
	ClientTime time.Time `gorm:"column:client_time;type:timestamp" json:"client_time"`
	Path       string    `gorm:"column:path;not null;type:text" json:"path"`
	Method     string    `gorm:"column:method;type:text;default:GET" json:"method"`
	IP         string    `gorm:"column:ip;type:inet" json:"ip"`
	UserAgent  string    `gorm:"column:user_agent;type:text" json:"user_agent"`
	Country    string    `gorm:"column:country;type:text" json:"country"`
	Region     string    `gorm:"column:region;type:text" json:"region"`
	City       string    `gorm:"column:city;type:text" json:"city"`
	Referer    string    `gorm:"column:referer;type:text" json:"referer"`
	Status     int64     `gorm:"column:status;type:smallint;default:200" json:"status"`
}

type DailyArticleStat struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Date      time.Time `gorm:"column:date;type:date;not null" json:"date"`
	Path      string    `gorm:"column:path;type:text;not null" json:"path"`
	UV        int64     `gorm:"column:uv;not null;type:int" json:"uv"`
	PV        int64     `gorm:"column:pv;not null;type:int" json:"pv"`
	CreatedAt time.Time `gorm:"column:created_at;not null;type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;type:timestamp" json:"updated_at"`
}

/* 指定表名 */
func (v VisitLog) TableName() string {
	return "visit_logs"
}

func (d DailyArticleStat) TableName() string {
	return "daily_article_stats"
}
