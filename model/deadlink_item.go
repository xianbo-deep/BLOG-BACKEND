package model

import "time"

type DeadLinkItem struct {
	ID    int64 `gorm:"primaryKey;autoIncrement;column:id"`
	RunID int64 `gorm:"column:run_id;type:bigint;not null;index"`

	FromPage   string    `gorm:"column:from_page;type:text;not null"`
	LinkURL    string    `gorm:"column:link_url;type:text;not null"`
	StatusCode int       `gorm:"column:status_code;type:int"`
	OK         bool      `gorm:"column:ok;type:tinyint(1);not null"`
	Err        string    `gorm:"column:err;type:text"`
	CheckedAt  time.Time `gorm:"column:checked_at;type:datetime;not null"`
}

func (DeadLinkItem) TableName() string {
	return "deadlink_item"
}
