package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	FullName  string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"` // stored as bcrypt hash
	IsAdmin   bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Movie struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string `gorm:"type:text"`
	DurationMin int    // duration in minutes
	ReleaseDate time.Time
	PosterURL   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Shows       []Show
}

type Show struct {
	ID          uint `gorm:"primaryKey"`
	MovieID     uint `gorm:"index"`
	Movie       Movie
	Hall        string
	StartTime   time.Time
	SeatsTotal  int
	SeatsBooked int
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Booking struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint `gorm:"index"`
	User       User
	ShowID     uint `gorm:"index"`
	Show       Show
	SeatsCount int
	TotalPrice float64
	Status     string // e.g., "pending", "confirmed", "cancelled"
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
