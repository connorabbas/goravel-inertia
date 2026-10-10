package models

import "time"

type PasswordResetToken struct {
	Email     string `gorm:"primaryKey"`
	TokenHash string
	ExpiresAt time.Time
}
