package httpdelivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

func (h *Handler) CreateEventWithOwner() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Location    string    `json:"location"`
			Slots       int       `json:"slots"`   
			DateFrom    time.Time `json:"date_from"`
			DateTo      time.Time `json:"date_to"`
			OwnerID     uuid.UUID `json:"owner_id"`
			Access      string    `json:"access"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		event, err := h.Orchestrator.CreateEventWithOwner(r.Context(), req.Title, req.Description, req.Location, req.Slots, req.DateFrom, req.DateTo, req.OwnerID, domain.Access(req.Access))
		if err != nil {
			http.Error(w, "failed to create event", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(event)
	}
}

func (h *Handler) PublishEventAtomic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			EventID uuid.UUID `json:"event_id"`
			RequesterID uuid.UUID `json:"requester_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		event, err := h.Orchestrator.PublishEventAtomic(r.Context(), req.EventID, req.RequesterID)
		if err != nil {
			http.Error(w, "failed to publish event", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(event)
	}
}

func (h *Handler) DraftEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			EventID uuid.UUID `json:"event_id"`
			RequesterID uuid.UUID `json:"requester_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		_, err := h.Orchestrator.DraftEvent(r.Context(), req.EventID, req.RequesterID)
		if err != nil {
			http.Error(w, "failed to draft event", http.StatusInternalServerError)
			return 
		}
		responce := fmt.Sprintf("event: %v was successfully drafted by %v", req.EventID, req.RequesterID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
	}
}

func (h *Handler) CancelEventWithCleanup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			EventID uuid.UUID `json:"event_id"`
			RequesterID uuid.UUID `json:"requester_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		event := h.Orchestrator.CancelEventWithCleanup(r.Context(), req.EventID, req.RequesterID)
		if event != nil {
			http.Error(w, "failed to cancel event", http.StatusInternalServerError)
			return
		}

		responce := fmt.Sprintf("event: %v was successfully canceled", req.EventID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
	}
}

func (h *Handler) DeleteEventWithPermissions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			EventID uuid.UUID `json:"event_id"`
			RequesterID uuid.UUID `json:"requester_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		event := h.Orchestrator.DeleteEventWithPermissions(r.Context(), req.EventID, req.RequesterID)
		if event != nil {
			http.Error(w, "failed to delete event", http.StatusInternalServerError)
			return
		}

		responce := fmt.Sprintf("event: %v was successfully deleted", req.EventID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
	}
}

func (h *Handler) GetEventByUserRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserRole string `json:"user_role"`
			RequesterID uuid.UUID `json:"requester_id"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		events, err := h.EvenService.GetEventsByUserRole(r.Context(), req.RequesterID, domain.SystemRole(req.UserRole))
		if err != nil {
			http.Error(w, "failed to get events", http.StatusInternalServerError)
			return
		}
		response := struct {
			Events []domain.Event `json:"events"`
			Count  int            `json:"count"`
			User   uuid.UUID      `json:"requester_id"`
		}{
			Events: events,
			Count:  len(events),
			User: req.RequesterID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func (h *Handler) GetEventsByLocatoin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Location string `json:"location"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		events, err := h.EvenService.GetEventsByLocation(r.Context(), req.Location)
		if err != nil {
			http.Error(w, "failed to get events", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(events)
	}
}

func (h *Handler) HasAailableSlots() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			EventID uuid.UUID `json:"event_id"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		isSlot, err := h.EvenService.HasAvailableSlots(r.Context(), req.EventID)
		if err != nil {
			http.Error(w, "failed to check slots", http.StatusInternalServerError)
			return 
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(isSlot)
	}
}