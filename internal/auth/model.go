package auth

import "time"

type RefreshToken struct {
	ID        int64     `gorm:"primaryKey"`
	UserID    int       `gorm:"not null;index"`
	Token     string    `gorm:"type:varchar(500);not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
	IsRevoked bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
