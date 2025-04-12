package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Content  string `json:"content"`
	UserID   uint   `json:"userId"`
	ChannelID uint  `json:"channelId"`
	User     User   `json:"user" gorm:"foreignKey:UserID"`
}