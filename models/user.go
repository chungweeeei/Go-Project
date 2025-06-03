package models

import (
	"time"
)

// gorm.Model definition
type User struct {
	Email     string    `gorm:"primaryKey;index;not null"`
	Username  string    `gorm:"not null"`
	Role      string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
