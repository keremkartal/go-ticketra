package domain

import (
	"time"
)

// User struct'ı veritabanındaki 'users' tablosunu temsil eder.
type User struct {
	ID        uint      `gorm:"primaryKey"` 
	Email     string    `gorm:"uniqueIndex;not null"` 
	Password  string    `gorm:"not null"`  
	Role      string    `gorm:"default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}