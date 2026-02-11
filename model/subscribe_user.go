package model

import "time"

type SubscribeUser struct {
	ID            int64     `gorm:"primaryKey;autoIncrement;column:id"`
	Email         string    `gorm:"column:email;type:text;not null;unique"`
	SubscribeTime time.Time `gorm:"column:subscribe_time;type:timestamp;not null"`
	Status        int       `gorm:"column:status;type:tinyint(1);not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:timestamp"`
	LastSentAt    time.Time `gorm:"column:last_sent_at;type:timestamp"`
	NotifyCount   int       `gorm:"column:notify_count;type:int;not null"`
}

func (s SubscribeUser) TableName() string {
	return "subscribe_user"
}
