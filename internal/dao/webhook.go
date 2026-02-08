package dao

import "gorm.io/gorm"

type GithubWebhookDao struct {
	db *gorm.DB
}

func NewGithubWebhookDao(db *gorm.DB) *GithubWebhookDao {
	return &GithubWebhookDao{db: db}
}
