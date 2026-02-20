package httpdelivery

import (
	"context"
	"user-service/internal/services"
)

type Handler struct {
	UserService services.UserService
	EvenService services.EventService
	ParticipantService services.ParticipantService
	AuthService services.AuthService
	Orchestrator services.OrchestratorService
	Ctx context.Context
}

func NewHandler(us services.UserService, ev services.EventService, par services.ParticipantService, au services.AuthService, or services.OrchestratorService, ctx context.Context) *Handler {
	return &Handler{
		UserService: us,
		EvenService: ev,
		ParticipantService: par,
		AuthService: au,
		Orchestrator: or,
		Ctx: ctx,
	}
}