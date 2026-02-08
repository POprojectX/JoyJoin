package domain

import (
	"time"

	"github.com/google/uuid"
)

//строгая типизация ролей, потом можно будет добавить еще какие то роли
type SystemRole string

const (
	RoleOwner     SystemRole = "Owner"
	RoleOrganizer SystemRole = "Organizer"
	RoleManager   SystemRole = "Manager"
	RoleGuest     SystemRole = "Guest"
	RoleStaff     SystemRole = "Staff"
)

// тут все статусы ивента + обьяснения для чего они нужны
type Status string

const (
	// черновик, виден ивент только тем кто работает над ним, например если еще ивент еще на садии планирования, логистики и тд. 
	// Можно редачить, регаться нет
	StatusDraft Status = "Draft"   
	StatusAnnounced Status = "Announced" //Анонсирован, билетов/регистрации пока нет, просто как новость об ивенте.
	StatusOngoing Status = "Ongoing" // ивент идет прямо сейчас, редачить нельзя место/дату, регаться уже нельзя
	StatusCompleted Status = "Completed" // уже завершен, присойдениться нельзя
	StatusCancelled Status = "Cancelled" // отменен, все активные роли аннулируются, отменить можно организатором
)

type ProfessionalRole struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"` // "Photographer", "DJ", "Chef" или кто то другой там
	CreatedAt time.Time `gorm:"index"`
}

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email string `gorm:"uniqueIndex;not null;size:255"`
	FirstName string    `gorm:"size:100"`
	LastName  string    `gorm:"size:100"`
	Password  []byte    `gorm:"not null"` // храним как bytes так как сам хеш у нас в байтах
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	
	EventsCount int `gorm:"-"`

	Participations []EventParticipant `gorm:"foreignKey:UserID"`
}

//решить надо с update методом! 
type UpdateUserInput struct {
	FirstName *string
	LastName *string
	Passwrod *string
}

type Event struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Status 		Status    `gorm:"not null"`
	Title       string    `gorm:"not null;size:200;index"`
	Description string    `gorm:"type:text"`
	Slots		int 	  `gorm:"not null"`
	AvailableSlots int	  `gorm:"not null;check:available_slots <= slots""`
	DateFrom    time.Time `gorm:"index"`
	DateTo		time.Time `gorm:"index"`
	Location    string    `gorm:"size:255"`
	OwnerID     uuid.UUID `gorm:"type:uuid;index;not null"`
	CreatedAt   time.Time `gorm:"index"`
	UpdatedAt   time.Time

	Participants []EventParticipant `gorm:"foreignKey:EventID"`
}

//DTO чисто для метода Update, что бы сделать его более еластичным (возможность обновлять еластично поля)
type UpdateEventInput struct {
	Status		*Status
	Title       *string
	Description *string
	Location    *string
	DateFrom    *time.Time
	DateTo 		*time.Time
	Slots		*int
}


type EventParticipant struct {
	ID       uint      `gorm:"primaryKey"`
	UserID   uuid.UUID  `gorm:"type:uuid;index:idx_user_event,unique;not null"`
	EventID  uuid.UUID  `gorm:"type:uuid;index:idx_user_event,unique;not null"`
	
	// Системная роль Owner, Organizer, Staff и так далее
	SystemRole         SystemRole `gorm:"index;not null;size:20"`
	
	// Специализация (только для Staff, может быть NULL)
	ProfessionalRoleID *uint             `gorm:"index"`
	ProfessionalRole   *ProfessionalRole `gorm:"foreignKey:ProfessionalRoleID"`
	
	// Доп поля
	JoinedAt time.Time `gorm:"autoCreateTime"`
	Notes    string    // заметки организатора о участнике, например если у кого то будет синдром дауна можно это пометить
	
	// Связи с preload чтобы подгружать связанные данные
	User   User   `gorm:"foreignKey:UserID"`
	Event  Event  `gorm:"foreignKey:EventID"`
}

// DTO
type UserWithEventsDTO struct {
	ID        uuid.UUID          `json:"id"`
	Email     string             `json:"email"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Events    []EventSummaryDTO  `json:"events"` // упрощённая структура, не вся тяжёлая модель
}
type EventSummaryDTO struct {
	EventID    uuid.UUID       `json:"event_id"`
	Title      string          `json:"title"`
	DateFrom   time.Time       `json:"date"`
	DateTo   time.Time       `json:"date"`
	SystemRole SystemRole      `json:"system_role"`
	JoinedAt   time.Time       `json:"joined_at"`
}