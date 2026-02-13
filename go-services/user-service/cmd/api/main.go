package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"user-service/internal/domain"
	"user-service/internal/storage"
)

func main() {
	config := &storage.ConfigDB{
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBname:   getEnv("DB_NAME", "joyjoin_db"),
		Port:     getEnv("DB_PORT", "5432"),
		SSLmode:  "disable",
	}

	db := storage.ConnectDB(config)
	if db == nil {
		log.Fatal("Failed to connect to database!")
	}

	err := db.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.EventParticipant{}, &domain.ProfessionalRole{})
	if err != nil {
		log.Fatal("Failed to automigrate!")
	}
	log.Println("Database connected and migrated successfully!")

	// Инициализация репозиториев
	//userRepo := repo.NewUserRepository(db)
	//eventRepo := repo.NewEventRepository(db)
	//participantRepo := repo.NewParticipantRepository(db)

	// Инициализация репозиториев
	// userRepo := repo.NewUserRepository(db, 10)
	// eventRepo := repo.NewEventRepository(db, 10)
	// participantRepo := repo.NewParticipantRepository(db)

	// // Инициализация сервисов
	// userService := services.NewUserService(userRepo)
	// eventService := services.NewEventService(eventRepo)
	// participantService := services.NewParticipantService(participantRepo, eventRepo, userRepo)
	// authService := services.NewAuthService(userRepo, getEnv("JWT_SECRET", "your-secret-key"))

	// // Оркестратор
	// orchestrator := services.NewEventOrchestrator(userService, eventService, participantService, authService)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	log.Println("🚀 Service running. Press Ctrl+C to stop.")
	<-sigChan
	
	log.Println("👋 Shutting down...")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}