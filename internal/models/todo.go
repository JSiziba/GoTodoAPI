package models

import (
	"time"

	"gorm.io/gorm"
)

// Todo represents a task in the todo list
type Todo struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Completed   bool           `json:"completed" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TodoCreate represents the structure used when creating a new Todo
type TodoCreate struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// TodoUpdate represents the structure used when updating a Todo
type TodoUpdate struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}
