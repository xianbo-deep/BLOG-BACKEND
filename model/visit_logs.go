package model

import "time"

/* 声明表结构 */
type VisitLog struct {
	ID         int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;type:timestamp" json:"created_at"`
	VisitTime  time.Time `gorm:"column:visit_time;not null;type:timestamp" json:"visit_time"`
	ClientTime time.Time `gorm:"column:client_time;type:timestamp" json:"client_time"`
	Path       string    `gorm:"column:path;not null;type:text" json:"path"`
	Method     string    `gorm:"column:method;type:text;default:GET" json:"method"`
	IP         string    `gorm:"column:ip;type:text" json:"ip"`
	UserAgent  string    `gorm:"column:user_agent;type:text" json:"user_agent"`
	Country    string    `gorm:"column:country;type:text" json:"country"`
	Region     string    `gorm:"column:region;type:text" json:"region"`
	City       string    `gorm:"column:city;type:text" json:"city"`
	Referer    string    `gorm:"column:referer;type:text" json:"referer"`
	Status     int64     `gorm:"column:status;type:smallint;default:200" json:"status"`
	VisitorID  string    `gorm:"column:visitor_id;type:text" json:"visitor_id"`
	Latency    int64     `gorm:"column:latency;type:int;default:0" json:"latency"`
	Medium     string    `gorm:"column:refr_medium;type:text" json:"medium"`
	Source     string    `gorm:"column:refr_source;type:text" json:"source"`
	Device     string    `gorm:"column:device;type:text" json:"device"`
	OS         string    `gorm:"column:os;type:text" json:"os"`
	Browser    string    `gorm:"column:browser;type:text" json:"browser"`
}

/* 指定表名 */
func (v VisitLog) TableName() string {
	return "visit_logs"
}
