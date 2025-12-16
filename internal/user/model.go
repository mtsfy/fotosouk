package user

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

	CreatedAt time.Time    `gorm:"autoCreateTime"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime"`
	DeletedAt sql.NullTime `gorm:"index"`
}
