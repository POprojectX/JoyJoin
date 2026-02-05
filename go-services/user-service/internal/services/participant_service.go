package services

import (
	"sync"
	"time"
	"user-service/internal/repo"
)

type ParticipantService interface {

}

type participantService struct {
	participantRepo repo.ParticipantRepository
	rateLimiter sync.RWMutex
	rateMap map[string]time.Time
	taskQueue chan func()
}

func NewParticipantService(r repo.ParticipantRepository) ParticipantService {
	s := &participantService{
		participantRepo: r,
		rateMap:   make(map[string]time.Time),
		taskQueue: make(chan func(), 1000), 
	}

	s.starterWorkerPool(10)

	return s
}

//используем утилиты
func (s *participantService) starterWorkerPool(workers int) {
	StarterWorkerPool(workers, s.taskQueue)
}

func (s *participantService) checkRateLimitPerEmail(email string) bool {
	return CheckRateLimitPerEmail(email, &s.rateLimiter, s.rateMap)
}