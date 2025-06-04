package models

import (
	"time"
)

type Task struct {
	ID          uint      `gorm:"primaryKey;not null`
	UserEmail   string    `gorm:not null;index`
	Title       string    `gorm:"type:varchar(200);not null"`
	Description string    `gorm:"type:text;not null"`
	Status      string    `gorm:"type:varchar(20);not null;default:'todo'"`
	Priority    string    `gorm:"type:varchar(10);not null;default:'medium'"`
	DueDate     time.Time `form:"type:date"`
	Category    string    `gorm:"type:varchar(50)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
