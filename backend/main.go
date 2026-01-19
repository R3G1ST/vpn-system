package main

import (
	"log"
	"xferant-vpn/backend/internal/api"
	"xferant-vpn/backend/internal/config"
	"xferant-vpn/backend/internal/database"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()
	
	log.Printf("🚀 Starting Xferant VPN Backend v1.0.0")
	log.Printf("🌐 Server address: %s", cfg.Server.Address)
	log.Printf("🗄️ Database: %s@%s:%s/%s", 
		cfg.Database.User, 
		cfg.Database.Host, 
		cfg.Database.Port, 
		cfg.Database.DBName)
	
	// Подключение к БД
	db, err := database.Init(cfg.Database)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	
	// Инициализация сервера
	server := api.NewServer(cfg, db)
	
	log.Printf("✅ Server starting on %s", cfg.Server.Address)
	if err := server.Start(); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}