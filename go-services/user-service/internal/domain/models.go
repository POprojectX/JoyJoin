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

type ProfessionalRole struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"` // "Photographer", "DJ", "Chef" или кто то другой там
}

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email string `gorm:"unique;not null"`
	FirstName string
	LastName string
	CreatedAt time.Time

	Participations []EventParticipant `gorm:"foreignKey:UserID"`
}

type Event struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string    `gorm:"not null"`
	Description string
	Date        time.Time
	Location    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	
	Participants []EventParticipant `gorm:"foreignKey:EventID"`
}

type EventParticipant struct {
	ID       uint      `gorm:"primaryKey"`
	UserID   uuid.UUID `gorm:"type:uuid;index;not null"`
	EventID  uuid.UUID `gorm:"type:uuid;index;not null"`
	
	// Системная роль Owner, Organizer, Staff и так далее
	SystemRole SystemRole `gorm:"not null"`
	
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