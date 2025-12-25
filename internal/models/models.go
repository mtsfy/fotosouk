package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID int `gorm:"primaryKey"`

	Username       string `gorm:"type:varchar(20);not null;unique"`
	FirstName      string `gorm:"type:varchar(255)"`
	LastName       string `gorm:"type:varchar(255)"`
	Email          string `gorm:"type:varchar(255);not null;unique"`
	HashedPassword string `gorm:"type:varchar(255);not null"`

	Images []Image `gorm:"foreignKey:UserID"`

	CreatedAt time.Time    `gorm:"autoCreateTime"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime"`
	DeletedAt sql.NullTime `gorm:"index"`
}

type Image struct {
	ID     int    `gorm:"primaryKey"`
	Url    string `gorm:"type:text;not null;"`
	UserID int    `gorm:"not null;index"`

	CreatedAt time.Time    `gorm:"autoCreateTime"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime"`
	DeletedAt sql.NullTime `gorm:"index"`
}

type RefreshToken struct {
	ID        int64     `gorm:"primaryKey"`
	UserID    int       `gorm:"not null;index"`
	Token     string    `gorm:"type:varchar(500);not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
	IsRevoked bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
