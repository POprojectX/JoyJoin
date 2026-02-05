package services

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailExists      = errors.New("email already exist")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrValidationFailed = errors.New("validation failed")
	ErrShortPassword    = errors.New("password is too short, must be more then 5 symbols!")
	ErrTooManyRequests  = errors.New("too many requests! (1 req/sec)")
	ErrContextCancelled = errors.New("context was cancelled")
	ErrNoEvents         = errors.New("user has no events")
	ErrNoEventTitle 	= errors.New("set the title for an event")
	ErrNoOwnerID 		= errors.New("event dont have OwnerID")
	ErrEventNotFound 	= errors.New("event does not exist")
	ErrEventNoParticipants = errors.New("event has no participants")
	ErrNoEventsWithRole = errors.New("no events with this role")
)