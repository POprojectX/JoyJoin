package main

import (
	"context"
	"log"
	"os"
	"user-service/internal/domain"
	httpdelivery "user-service/internal/http"
	"user-service/internal/repo"
	"user-service/internal/services"
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

	//Инициализация репозиториев
	userRepo := repo.NewUserRepository(db, 10)
	eventRepo := repo.NewEventRepository(db, 10)
	participantRepo := repo.NewParticipantRepository(db)

	// Инициализация сервисов
	userService := services.NewUserService(userRepo)
	eventService := services.NewEventService(eventRepo)
	participantService := services.NewParticipantService(participantRepo, eventRepo, userRepo)
	authService := services.NewAuthService(userRepo, getEnv("JWT_SECRET", "your-secret-key"))

	// Оркестратор
	orchestrator := services.NewEventOrchestrator(userService, eventService, participantService, authService)

	ctx := context.Background()

	//указатель на authService нужен потому что у нас все сервисвы сделаны через интерфейсы, которые сами по себе уже ссылочный тип,
	//а authService работает через структуру поэтому нужен указатель на структуру
	h := httpdelivery.NewHandler(userService, eventService, participantService, *authService, orchestrator, ctx)

	httpdelivery.Run(h, ":8080")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}