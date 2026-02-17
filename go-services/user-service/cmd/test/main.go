package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repo"
	"user-service/internal/services"
	"user-service/internal/storage"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func test(name string, err error) {
	if err != nil {
		log.Printf("FAIL: %s - %v", name, err)
	} else {
		log.Printf("PASS: %s", name)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute) //ставим на 10 минут потому что у нас один контекст на все тесты
	defer cancel()

	config := &storage.ConfigDB{
		Host:     getEnv("DB_HOST", "test-db"),
		User:     getEnv("DB_USER", "test"),
		Password: getEnv("DB_PASSWORD", "test"),
		DBname:   getEnv("DB_NAME", "joyjoin_test_db"),
		Port:     getEnv("DB_PORT", "5432"),
		SSLmode:  "disable",
	}

	var db *gorm.DB
	var err error
	
	for i := 0; i < 10; i++ {
		db = storage.ConnectDB(config)
		if db != nil {
			sqlDB, err := db.DB()
			if err == nil && sqlDB.Ping() == nil {
				break
			}
		}
		log.Printf("Waiting for database... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}
	if db == nil {
		log.Fatal("Failed to connect to database after multiple attempts!")
	}

	db.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.EventParticipant{}, &domain.ProfessionalRole{})
	log.Println("Database connected and migrated successfully!")

	userRepo := repo.NewUserRepository(db, 10)
	eventRepo := repo.NewEventRepository(db, 10)
	participantRepo := repo.NewParticipantRepository(db)

	userService := services.NewUserService(userRepo)
	eventService := services.NewEventService(eventRepo)
	participantService := services.NewParticipantService(participantRepo, eventRepo, userRepo)
	authService := services.NewAuthService(userRepo, "test-secret")
	orchestrator := services.NewEventOrchestrator(userService, eventService, participantService, authService)

	log.Printf("=======================start testing====================")

	// ============ USER SERVICE =====================================================================================================
	log.Println("--- User Service ---")
	// 1. Регистрация
	email1 := fmt.Sprintf("user1_%d@test.com", time.Now().Unix())
	user1, _, err := orchestrator.RegisterAndCreateProfile(ctx, email1, "password123", "Bob", "McLover")
	test("User registration", err)

	// 2. Повторная регистрация (должна упасть)
	_, _, err = orchestrator.RegisterAndCreateProfile(ctx, email1, "password123", "Bob", "McLover")
	if err != nil {
		test("Blocking already used email ✅", nil)
	} else {
		test("Blocking already used email", fmt.Errorf("Should have failed"))
	}

	// 3. Логин
	time.Sleep(time.Second * 2)
	_, token, err := orchestrator.LoginAndGetProfile(ctx, email1, "password123")
	test("Login", err)
	if err == nil && token != "" {
		preview := token
		if len(token) > 20 {
			preview = token[:20]
		}
		log.Printf("   Token: %s...", preview)
	}

	// 4. Логин с неправильным паролем
	_, _, err = orchestrator.LoginAndGetProfile(ctx, email1, "wrongpass")
	if err != nil {
		test("Blocking wrong password", nil)
	} else {
		test("Blocking wrong password", fmt.Errorf("Should have failed"))
	}

	// 5. Получение профиля
	profile, err := orchestrator.GetUserFullProfile(ctx, user1.ID)
	test("Getting user profile", err)
	log.Printf("   User: %s %s", profile.FirstName, profile.LastName)

	//6. Обновление профиля
	newName := "Mike"
	_, err = userService.Update(ctx, domain.UpdateUserInput{FirstName: &newName}, user1.ID)
	test("Updating user profile", err)

	//7. ПОлучение пользователя с событиями
	userWithEvents, err := userService.GetUserWithEvents(ctx, user1.ID)
	test("Getting user with events", err)
	log.Printf("   User with events: %d events", len(userWithEvents.Events))

	// ============ EVENT SERVICE ========================================================================================================
	log.Println("--- Event Service ---")

	// Создаём второго пользователя для тестов
	email2 := fmt.Sprintf("user2_%d@test.com", time.Now().Unix())
	user2, _, err := orchestrator.RegisterAndCreateProfile(ctx, email2, "password123", "Tester", "OfTests")
	if err != nil || user2 == nil {
		log.Fatalf("Failed to create user2: %v", err)
	}

	// 8. Создание события
	event, err := orchestrator.CreateEventWithOwner(
		ctx, "Party", "Cool Party", "Opole",
		10, time.Now().Add(24*time.Hour), time.Now().Add(48*time.Hour),
		user2.ID, domain.AccessPublic,
	)
	test("Create event", err)
	log.Printf("   Event ID: %s, Status: %s", event.ID, event.Status)

	// 9. Проверка статуса Draft
	if event.Status == domain.StatusDraft {
		test("Draft Status by default ✅", nil)
	} else {
		test("Draft Status by default", fmt.Errorf("expected Draft, got %s", event.Status))
	}
	time.Sleep(2 * time.Second)

	// 10. Публикация события
	published, err := orchestrator.PublishEventAtomic(ctx, event.ID, user2.ID)
	test("Publishing event", err)
	if published.Status == domain.StatusAnnounced {
		test("Status Announced after publishing", nil)
	} else {
		test("Status Announced after publishing", fmt.Errorf("expected Announced, got %s", published.Status))
	}

	// 11. Попытка публикации не-владельцем
	_, err = orchestrator.PublishEventAtomic(ctx, event.ID, user1.ID)
	if err != nil {
		test("Blocking publishing event by non-owner", nil)
	} else {
		test("Blocking publishing event by non-owner", fmt.Errorf("should have failed"))
	}

	// 12. Обновление события
	newTitle := "Super Party PDiddy Edition"
	_, err = eventService.Update(ctx, domain.UpdateEventInput{Title: &newTitle}, event.ID)
	test("Updating event title", err)

	// 13. Получение участников (пока только owner)
	participants, err := eventService.GetEventParticipants(ctx, event.ID)
	test("Getting event participants (only owner)", err)
	log.Printf("   Participants: %d", len(participants))

	// ============ SLOTS & PARTICIPANTS ============
	log.Println("\n--- Slots & Participants ---")
	time.Sleep(2 * time.Second)

	// 14. Присоединение как гость (user1)
	err = orchestrator.JoinEventAsGuestPublic(ctx, user1.ID, event.ID)
	test("Joining event as guest", err)

	time.Sleep(2 * time.Second)
	// 15. Проверка слотов
	eventWithSlot, _ := eventService.GetByID(ctx, event.ID)
	log.Printf("   Total Slots: %d, Available: %d", eventWithSlot.Slots, eventWithSlot.AvailableSlots)
	if eventWithSlot.AvailableSlots == 9 {
		test("Slot occupied (9 left)", nil)
	} else {
		test("Slot occupied", fmt.Errorf("expected 9 available, got %d", eventWithSlot.AvailableSlots))
	}

	time.Sleep(2 * time.Second)
	// 16. Повторное присоединение (должно упасть)
	err = orchestrator.JoinEventAsGuestPublic(ctx, user1.ID, event.ID)
	if err != nil {
		test("Blocking repeated joining ✅", nil)
	} else {
		test("Blocking repeated joining", fmt.Errorf("should have failed"))
	}

	time.Sleep(2 * time.Second)
	// 17. Создаём 10 гостей и занимаем все слоты
	log.Println("   Creating 10 guests to fill all slots...")
	guests := make([]uuid.UUID, 10)
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		email := fmt.Sprintf("guest%d_test17@test.com", i)
		g, _, errReg := orchestrator.RegisterAndCreateProfile(ctx, email, "password123123", "Guest", fmt.Sprintf("№%d", i))
		if errReg != nil || g == nil {
			log.Printf("   Warning: failed to create guest %d: %v", i, errReg)
			continue
		}
		guests[i] = g.ID
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)
	// Занимаем оставшиеся 9 слотов
	for i := 0; i < 9; i++ {
		orchestrator.JoinEventAsGuestPublic(ctx, guests[i], event.ID)
	}

	time.Sleep(2 * time.Second)
	// 18. Проверяем что слоты кончились
	fullEvent, _ := eventService.GetByID(ctx, event.ID)
	if fullEvent.AvailableSlots == 0 {
		test("All slots are occupied (0 available) ✅", nil)
	} else {
		test("All slots are occupied", fmt.Errorf("expected 0, got %d", fullEvent.AvailableSlots))
	}

	time.Sleep(2 * time.Second)
	// 19. Попытка занять слот когда их нет
	err = orchestrator.JoinEventAsGuestPublic(ctx, guests[9], event.ID)
	if err != nil {
		test("Blocking joining when no slots available ✅", nil)
	} else {
		test("Blocking joining when no slots available", fmt.Errorf("should have failed"))
	}

	time.Sleep(2 * time.Second)
	// 20. Освобождение слота (guest покидает событие)
	err = orchestrator.LeaveEventAndFreeSlot(ctx, guests[0], event.ID)
	test("Leaving event + freeing slot", err)

	afterLeave, _ := eventService.GetByID(ctx, event.ID)
	if afterLeave.AvailableSlots == 1 {
		test("Slot freed (1 available) ✅", nil)
	} else {
		test("Slot freed", fmt.Errorf("expected 1, got %d", afterLeave.AvailableSlots))
	}

	// ============ ROLES =============================================================================================================
	log.Println("\n--- Role Management ---")

	time.Sleep(2 * time.Second)
	// 21. Назначение Staff (owner назначает)
	staffEmail := fmt.Sprintf("staff_%d@test.com", time.Now().Unix())
	staff, _, err := orchestrator.RegisterAndCreateProfile(ctx, staffEmail, "password123123", "Employee", "Test")
	if err != nil || staff == nil {
		log.Fatalf("Failed to create staff: %v", err)
	}
	
	err = orchestrator.AssignStaffWithOutSlotCheck(ctx, user2.ID, staff.ID, event.ID, nil)
	test("Staff assigned", err)

	time.Sleep(2 * time.Second)
	// 22. Проверка роли
	var role domain.SystemRole
	for i := 0; i < 5; i++ {
		role, err = participantService.GetUserRoleInEvent(ctx, staff.ID, event.ID)
		if err == nil && role == domain.RoleStaff {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if role == domain.RoleStaff {
		test("Staff role confirmed ✅", nil)
	} else {
		test("Staff role not confirmed", fmt.Errorf("expected Staff, got %s", role))
	}

	time.Sleep(2 * time.Second)
	// 23. Смена роли (Guest -> Organizer)
	// Сначала добавляем гостя, потом меняем роль
	orgEmail := "org_TEST23@test.com"
	org, _, err := orchestrator.RegisterAndCreateProfile(ctx, orgEmail, "password123123", "Organizer", "Junior")
	if err != nil || org == nil {
		log.Fatalf("Failed to create org: %v", err)
	}
	err = orchestrator.JoinEventAsGuestPublic(ctx, org.ID, event.ID)
	time.Sleep(500 * time.Millisecond)
	if err != nil {
		log.Printf("   Warning: org failed to join: %v", err)
	}else {
		var participant *domain.EventParticipant
		for i := 0; i < 5; i++ {
			participant, err = participantService.GetParticipantByUserID(ctx, org.ID, event.ID)
			if err == nil && participant != nil {
				break
			}
			time.Sleep(500 * time.Microsecond)
		}
		if err != nil || participant == nil {
			log.Printf("   Warning: org not found as participant: %v", err)
		} else {
			err = participantService.ChangeRole(ctx, user2.ID, participant.ID, domain.RoleOrganizer)
			test("Changing role to Organizer", err)
		}
	}

	time.Sleep(2 * time.Second)
	// 24. Проверка что Organizer может назначать Staff
	staff2Email := "staff2_24@test.com"
	staff2, _, _ := orchestrator.RegisterAndCreateProfile(ctx, staff2Email, "password123123", "Employee2", "Test")
	
	err = orchestrator.AssignStaffWithOutSlotCheck(ctx, org.ID, staff2.ID, event.ID, nil)
	test("Organizer assigns Staff", err)

	time.Sleep(2 * time.Second)
	// 25. Проверка что Organizer НЕ может менять роль Owner
	ownerParticipant, _ := participantService.GetParticipantByUserID(ctx, user2.ID, event.ID)
	err = participantService.ChangeRole(ctx, org.ID, ownerParticipant.ID, domain.RoleGuest)
	if err != nil {
		test("Organizer cant change Owner role ✅", nil)
	} else {
		test("Organizer cant change Owner role", fmt.Errorf("should have failed"))
	}

	// ============ PRIVATE EVENTS =====================================================================================================
	log.Println("\n--- Private Events ---")

	time.Sleep(2 * time.Second)
	// 26. Создание приватного события
	privateEvent, _ := orchestrator.CreateEventWithOwner(
		ctx, "Private Party", "Only for hood", "Opole",
		5, time.Now().Add(24*time.Hour), time.Now().Add(48*time.Hour),
		user2.ID, domain.AccessPrivate,
	)
	orchestrator.PublishEventAtomic(ctx, privateEvent.ID, user2.ID)

	time.Sleep(2 * time.Second)
	// 27. Попытка публичного присоединения к приватному
	randomUserEmail := fmt.Sprintf("random_%d@test.com", time.Now().Unix())
	randomUser, _, _ := orchestrator.RegisterAndCreateProfile(ctx, randomUserEmail, "password123123", "Random", "User")
	
	err = orchestrator.JoinEventAsGuestPublic(ctx, randomUser.ID, privateEvent.ID)
	if err != nil {
		test("Blocking public access to private event", nil)
	} else {
		test("Blocking public access to private event", fmt.Errorf("should have failed"))
	}

	time.Sleep(2 * time.Second)
	// 28. Приватное присоединение по приглашению (через owner)
	err = orchestrator.JoinEventAsGuestPrivate(ctx, user2.ID, randomUser.ID, privateEvent.ID)
	test("Private joining via invitation", err)

	// ============ CANCELLATION & DELETION =================================================================================
	log.Println("\n--- Cancellation & Deletion ---")

	time.Sleep(2 * time.Second)
	// 29. Отмена события
	err = orchestrator.CancelEventWithCleanup(ctx, privateEvent.ID, user2.ID)
	test("Canceling event", err)

	cancelled, _ := eventService.GetByID(ctx, privateEvent.ID)
	if cancelled.Status == domain.StatusCancelled {
		test("Status Cancelled ✅", nil)
	} else {
		test("Status Cancelled", fmt.Errorf("expected Cancelled, got %s", cancelled.Status))
	}

	time.Sleep(2 * time.Second)
	// 30. Удаление пользователя с очисткой
	err = orchestrator.DeleteUserAndCleanup(ctx, randomUser.ID)
	test("Deleting user with cleanup", err)

	// Проверяем что пользователь удалён
	_, err = userService.GetByID(ctx, randomUser.ID)
	if err != nil {
		test("User deleted successfully", nil)
	} else {
		test("User deleted successfully", fmt.Errorf("user still exists"))
	}

	time.Sleep(2 * time.Second)

	// ============ КОНКУРЕНТНЫЙ ТЕСТ ==================================================================================================
	log.Println("\n--- Concurrent Slots Test ---")

	// Создаём событие с 5 слотами
	concurrentEvent, _ := orchestrator.CreateEventWithOwner(
		ctx, "Concurrent Event", "Test", "Opole",
		5, time.Now().Add(24*time.Hour), time.Now().Add(48*time.Hour),
		user2.ID, domain.AccessPublic,
	)
	orchestrator.PublishEventAtomic(ctx, concurrentEvent.ID, user2.ID)
	time.Sleep(2 * time.Second)

	// Создаём 10 гостей
	concurrentGuests := make([]uuid.UUID, 10)
	for i := 0; i < 10; i++ {
		email := fmt.Sprintf("concurrent%d_30@test.com", i)
		g, _, errReg := orchestrator.RegisterAndCreateProfile(ctx, email, "password123123", "Concurrent", fmt.Sprintf("№%d", i))
		if errReg != nil {
			log.Printf("   Warning: failed to create concurrent guest %d: %v", i, errReg)
		}
		concurrentGuests[i] = g.ID
	}

	// Одновременно пытаемся занять слоты
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		if concurrentGuests[i] == uuid.Nil {
			continue
		}
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			err := orchestrator.JoinEventAsGuestPublic(ctx, concurrentGuests[idx], concurrentEvent.ID)
			mu.Lock()
			if err == nil {
				successCount++
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	if successCount == 5 {
		test(fmt.Sprintf("Concurrent test: %d successes, 5 failures ✅", successCount), nil)
	} else {
		test("Concurrent test", fmt.Errorf("expected 5 successes, got %d", successCount))
	}

	finalEvent, _ := eventService.GetByID(ctx, concurrentEvent.ID)
	log.Printf("   Total: slots %d, available %d", finalEvent.Slots, finalEvent.AvailableSlots)
	
	// ============ END TESTS ============
	log.Println("\n=== 📊 All tests completed ===")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}