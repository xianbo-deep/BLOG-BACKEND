package model

import "time"

/* 声明表结构 */
type VisitLog struct {
	ID          int64     `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;type:timestampz"`
	VisitTime   time.Time `gorm:"column:visit_time;not null;type:timestampz"`
	ClientTime  time.Time `gorm:"column:client_time;type:timestamp"`
	Path        string    `gorm:"column:path;not null;type:text"`
	Method      string    `gorm:"column:method;type:text;default:GET"`
	IP          string    `gorm:"column:ip;type:text"`
	UserAgent   string    `gorm:"column:user_agent;type:text"`
	Country     string    `gorm:"column:country;type:text"`
	CountryCode string    `gorm:"column:country_code;type:text"`
	CountryEN   string    `gorm:"column:country_en;type:text"`
	Region      string    `gorm:"column:region;type:text"`
	RegionCode  string    `gorm:"column:region_code;type:text"`
	RegionEN    string    `gorm:"column:region_en;type:text"`
	City        string    `gorm:"column:city;type:text"`
	CityEN      string    `gorm:"column:city_en;type:text"`
	Referer     string    `gorm:"column:referer;type:text"`
	Status      int64     `gorm:"column:status;type:smallint;default:200"`
	VisitorID   string    `gorm:"column:visitor_id;type:text"`
	Latency     int64     `gorm:"column:latency;type:int;default:0"`
	Medium      string    `gorm:"column:refr_medium;type:text"`
	Source      string    `gorm:"column:refr_source;type:text"`
	Device      string    `gorm:"column:device;type:text"`
	OS          string    `gorm:"column:os;type:text"`
	Browser     string    `gorm:"column:browser;type:text"`
	Lat         float64   `gorm:"column:lat;type:double precision;default:999"`
	Lon         float64   `gorm:"column:lon;type:double precision;default:999"`
}

/* 指定表名 */
func (v VisitLog) TableName() string {
	return "visit_logs"
}
