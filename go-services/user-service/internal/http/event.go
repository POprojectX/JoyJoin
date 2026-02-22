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

