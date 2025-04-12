package models

import "gorm.io/gorm"

type Server struct {
	gorm.Model
	Name    string    `json:"name"`
	UserID  uint      `json:"userId"`
	Channels []Channel `json:"channels"`
}