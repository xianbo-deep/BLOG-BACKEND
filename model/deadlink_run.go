package model

import "time"

type DeadLinkRun struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;column:id"`
	StartedAt    time.Time `gorm:"column:started_at;type:timestamp;not null"`
	FinishedAt   time.Time `gorm:"column:finished_at;type:timestamp;not null"`
	PagesScanned int       `gorm:"column:pages_scanned;type:int"`
	DeadLinkCnt  int       `gorm:"column:deadlink_cnt;type:int"`
	LinksChecked int       `gorm:"column:links_checked;type:int"`
}

func (DeadLinkRun) TableName() string {
	return "deadlink_run"
}
