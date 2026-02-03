package main

import (
	"log"
	"os"
	"user-service/internal/domain"
	"user-service/internal/storage"
)

func main() {
	config := &storage.ConfigDB{
		Host:     os.Getenv("DB_HOST"),
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        DBname:   os.Getenv("DB_NAME"),
        Port:     os.Getenv("DB_PORT"),
        SSLmode:  "disable",
        TimeZone: "Europe/Opole",
	}

	db := storage.ConnectDB(config)
	if db == nil {
		log.Fatal("Failed to connect to database!")
	}

	err := db.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.EventParticipant{})
	if err != nil {
		log.Fatal("Failed to automigrate!")
	}
	log.Println("Database connected and migrated successfully!")

	// Инициализация репозиториев
	//userRepo := repo.NewUserRepository(db)
	//eventRepo := repo.NewEventRepository(db)
	//participantRepo := repo.NewParticipantRepository(db)


}