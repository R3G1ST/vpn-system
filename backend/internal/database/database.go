package database

import (
	"log"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"xferant-vpn/backend/internal/config"
)

func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := buildDSN(cfg)
	
	// Настройка логгера GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}
	
	// Получаем соединение с базой данных
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	// Настройка пула соединений
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	// Auto migrate моделей
	err = db.AutoMigrate(
		&User{},
		&VPNConfig{},
		&ApiKey{},
		&Payment{},
		&Server{},
	)
	if err != nil {
		return nil, err
	}
	
	log.Println("✅ Database connection established")
	return db, nil
}

func buildDSN(cfg config.DatabaseConfig) string {
	return "host=" + cfg.Host +
		" user=" + cfg.User +
		" password=" + cfg.Password +
		" dbname=" + cfg.DBName +
		" port=" + cfg.Port +
		" sslmode=" + cfg.SSLMode
}

// Модели базы данных
type User struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Username     string    `gorm:"not null"`
	PasswordHash string    `gorm:"not null"`
	Status       string    `gorm:"default:'active'"`
	TrafficLimit int64     `gorm:"default:1073741824"` // 1GB в байтах
	TrafficUsed  int64     `gorm:"default:0"`
	ExpireDate   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type VPNConfig struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `gorm:"not null;index"`
	Protocol  string    `gorm:"not null"` // vless, vmess, trojan
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
	Type      string    `gorm:"not null"` // read, write, admin
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	ExpiresAt time.Time
	LastUsed  time.Time
}

type Payment struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      string    `gorm:"not null;index"`
	Amount      float64   `gorm:"not null"`
	Currency    string    `gorm:"default:'RUB'"`
	Status      string    `gorm:"default:'pending'"` // pending, completed, failed
	PaymentID   string    `gorm:"uniqueIndex"`
	Provider    string    `gorm:"not null"` // yookassa, cloudpayments
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Server struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"not null"`
	Hostname    string    `gorm:"not null"`
	IPAddress   string    `gorm:"not null"`
	Status      string    `gorm:"default:'online'"` // online, offline, maintenance
	Location    string    `gorm:"not null"`
	LoadPercent int       `gorm:"default:0"`
	LastUpdate  time.Time
}