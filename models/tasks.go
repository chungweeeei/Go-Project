package models

import (
	"time"
)

// gorm.Model definition
type Task struct {
	ID          uint      `json:"id";gorm:"primaryKey;not null`
	UserEmail   string    `json:"user_email";gorm:not null;index`
	Title       string    `json:"title";gorm:"type:varchar(200);not null"`
	Description string    `json:"description";gorm:"type:text;not null"`
	Status      string    `json:"status";gorm:"type:varchar(20);not null;default:'todo'"`
	Priority    string    `json:"priority";gorm:"type:varchar(10);not null;default:'medium'"`
	DueDate     time.Time `json:"due_date";form:"type:date"`
	Category    string    `json:"category";gorm:"type:varchar(50)"`
	CreatedAt   time.Time `json:"created_at";gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at";gorm:"autoUpdateTime"`
}
