package model

import "time"

type SubscribeUser struct {
	ID            int64     `gorm:"primaryKey;autoIncrement;column:id"`
	Telephone     string    `gorm:"column:telephone;not null;unique"`
	Email         string    `gorm:"column:email;type:text;not null"`
	VisitorID     string    `gorm:"column:visitor_id;type:text;not null"`
	SubscribeTime time.Time `gorm:"column:subscribe_time;type:datetime;not null"`
	Status        uint8     `gorm:"column:status;type:tinyint(1);not null"`
	UpdateAt      time.Time `gorm:"column:update_at;type:datetime"`
}

func (s SubscribeUser) TableName() string {
	return "subscribe_user"
}
