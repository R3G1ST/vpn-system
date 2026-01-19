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
	
	log.Printf("🔗 Connecting to database: %s@%s:%s/%s", 
		cfg.User, cfg.Host, cfg.Port, cfg.DBName)
	
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}
	
	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	
	log.Println("✅ Database connection established successfully")
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