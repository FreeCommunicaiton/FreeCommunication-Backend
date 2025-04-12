package models

import "gorm.io/gorm"

type Channel struct {
	gorm.Model
	Name     string    `json:"name"`
	ServerID uint      `json:"serverId"`
	Messages []Message `json:"messages"`
}