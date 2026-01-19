package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Username     string    `gorm:"not null"`
	PasswordHash string    `gorm:"not null"`
	Status       string    `gorm:"default:'active'"`
	TrafficLimit int64     `gorm:"default:1073741824"`
	TrafficUsed  int64     `gorm:"default:0"`
	ExpireDate   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type VPNConfig struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `gorm:"not null;index"`
	Protocol  string    `gorm:"not null"`
	Config    string    `gorm:"type:text;not null"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ApiKey struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `gorm:"not null"`
	Key       string    `gorm:"uniqueIndex;not null"`
	Secret    string    `gorm:"not null"`
	Type      string    `gorm:"not null"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Server struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"not null"`
	Hostname    string    `gorm:"not null"`
	IPAddress   string    `gorm:"not null"`
	Status      string    `gorm:"default:'online'"`
	Location    string    `gorm:"not null"`
	LoadPercent int       `gorm:"default:0"`
	LastUpdate  time.Time
}