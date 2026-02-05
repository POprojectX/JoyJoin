package services

import (
	"sync"
	"time"
)

//тут мы запускаем наш воркер пул
func StarterWorkerPool(workers int, taskQueue <- chan func()) {
	for i:= 0; i < workers; i++ {
		go func(id int) {
			for task := range taskQueue {
				task()
			}
		}(i)
	}
}

//тут будет лимитер который будет следить за лемитом на каждый email, что бы один еблан не положил сервак.
//Если будем писать тесты: один пользователь (то есть с одного мыла) может делать ОДИН запрос В СЕКУНДУ!
func CheckRateLimitPerEmail(key string, rateLimiter *sync.RWMutex, rateMap map[string]time.Time) bool {
	rateLimiter.Lock()
	defer rateLimiter.Unlock()

	now := time.Now()
	if lastTime, exists := rateMap[key]; exists {
		if now.Sub(lastTime) < time.Second {
			return false
		}
	}
	rateMap[key] = now
	return true
}

//утилиты для Update метода
func ApplyNotNil [T any] (target *T, source *T) bool {
	if source != nil {
		*target = *source
		return false
	}
	return true
}

func CollectUpdates(updates map[string]interface{}, condition bool, key string, value interface{}) {
	if condition {
		updates[key] = value
	}
}